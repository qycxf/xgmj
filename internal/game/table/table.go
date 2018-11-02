//基础牌桌功能

package table

import (
	"time"

	"strconv"

	"fmt"
	"runtime"

	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/maj"
	"qianuuu.com/xgmj/internal/game/seat"
	. "qianuuu.com/xgmj/internal/mjcomn"
	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/qo"
	"qianuuu.com/player"
)

// TimerAction 桌子操作
type TimerAction struct {
	timerType consts.TimerType
	param     int
}

// Table 桌子
type Table struct {
	ID     int // 桌子 ID
	Name   string
	finish chan struct{}

	TableCfg *config.TableCfg //桌子属性配置

	uuid             string            //每局牌唯一标志
	state            consts.TableState //桌子状态
	updateTime       int64             //记录某个倒计时启动时间点
	createTime       int64             //创建桌子时间
	seats            seat.Seats        //玩家位置管理
	Majhong          *maj.Majong       //一局牌管理
	robots           *Robots           //机器人管理
	ReqDismissTime   int64             //请求解散牌桌时间点
	ReqDismissSeatID int               //请求解散牌桌玩家seatID
	AgreeDismissArr  []int             //是否同意 解散房间信息

	tableAction chan consts.TimerType // 桌子内触发操作
	timerAction chan TimerAction

	autoChangeState  *time.Timer //自动切换桌子状态计时器
	autoDismissTable *time.Timer //超时解散牌桌计时器

	remainTime     int32        //倒计时剩余时间
	ExecOptInfo    *ExecOptInfo //位置执行操作信息
	SaveOpt        []int        //多人操作时 暂存先发送的操作,等待高优先级的操作响应
	FangSeatID     int          //房主位置
	PresentUID     int          //如果是赠送房间,标记赠送者uid
	PresentName    string       //如果是赠送房间,标记赠送者昵称
	GameHappend    bool         //标记游戏第一局是否打完
	SelectThreeArr [][]int      //临时保存各位置选的三张牌值
	DirectOver     bool         //是否直接结束,游戏中途解散导致的游戏结束

	handler    Handler    //调用外部接口
	tableInter TableInter //牌桌接口
	ConsumeID  int

	// 同步处理模块
	ql *qo.Qo
}

// Go 异步调用 (线性，有顺序)
func (t *Table) Go(fn func()) {
	if t == nil {
		return
	}
	if t.ql == nil {
		t.ql = qo.New()
	}
	t.ql.Go(fn)
}

// NewTable 创建一个桌子
func NewTable(_tableID int, _robots *Robots, _tableCfg *config.TableCfg) *Table {
	ret := &Table{
		TableCfg:    _tableCfg,
		ID:          _tableID,
		finish:      make(chan struct{}),
		robots:      _robots,
		seats:       make(map[int]*seat.Seat),
		tableAction: make(chan consts.TimerType),
		timerAction: make(chan TimerAction),
		PresentName: "",
	}

	_tableCfg.Print()
	ret.Init()

	return ret
}

//初始化
func (t *Table) Init() {

	//牌桌自动切换状态
	t.autoChangeState = time.AfterFunc(time.Second, func() {
		t.tableAction <- consts.TimerAutoChangeState
	})
	t.autoChangeState.Stop()

	//牌桌超时自动解散
	t.autoDismissTable = time.AfterFunc(time.Second, func() {
		t.tableAction <- consts.TimerAutoDismisTable
	})
	t.autoDismissTable.Stop()

	t.Majhong = maj.NewMajong(t.TableCfg) //一局麻将牌管理对象
	t.SaveOpt = make([]int, 0)
	t.SelectThreeArr = make([][]int, 4)
	t.DirectOver = false
	t.GameHappend = false
	t.remainTime = 0
	t.uuid = "uuid"
	t.createTime = time.Now().Unix() * 1000 //转成毫秒

	//初始化桌子状态 - 空闲
	t.setState(consts.TableStateIdle)

	t.ClearDismissInfo()
}

// Serve 桌子逻辑
func (t *Table) Serve() {

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1<<16)
			n := runtime.Stack(trace, true)
			logs.Info("%v", fmt.Errorf("tableId:%v-------------------->panic recover\n %v\n stack trace %d bytes\n %s",
				t.ID, err, n, trace[:n]))
		}
	}()

	// 消息监听
	for {

		select {

		//定时执行 间隔1秒
		case <-time.After(time.Second * 1):
			//机器人游戏逻辑
			t.Go(func() {
				//if t.GetState() > consts.TableStateDealCard && t.GetState() > consts.TableStateShowResult {
				t.buhuaLogic()
				//}
				t.tingLogic()
				t.robotLogic()

			})
			break

			// 桌子内部消息
		case action := <-t.tableAction:
			t.Go(func() {
				switch action {

				//超时自动解散牌桌
				case consts.TimerAutoDismisTable:
					logs.Info("tableId:%d -----------------> TimerAutoDismisTable ==>  超时解散牌桌", t.ID)
					for k, _ := range t.AgreeDismissArr {
						t.AgreeDismissArr[k] = consts.Yes //超时全部同意
					}
					t.DisMissTable(consts.DefaultIndex)
					break

					//游戏状态计时器,自动切换状态
				case consts.TimerAutoChangeState:

					logs.Info("tableId:%v TimerAutoChangeState ========================>  t.GetState():%v", t.ID, t.GetState())
					if t.GetState() == consts.TableStateWaiteReady {

						logs.Info("tableId:%v ========================>等待准备时间结束,游戏自动开始", t.ID)
						t.autoReady()

					} else if t.GetState() == consts.TableStateDealCard {

						logs.Info("tableId:%v ========================>等待客户端发牌时间结束等待轮流补花", t.ID)
						//if t.TableCfg.TableType == TableType_HYMJ {
						//	logs.Info("tableId:%v ========================>所有玩家开始补花", t.ID)
						////	所有玩家开始轮流补花
						//	t.LunLiuBuHua()
						//	logs.Info("tableId:%v ========================>玩家开始思考", t.ID)
						//}
						//庄家开始思考
						bu := true
						if t.TableCfg.TableType == TableType_BBMJ && t.TableCfg.DaiHua <= 0 {
							bu = false
						}
						if bu {
							time.Sleep(time.Second * 1)
							t.LunLiuBuHua()

						}
						t.tableInter.DPlayerFirstThink()

					} else if t.GetState() == consts.TableStateWaiteSend {

						logs.Info("tableId:%v ========================>出牌时间结束,系统自动出牌", t.ID)
						t.autoSendBySys()
					} else if t.GetState() == consts.TableStateWaiteThink {

						logs.Info("tableId:%v ========================>思考时间结束,系统自动思考", t.ID)
						t.autoThinkBySys()
					} else if t.GetState() == consts.TableStateZhuaNiao {

						logs.Info("tableId:%v========================>抓鸟时间结束,进入结算", t.ID)
						t.GameOver()

					} else if t.GetState() == consts.TableStateShowResult {

						logs.Info("tableId:%v========================>结算状态时间展示结束", t.ID)
						t.Reset() //先重置牌桌游戏数据,再发送TableInfo
						t.setState(consts.TableStateIdle)
					}

					break

				default:
					break
				}

			})
		case <-t.finish:
			logs.Info("tableId:%v,~~~~-t.finish~~~~", t.ID)
			return
		}

	}
}

// 玩家进入桌子
func (t *Table) JoinTable(_player *player.Player) {

	logs.Info("tableId:%v--------------->%v 进入桌子", t.ID, _player.String())
	_player.SetTableID(t.ID)
	idleSeatID := t.SearchIdleSeat()

	if idleSeatID == -1 {
		logs.Info("tableId:%v,******************************** error ! searchIdleSeat ,no idle SeatID!", t.ID)
		return
	}

	//寻找座位坐下
	_seat := seat.NewSeat(idleSeatID, _player)
	t.sitdown(_seat)

	//记录房主位置
	if _player.IsFangZhu() {
		t.FangSeatID = idleSeatID
	}

	//准备
	t.SeatReady(_player.ID())
}

//玩家坐到一个位置中
func (t *Table) sitdown(seat *seat.Seat) {

	seatId := seat.GetId()
	seat.SetState(consts.SeatStateIdle) //初始为空闲状态
	t.seats[seatId] = seat
	logs.Info("tableId:%v------------------>%v 坐到%d号位置,t.seats %v", t.ID, seat.GetPlayer().String(), seatId, t.seats)

	//发送玩家坐下信息
	t.SendTableInfo()

	testuid := seat.GetPlayer().ID()

	if config.TestData().HasId(testuid) {
		logs.Info("--------------------------2222----------->testuid: %v,t.Majhong.TableCfg.TestSeatId:%v",
			testuid, t.Majhong.TableCfg.TestSeatId)
		if t.Majhong.TableCfg.TestSeatId == consts.DefaultIndex { //只设置一个位置
			t.Majhong.TableCfg.TestSeatId = seatId
			logs.Info("--------------------------4444----------->testuid: %v,t.Majhong.TableCfg.TestSeatId:%v",
				testuid, t.Majhong.TableCfg.TestSeatId)
		}
	}
}

//超时系统自动准备或踢出 (第一局超时踢出,后面每局自动准备)
func (t *Table) autoReady() {

	for _, v := range t.seats {
		//未准备的玩家直接准备
		if v.GetState() == consts.SeatStateIdle {
			t.SeatReady(v.GetPlayer().ID())
		}
	}
}

//func (t *Table) autoLunLiuBuHua() {
//
//	t.LunLiuBuHua()
//	for _, v := range t.seats {
//		//未准备的玩家直接准备
//		if v.GetState() == consts.SeatStateIdle {
//			t.SeatReady(v.GetPlayer().ID())
//		}
//	}
//}

// 玩家准备
func (t *Table) SeatReady(_uid int) {

	seat := t.seats.GetSeatByUID(_uid)
	logs.Info("tableId:%v-------------------> %v 准备游戏! ,_uid:%v,t.seats:%v", t.ID, seat, _uid, t.seats)

	//设置玩家准备状态
	seat.SetState(consts.SeatStateGameReady)
	seat.NotReadyCt = 0 //清除不准备次数

	waiteReady := false //是否进入等待装备
	startNow := false   //是否立即开始

	readyCt := t.getReadyCt()
	if readyCt == t.TableCfg.PlayerCt-1 {
		waiteReady = true
	} else if readyCt == t.TableCfg.PlayerCt {
		startNow = true
	}

	//切换到等待准备
	if t.GetState() == consts.TableStateIdle {
		if waiteReady {
			t.setState(consts.TableStateWaiteReady) //等待玩家准备
			return
		}
	}

	//立即开始游戏
	if startNow {
		logs.Info("tableId:%v, ----------------> 所有玩家均已准备,游戏立即开始 <----------------", t.ID)
		t.autoChangeState.Stop() //停止准备倒计时
		t.tableInter.GameStart()

		return
	}

	t.SendTableInfo()
}

// Reset 重置游戏中产生的一局 数据数据
func (t *Table) Reset() {

	//重置牌局一局数据
	t.Majhong.Reset()

	for _, v := range t.seats {
		v.Reset()
	}
	t.ExecOptInfo = NewExecOptInfo()
}

////一局游戏开始
//func (t *Table) gameStart() {
//
//	if t.Majhong.GameCt == 0 {
//		t.Majhong.GameCt++ //游戏局数+1 (第一局 游戏开始时增加)
//	}
//
//	logs.Info("------------->gameStart(%v %v)<------------- 游戏开始!", t.ID, t.Majhong.GameCt)
//
//	//检测玩家Ip
//	//t.chkIp()
//
//	//确定庄家
//	t.tableInter.MakeDSeat()
//
//	//确定完庄家后再 清除部分上一局用于下局判断庄家位置的的临时数据
//	t.Majhong.ClearLastTmpData()
//
//	//发牌
//	t.Majhong.DealCard()
//
//	//发牌状态
//	t.setState(consts.TableStateDealCard)
//}

//检测玩家IP是否相同
func (t *Table) chkIp() {
	// strTmp := ""
	// ipArr := make([]string, t.TableCfg.PlayerCt)
	// // for _, v := range t.GetSeats() {
	// // ipArr[v.GetId()] = v.GetPlayer().GetLastLoginIP()
	// // }

	// seatUseArr := make([]int, t.TableCfg.PlayerCt)
	// for i := 0; i < len(ipArr); i++ {
	// 	if seatUseArr[i] == consts.Yes { //已经提示相同的不再计算
	// 		continue
	// 	}
	// 	ip1 := ipArr[i]
	// 	str := "[" + t.GetSeats().GetSeatBySeatId(i).GetPlayer().NickName() + "]"
	// 	findSameIp := false
	// 	for j := i + 1; j < len(ipArr); j++ {
	// 		ip2 := ipArr[j]
	// 		if ip2 == ip1 {
	// 			seatUseArr[i] = consts.Yes
	// 			seatUseArr[j] = consts.Yes
	// 			if !findSameIp {
	// 				findSameIp = true
	// 			}
	// 			str += "[" + t.GetSeats().GetSeatBySeatId(j).GetPlayer().NickName() + "]"
	// 		}
	// 	}
	// 	if findSameIp {
	// 		str += " IP相同! "
	// 		strTmp += str
	// 	}
	// }
	// if strTmp != "" {
	// 	t.SendTopTip(strTmp)
	// }

}

// 检测\验证 操作合法性
func (t *Table) CheckOpt(_seat *seat.Seat, _optType int, _cardData int) int {

	_seatID := _seat.GetId()
	result := consts.Success

	if _optType == OptTypeTing {
		return result
	}
	//出牌
	if _optType == OptTypeSend {
		if t.GetState() != consts.TableStateWaiteSend { //检测状态
			result = consts.ErrorTableStateSend
		} else if _seatID != t.Majhong.CurtSenderIndex { //检测位置
			result = consts.ErrorNotCurrentSenderIndex
		} else {
			//检测牌值
			hasPai := t.Majhong.CMajArr[_seatID].IsHandHasSamePai(_cardData)
			if !hasPai {
				result = consts.ErrorSendNoThisCard
			}

		}
	} else {
		//碰杠胡弃  判断当前状态 以及是否在思考列表中
		if t.GetState() != consts.TableStateWaiteThink {
			result = consts.ErrorTableStateThinkOpt
		} else {

			//用操作判断验证判断 方式 (验证碰\杠\胡)------------------------------------
			canOpt := t.Majhong.CMajArr[_seatID].CanOpt(_optType) //操作是否可用,如果已操作过,则不能重复操作
			if !canOpt {
				result = consts.ErrorSeatOptNoThisOpt
				return result
			}

			//用思考位置验证判断 方式 (验证碰\杠\胡\弃)------------------------------------
			//多人操作,碰杠 不删除思考位置 , 胡弃 删除思考位置
			thinkerCt := t.Majhong.GetThinkerCt()
			if thinkerCt == 1 {
				if !t.Majhong.HasThinker(_seatID) {
					result = consts.ErrorNotCurrentThinkIndex
				}
			} else {
				//主要是处理客户端连续发送的问题
				if t.Majhong.HasThinker(_seatID) {
					//如果有已经操作且被保存则返回不能重复操作 (碰 杠)
					if len(t.SaveOpt) > 0 {
						if t.SaveOpt[0] == _seatID {
							result = consts.ErrorHasSaveThinkOpt
						}
					}
				} else {
					//胡 弃
					result = consts.ErrorNotCurrentThinkIndex
				}
			}
		}
	}
	return result
}

// 执行位置操作, 需要注意的问题 :
// 1.如果是单人思考,则操作后清空操作且删除思考位置  2.如果多人思考 [碰\杠] 清空操作但不删除位置, [胡\取消] 清空操作且删除位置
// 2._isSave:true 执行保存的操作
func (t *Table) SeatOpt(_seat *seat.Seat, _opt int, _data int, _isSave bool) {

	seatId := _seat.GetId()
	logs.Info("------------->tabeId:%v,玩家操作,%v ->%s,_data:%v,player=%v", t.ID, _seat, GetOptName(_opt), _data, _seat.GetPlayer().String())

	optCmaj := t.Majhong.CMajArr[seatId]

	//检测多人思考的情况,出牌不检测
	if _opt != OptTypeSend {
		isReturn := t.chkMultThinker(_seat, _opt, _data)

		//操作被保存或被处理,不再继续执行该操作,发送TableInfo刷新
		if isReturn {
			t.SendTableInfo()
			return
		}

		//停止状态计时器,如果多人操作 有操作被返回,则不能中断计时器;有人仍然需要继续思考
		if t.Majhong.GetThinkerCt() == 0 {
			t.autoChangeState.Stop()
		}
	}

	//保存当前操作的玩家位置及操作值,客户端使用,如果操作被返回,则不设置 t.SeatOptInfo的值,因此放到这里保存
	t.ExecOptInfo = NewExecOptInfo()
	t.ExecOptInfo.OptSeatId = seatId
	t.ExecOptInfo.OptType = _opt
	t.ExecOptInfo.OptData = _data

	//执行操作
	switch _opt {

	case OptTypeSend:            //出牌
		t.autoChangeState.Stop() //出牌后立即停止计时器
		t.send(_seat, _data)
		break

	case OptTypeHu: //胡牌

		_huType := optCmaj.HuType
		t.ExecOptInfo.OptDetail = _huType

		//记录胡牌顺序
		t.Majhong.HuSeq++
		if t.Majhong.HuSeq == 1 {
			t.Majhong.FirstHuSeatID = _seat.GetId() //记录第一个胡牌玩家
		}
		t.Majhong.CMajArr[_seat.GetId()].HuSeq = t.Majhong.HuSeq
		t.Majhong.HasHuArr = append(t.Majhong.HasHuArr, seatId) //记录胡牌玩家位置

		//如果是放炮,加上放炮位置
		if _huType == HUTYPE_JIEPAO {
			//判断是否抢杠胡 ,因为每个人只能胡一次,所以这里可以根据 胡牌详细判断
			if t.Majhong.CMajArr[_seat.GetId()].HuTypeDetail == HUTYPE_DETAIL_QIANGGANG {
				t.ExecOptInfo.DianPaoSeatId = t.Majhong.LastMianGangSeatID //上次面杠位置

			} else if t.Majhong.CMajArr[_seat.GetId()].HuTypeDetail == HUTYPE_DETAIL_GANG_SHANG_PAO {
				t.ExecOptInfo.DianPaoSeatId = t.Majhong.LastSenderSeatID //上次点杠(出牌)位置
				t.Majhong.CMajArr[t.Majhong.LastSenderSeatID].AddPxScore("点"+strconv.Itoa(t.Majhong.HuSeq)+"胡", 0)

			} else {
				t.ExecOptInfo.DianPaoSeatId = t.Majhong.LastSenderSeatID
				t.Majhong.CMajArr[t.Majhong.LastSenderSeatID].AddPxScore("点"+strconv.Itoa(t.Majhong.HuSeq)+"胡", 0)
			}
		}

		t.tableInter.HuPai(_seat)
		break

	case OptTypePeng: //碰牌

		t.tableInter.Peng(seatId)
		break

	case OptTypeGang: //杠牌

		t.tableInter.Gang(seatId, _data, _isSave)
		break
	case OptTypeBu: //补花
		t.Buhua(seatId)
		break
	case OptTypeCancel: //取消思考(手动取消或 超时取消)
		t.cancer(_seat, _data)
		break
	case OptTypeTing: //报听操作
		t.BaoTing(seatId)
		break

	default:
		break
	}

}

//执行操作数据对象
type ExecOptInfo struct {
	OptSeatId     int   //操作位置
	OptType       int   //操作类型
	OptData       int   //操作对应的牌值
	OptDetail     int   //操作详细类型
	DianPaoSeatId int   //点炮位置
	HupxIdArr     []int //胡牌牌型id数组
}

func NewExecOptInfo() *ExecOptInfo {
	ret := ExecOptInfo{
		OptSeatId:     consts.DefaultIndex,
		OptType:       consts.DefaultIndex,
		OptData:       consts.DefaultIndex,
		OptDetail:     consts.DefaultIndex,
		DianPaoSeatId: consts.DefaultIndex,
		HupxIdArr:     make([]int, 0),
	}
	return &ret
}

//发送完客户端需要的临时数据,发送后清除
func (t *Table) clearCltData() {

	//发送完操作信息后重置操作数据
	t.ExecOptInfo = NewExecOptInfo()

}

// 多位置(2~3个位置)同时思考
func (t *Table) chkMultThinker(_seat *seat.Seat, _opt int, _data int) bool {
	length := len(t.Majhong.CurtThinkerArr)
	logs.Info("tableId:%v,--------------->chkMultThinker,t.Majhong.CurtThinkerArr:%v ", t.ID, t.Majhong.CurtThinkerArr)
	if length == 1 {
		return false
	}

	//查看当前位置操作值,按情况分类 1:碰胡(胡) 2:杠胡(胡) 3:胡胡(胡)
	//1) 碰或杠的人先操作,保存操作
	if _opt == OptTypePeng || _opt == OptTypeGang {

		if _opt == OptTypePeng {
			//碰牌操作检测过手胡
			t.chkGSH(_seat.GetId())
		}
		//暂存碰杠操作
		t.SaveOpt = []int{_seat.GetId(), _opt, _data}
		logs.Info("tableId:%v,--------------->多人操作,位置 %v 操作被保存,操作数组:%v ", t.ID, _seat.GetId(), t.Majhong.CMajArr[_seat.GetId()].OptInfo)

		//删除位置操作
		t.Majhong.CMajArr[_seat.GetId()].ResetOptInfo() //清空本位置操作 ,不删除思考者,用于第二个玩家操作时判断
		return true
	} else if _opt == OptTypeHu {
		//是否一炮多响
		if false {
			//if t.TableCfg.DianpaoHu == consts.Yes {
			//if t.TableCfg.TableType == TableType_FYMJ {
			lastSenderSeatID := t.Majhong.LastSenderSeatID
			huTypeDetail := t.Majhong.CMajArr[_seat.GetId()].HuTypeDetail
			if huTypeDetail == HUTYPE_DETAIL_QIANGGANG {
				lastMianGangSeatID := t.Majhong.LastMianGangSeatID
				lastSenderSeatID = lastMianGangSeatID
			}

			//查看当前位置优先级
			if lastSenderSeatID != consts.DefaultIndex {

				t.Majhong.CMajArr[_seat.GetId()].OptInfo.HuCard = _data
				//t.Majhong.CMajArr[_seat.GetId()].HuMCard = NewMCard(_data)

				//将当期保存的面杠操作删除
				logs.Info("tableId:%v,--------------->多人操作,阜阳麻将抢杠胡!删除保存的面杠操作 t.SaveOpt:%v ", t.ID, t.SaveOpt)
				t.SaveOpt = make([]int, 0)
				//t.Majhong.RemoveThinker(lastSenderSeatID)

				//  获取点炮玩家右边思考的的玩家
				rightArr := t.Majhong.SearchArrRight(lastSenderSeatID)
				_index := consts.DefaultIndex
				for i := 0; i < len(rightArr); i++ {
					if rightArr[i] == _seat.GetId() {
						_index = i
					}
				}

				logs.Info("tableId:%v,---------------> rightArr:%v , _index:%v", t.ID, rightArr, _index)

				if _index != consts.DefaultIndex {

					//将当期位置后面可胡的玩家操作删除
					for i := _index + 1; i < len(rightArr); i++ {
						if t.Majhong.HasThinker(rightArr[i]) {
							//先查看是否有操作被保存
							if len(t.SaveOpt) > 0 && t.SaveOpt[0] == rightArr[i] {
								logs.Info("tableId:%v,--------------->多人胡牌操作,位置 %v 已保存的操作被删除,操作数组:%v ", t.ID, rightArr[i], t.SaveOpt)
								t.SaveOpt = make([]int, 0)
								t.Majhong.RemoveThinker(rightArr[i])
							} else {
								t.Majhong.CMajArr[rightArr[i]].ResetOptInfo()
								t.Majhong.RemoveThinker(rightArr[i])
								logs.Info("tableId:%v,--------------->多人胡牌操作,位置 %v 优先级较低 ,CurtThinkerArr:%v", t.ID, rightArr[i], t.Majhong.CurtThinkerArr)
							}
						}
					}

					//查找 _index 之前是否有人可胡
					for i := 0; i < _index; i++ {
						if t.Majhong.HasThinker(rightArr[i]) {

							//暂存胡牌操作
							t.SaveOpt = []int{_seat.GetId(), _opt, t.Majhong.LastSendCard.GetData()}
							logs.Info("tableId:%v,--------------->多人胡牌操作,位置 %v 操作被保存,操作数组:%v ", t.ID, _seat.GetId(), t.Majhong.CMajArr[_seat.GetId()].OptInfo)

							//删除位置操作
							t.Majhong.CMajArr[_seat.GetId()].ResetOptInfo() //清空本位置操作 ,不删除思考者,用于第二个玩家操作时判断
							return true
						}
					}

				} else {
					logs.Info("tableId:%v,******************* [error] *******************>多人胡牌操作,位置 %v rightArr:%v ", t.ID, _seat.GetId(), rightArr)

				}

			}
			//}
		}

		//检测其他位置操作,如果不是胡,则直接删除操作
		for _, v := range t.Majhong.CurtThinkerArr {
			if v != _seat.GetId() {
				//先查看是否有操作被保存
				if len(t.SaveOpt) > 0 && t.SaveOpt[0] == v {
					logs.Info("tableId:%v,--------------->多人操作,位置 %v 已保存的操作被删除,操作数组:%v ", t.ID, v,
						t.SaveOpt)
					t.SaveOpt = make([]int, 0)
					t.Majhong.RemoveThinker(v)
				} else {

					//是否是胡操作,不是则直接删除操作
					canPeng := t.Majhong.CMajArr[v].CanOpt(OptTypePeng)
					if canPeng {
						logs.Info("tableId:%v,--------------->多人操作,位置 %v 未进行的碰操作被删除,操作数组:%v ", t.ID, v, t.Majhong.CMajArr[v].OptInfo)
						t.Majhong.CMajArr[v].RemoveOpt(OptTypePeng)
					}

					canGang := t.Majhong.CMajArr[v].CanOpt(OptTypeGang)
					if canGang {
						logs.Info("tableId:%v,--------------->多人操作,位置 %v 未进行的杠操作被删除,操作数组:%v ", t.ID, v, t.Majhong.CMajArr[v].OptInfo)
						t.Majhong.CMajArr[v].RemoveOpt(OptTypeGang)
					}

					//如果该玩家也能胡,则返回
					canHu := t.Majhong.CMajArr[v].CanOpt(OptTypeHu)
					if canHu {

					} else {
						t.Majhong.CMajArr[v].ResetOptInfo()
						t.Majhong.RemoveThinker(v)
						logs.Info("tableId:%v,--------------->多人操作,位置 %v 没有其他操作,从思考列表中删除 ,CurtThinkerArr:%v", t.ID, v, t.Majhong.CurtThinkerArr)
					}
				}
			}
		}

		//多人操作,取消操作(这里仅处理 [玩家手动取消] )
	} else if _opt == OptTypeCancel {

		//取消操作检测过手胡
		t.chkGSH(_seat.GetId())
		t.chkGSP(_seat.GetId())

		//取消操作 取消抢杠胡
		t.chkQGH(_seat)

		//如果是碰\杠 玩家取消,则移除操作\思考者;如果是胡玩家取消,先移除操作\思考者,
		// 然后判断是否有保存,有则执行,调用单人思考方法
		t.Majhong.CMajArr[_seat.GetId()].ResetOptInfo() //清空本位置操作
		t.Majhong.RemoveThinker(_seat.GetId())          //取消操作则删除思考人位置

		logs.Info("tableId:%v,--------------->consts.SeatOptCancer t.SaveOpt:%v", t.ID, t.SaveOpt)
		//可能有三个玩家在思考,一个胡放弃了,还需等待另一个胡放弃,所以必须当前没有思考者,才检测执行保存的操作
		if t.Majhong.GetThinkerCt() <= 1 {
			if len(t.SaveOpt) > 0 { //如果有保存操作,则执行保存操作
				//IMT 且被保存的操作的位置,必须是 thinker 中的位置值
				if t.SaveOpt[0] == t.Majhong.CurtThinkerArr[0] {
					//TODO 这种情况下没有发送操作信息,客户端回放时未展示放弃操作
					_seat := t.seats[t.SaveOpt[0]]
					t.SeatOpt(_seat, t.SaveOpt[1], t.SaveOpt[2], true)
					t.SaveOpt = make([]int, 0) //清空保存的操作
					//return true
				}
			}
		} else {
			//继续等待其它思考者 操作
		}

		return true //放弃或系统超时放弃都 直接返回,因为这里已经处理完放弃操作.
	} else {
		logs.Info("tableId:%v,*********************error!!! chkMultThinker _seat:%v,_opt:%v,_data:%v", t.ID, _seat, _opt, _data)
	}

	return false
}

//放弃操作时检测过手胡
func (t *Table) chkGSH(_seatId int) {

	//如果是放弃胡操作,则标记过手胡显示
	canHu := t.Majhong.CMajArr[_seatId].CanOpt(OptTypeHu)
	if canHu {
		t.Majhong.CMajArr[_seatId].GuoShouHu = true                                      //标记过手胡
		t.Majhong.CMajArr[_seatId].GuoShouHuCard = t.Majhong.LastSendCard.Clone()        //标记过手胡的牌
		t.Majhong.CMajArr[_seatId].GuoShouFanCt = t.Majhong.GetSeatFanCt(_seatId, false) //记录番数
		logs.Info("tableId:%v,--------------->位置%v 放弃操作,被标记过手状态,GuoShouFanCt:%v", t.ID, _seatId, t.Majhong.CMajArr[_seatId].GuoShouFanCt)
	}
}

//放弃操作时检测过手碰
func (t *Table) chkGSP(_seatId int) {

	//如果是放弃胡操作,则标记过手胡显示
	canHu := t.Majhong.CMajArr[_seatId].CanOpt(OptTypePeng)
	if canHu {
		t.Majhong.CMajArr[_seatId].GuoShouPeng = true                               //标记过手胡
		t.Majhong.CMajArr[_seatId].GuoShouPengCard = t.Majhong.LastSendCard.Clone() //标记过手碰的牌
		//t.Majhong.CMajArr[_seatId].GuoShouFanCt = t.Majhong.GetSeatFanCt(_seatId, false) //记录番数
		logs.Info("tableId:%v,--------------->位置%v 放弃操作,被标记过手状态", t.ID, _seatId)
	}
}

//放弃操作时清除 抢杠胡 信息
func (t *Table) chkQGH(_seat *seat.Seat) {
	if t.Majhong.CMajArr[_seat.GetId()].HuTypeDetail == HUTYPE_DETAIL_QIANGGANG {
		logs.Info("tableId:%v,--------------->位置%v 放弃操作,抢杠胡状态被删除", t.ID, _seat.GetId())
		t.Majhong.CMajArr[_seat.GetId()].HuTypeDetail = consts.DefaultIndex
	}
}

//玩家轮流补花
func (t *Table) LunLiuBuHua() {
	dSeatID := t.Majhong.DSeatID //庄家座位
	seatIdArr := make([]int, t.TableCfg.PlayerCt)
	seatIdArr = append(seatIdArr, dSeatID)
	//nextId := -1
	//for i := 0; i < t.TableCfg.PlayerCt; i++ {
	//	if nextId == -1 {
	//		nextId = t.Majhong.GetNextSeatID(dSeatID)
	//	} else {
	//		nextId = t.Majhong.GetNextSeatID(nextId)
	//	}
	//	seatIdArr = append(seatIdArr, nextId)
	//}
	for i := 0; i < t.TableCfg.PlayerCt; i++ {
		if i != dSeatID {
			seatIdArr = append(seatIdArr, i)
		}
	}

	falg := true
	for falg {
		falg = false
		for i := 0; i < len(seatIdArr); i++ {
			logs.Info("tableId:%v============i:%v==============seatIdArr[i]：%v补花", t.TableCfg.TableId, i, seatIdArr[i])
			if seatIdArr[i] < 0 || (seatIdArr[i] > t.TableCfg.PlayerCt-1) {
				continue
			}
			_cmaj := t.Majhong.CMajArr[seatIdArr[i]]
			_cmaj.BuHua = maj.NewBuHua() //清空花牌
			if t.TableCfg.TableType == TableType_BBMJ {
				_cmaj.Check_BB_BuHua()
			} else if t.TableCfg.TableType == TableType_HYMJ {
				_cmaj.Check_HY_BuHua(t.TableCfg.FengLing)
			}
			t.BuHuaCards(_cmaj, _cmaj.BuHua)
			logs.Info("tableId:%v==========================位置：%v补花", t.TableCfg.TableId, seatIdArr[i])

			_cmaj.BuHua = maj.NewBuHua() //清空花牌
			if t.TableCfg.TableType == TableType_BBMJ {
				_cmaj.Check_BB_BuHua()
			} else if t.TableCfg.TableType == TableType_HYMJ {
				_cmaj.Check_HY_BuHua(t.TableCfg.FengLing)
			}

			if len(_cmaj.BuHua.HandHuaCards) == 0 {
				seatIdArr[i] = 999
			} else {
				falg = true
			}
		}
		logs.Info("================本轮位置补花结束！================")

	}
	//立即开始游戏
	logs.Info("tableId:%v, ----------------> 所有玩家均已补花,庄家开始思考 <----------------", t.ID)
	//t.autoChangeState.Stop() //停止准备倒计时
	//t.setState(consts.TableStateChangeCard)

	t.SendTableInfo()

}

//设置玩家报听
func (t *Table) BaoTing(_seatId int) {
	_cmaj := t.Majhong.CMajArr[_seatId]
	if !_cmaj.IsFristSend {
		_cmaj.YinZi = true
	}
	_cmaj.IsTing = true
	_cmaj.ResetOptInfo() //清空本位置思考操作
	t.tableInter.SendCheckTing()
	//	修改报听之后的可以打出的牌
	t.TingDanCard()
	t.SendTableInfo()

}

// 选择报听后只能胡 单张牌
func (t *Table) TingDanCard() {

	_seatId := t.Majhong.CurtSenderIndex
	_cmaj := t.Majhong.CMajArr[_seatId]

	_SendTipArr := make([]*SendTip, 0)
	for _, v := range _cmaj.SendTipArr {
		if len(v.HuCards) == 1 {
			_SendTipArr = append(_SendTipArr, v)
		}
	}
	_cmaj.SendTipArr = _SendTipArr
}

////玩家补花
func (t *Table) Buhua(_seatId int) {

	_cmaj := t.Majhong.CMajArr[_seatId]
	//手牌中的花牌（东。中。发。白）
	_buCard := make([]int, 0)
	for _, v := range _cmaj.BuHua.HandHuaCards {
		_buCard = append(_buCard, v)
	}
	//设置该座位所有补花牌，根据这个判断是否能接炮
	for _, v := range _buCard {
		_cmaj.BuHua.BuCards = append(_cmaj.BuHua.BuCards, v)
	}

	logs.Info("tableId:%v-------------------->位置:%v补花,牌`:%v", t.ID, _seatId, _buCard)

	//执行补花操作
	t.BuHuaCards(_cmaj, _cmaj.BuHua)
	_cmaj.ResetOptInfo() //清空本位置思考操作

	_cmaj.LastFetchMCard = nil //  将LastFetchCard设置为空,防止发送tableinfo 时被移动,导致显示错误
	t.Majhong.SetCurtSenderIndex(_seatId)

	isSelfThink := t.tableInter.ThinkSelfPai(HUTYPE_DETAIL_GANG_SHANG_HUA, false)
	//需要思考
	if isSelfThink {
		t.setState(consts.TableStateWaiteThink)
	} else {
		t.tableInter.SendCheckTing()
		//无须思考,等待出牌
		t.setState(consts.TableStateWaiteSend)
	}
}

//补牌
func (t *Table) BuHuaCards(cmajs *maj.CMaj, seatBuHua *maj.BuHua) {
	ct := 0
	for _, v := range cmajs.GetHandPai() {
		ishua := false
		if t.TableCfg.TableType == TableType_BBMJ {
			ishua = v.IsBBHuaPai()
		} else if t.TableCfg.TableType == TableType_HYMJ {
			ishua = v.IsHYHuaPai(t.TableCfg.FengLing)
		}
		if ishua {
			//删除花牌
			logs.Info("手牌删除花牌：%v", v)
			seatBuHua.BuCards = append(seatBuHua.BuCards, v.GetData())
			cmajs.HuaPaiArr = append(cmajs.HuaPaiArr, v)
			cmajs.RemoveHandCard(NewMCard(v.GetData()))
			//cmajs.BuhuaCt++
			ct++
		}
	}
	for i := 0; i < ct; i++ {
		card := t.Majhong.FetchACard(cmajs.SeatID, false)
		t.Majhong.CMajArr[cmajs.SeatID].AddHandPai(card, true)

	}

	//time.Sleep(time.Second * 1)
	cmajs.SortHandPai() //补完花后排序手牌
	logs.Info("位置：%v:玩家补花后手牌：%v", cmajs.SeatID, cmajs.GetHandPai())
}

//玩家取消操作,这里只处理 思考人只有一个玩家的情况
func (t *Table) cancer(_seat *seat.Seat, _paiData int) {

	if len(t.Majhong.CurtThinkerArr) > 1 {
		//这里只处理只有一个位置思考的情况,多位置思考 前置方法处理
		logs.Info("tableId:%v *************************cancer error  CurtThinkerArr:%v,_seat:%v,_paiData:%v",
			t.ID, t.Majhong.CurtThinkerArr, _seat, _paiData)
		return
	}

	logs.Info("tableId:%v,len(t.Majhong.CurtThinkerArr):%v", t.ID, len(t.Majhong.CurtThinkerArr))

	//取消操作检测过手胡
	t.chkGSH(_seat.GetId())
	t.chkGSP(_seat.GetId())

	//取消操作 取消抢杠胡
	t.chkQGH(_seat)

	t.Majhong.CMajArr[_seat.GetId()].ResetOptInfo() //清空操作
	t.Majhong.RemoveThinker(_seat.GetId())          //删除思考位置

	//发送取消
	t.SendTableInfo()

	//判断是思考自己还是思考他人
	if t.Majhong.CurtSenderIndex == _seat.GetId() { //思考自己
		logs.Info("tableId:%v,放弃思考 自己的牌", t.ID)
		//轮到自己出牌
		t.tableInter.SendCheckTing()

	} else {
		//思考他人
		logs.Info("tableId:%v,放弃思考 他人的牌 t.Majhong.CurtThinkerArr:%v", t.ID, t.Majhong.CurtThinkerArr) //这里可能是他人出的牌 或 抢杠胡的情况

		//如果有一炮多响,且之前有人已经点胡 ,已经将这张牌添加到胡牌人玩家手中,则后面的人取消则不再添加
		if len(t.Majhong.HasHuArr) > 0 {
			logs.Info("tableId:%v,之前有人已经点胡,本次取消操作不再添加牌:%v ", t.ID, t.Majhong.HasHuArr)

			//这里判断游戏结束,一炮多响,已经有人胡,最后一个可胡玩家取消,则游戏结束
			if len(t.Majhong.CurtThinkerArr) == 0 { //已经无人思考,即最后一个可胡玩家如果取消
				t.tableInter.LastHuThinkerCancer()
				return
			}
		} else {
			//当前只有一个人思考,放弃 ,将出牌人的牌添加到出牌数值中,并切换到下个人拿牌
			logs.Info("tableId:%v,位置放弃操作,将上次出的牌添加到 出牌人 出牌数组中去! ", t.ID)
			t.Majhong.CMajArr[t.Majhong.LastSenderSeatID].AddOutPai(t.Majhong.LastSendCard)
		}

		if t.tableInter.ChkOver() {
			logs.Info("tableId:%v,放弃思考别人打出的最后一张牌", t.ID)
			return
		}

		//如果有一炮多响,且之前有人已经点胡
		if len(t.Majhong.HasHuArr) > 0 {
			//可能 有1~2 个人点了胡牌,从放炮人左手起开始查找,找到按逆时针顺序,胡牌人是最大的位置的
			leftArr := t.Majhong.SearchArrLeft(t.Majhong.CurtSenderIndex)
			for _, v := range leftArr {
				if HasElement(t.Majhong.HasHuArr, v) {
					logs.Custom(logs.TableTag, "tableId:%v, 位置%v取消操作,一炮多响,之前有人已经点胡, arr:%v,v:%v", t.ID, _seat.GetId(), t.Majhong.HasHuArr, v)
					t.Majhong.SetCurtSenderIndex(v)
					break
				}
			}
		}

		logs.Info("tableId:%v------------------->思考他人操作取消,当前出牌人位置 t.Majhong.CurtSenderIndex:%v", t.ID, t.Majhong.CurtSenderIndex)
		t.playerFetchPai(true)
		t.tableInter.FetcherStartThink(consts.DefaultIndex)

	}
}

//玩家出牌
func (t *Table) send(_seat *seat.Seat, _paiData int) {
	seatID := _seat.GetId()
	sendCard := NewMCard(_paiData)

	logs.Info("tableId:%v------------------->位置:%v出牌,牌=%v", t.ID, seatID, sendCard)

	t.Majhong.LastSenderSeatID = seatID                //保存出牌人位置
	t.Majhong.LastSendCard = sendCard.Clone()          //保存这张牌
	t.Majhong.CMajArr[seatID].RemoveHandCard(sendCard) //删除这张手牌
	t.Majhong.CMajArr[seatID].LastSendMCard = sendCard.Clone()

	//如果这是最后一张打出的牌,则这张牌为海底胡牌型
	if t.Majhong.GetRemainPaiCt() == 0 {
		//t.Majhong.LastSendCard.AddPx(consts.EXTPXID_HAI_DI_HU) //标记 海底胡
	}
	if t.Majhong.CMajArr[seatID].GangShangPao {
		//t.Majhong.LastSendCard.AddPx(consts.EXTPXID_GANG_PAO) //打出的牌 标记杠上炮
		t.Majhong.CMajArr[seatID].GangShangPao = false //立即清除状态
	}

	//发送出牌动作
	t.SendTableInfo()

	//清空思考者位置
	t.Majhong.ClearThinker()

	//添加思考位置
	t.addOtherThinker(seatID)

	otherThink := t.tableInter.ThinkOtherPai(sendCard)

	t.tableInter.UpdateTingCards(seatID)

	//报听后出牌次数累加
	if t.Majhong.CMajArr[seatID].IsTing {
		t.Majhong.CMajArr[seatID].HowTingCt++
	}
	t.Majhong.CMajArr[seatID].IsFristSend = true
	//切换到等待人思考
	if otherThink {

		t.setState(consts.TableStateWaiteThink)

	} else {
		ishua := false
		if t.TableCfg.TableType == TableType_BBMJ && t.TableCfg.DaiHua > 0 {
			ishua = sendCard.IsBBHuaPai()
		} else if t.TableCfg.TableType == TableType_HYMJ {
			ishua = sendCard.IsHYHuaPai(t.TableCfg.FengLing)
		}
		if ishua {
			t.Majhong.CMajArr[seatID].AddHuaPaiArr(sendCard)
		} else {
			t.Majhong.CMajArr[seatID].AddOutPai(sendCard)
		}

		if t.tableInter.ChkOver() {
			return
		}

		t.playerFetchPai(true)

		t.tableInter.FetcherStartThink(consts.DefaultIndex)
	}
}

//添加 以 _seatID 开始,其他思考玩家位置到 牌局思考位置列表中
func (t *Table) addOtherThinker(_seatID int) {
	//添加其它思考位置,_seatID
	for i := _seatID + 1; i < t.TableCfg.PlayerCt; i++ {
		if i < t.TableCfg.PlayerCt && t.seats[i].GetState() != consts.SeatStateGameHasHu {
			t.Majhong.AddThinker(i)
		}
	}
	for i := 0; i < _seatID; i++ {
		if t.seats[i].GetState() != consts.SeatStateGameHasHu {
			t.Majhong.AddThinker(i)
		}
	}
}

//根据位置及杠牌值计算 杠牌类型
func (t *Table) getGangType(seatID int, _paiData int) int {
	_gangType := -1
	_gangCard := NewMCard(_paiData)

	if seatID == t.Majhong.CurtSenderIndex {
		count := t.Majhong.CMajArr[seatID].HandEqualPaiCt(_gangCard)
		if count == 4 {
			_gangType = GANGTYPE_AN //暗杠
		} else if count == 1 {
			_gangType = GANGTYPE_MIAN //面杠
		}

	} else {
		_gangType = GANGTYPE_ZHI //直杠
	}
	return _gangType
}

//刮风下雨 积分计算 (暗杠\面杠,扣除所有未胡玩家积分)
func (t *Table) GuaFXiay_AM(_seatID int, _score int) {

	gfxyArr := make([]int, t.TableCfg.PlayerCt)
	for _, v := range t.seats {
		if v.GetId() != _seatID && v.GetState() != consts.SeatStateGameHasHu {
			//暂存刮风下雨的位置所得,在单局游戏结束后计算,如果结束后未听牌,则不计算刮风下雨所得
			gfxyArr[v.GetId()] = _score //记录该位置需要被扣除的积分
		}
	}
	//保存到 GfxyRec 记录中去
	t.Majhong.CMajArr[_seatID].GfxyRec = append(t.Majhong.CMajArr[_seatID].GfxyRec, gfxyArr)
}

//刮风下雨 积分计算 (直杠,扣除点杠玩家积分)
func (t *Table) GuaFXiay_Z(_seatID int, _score int, _dianGangSeatID int) {

	gfxyArr := make([]int, t.TableCfg.PlayerCt)
	gfxyArr[_dianGangSeatID] = _score //记录该位置需要被扣除的积分

	//保存到 GfxyRec 记录中去
	t.Majhong.CMajArr[_seatID].GfxyRec = append(t.Majhong.CMajArr[_seatID].GfxyRec, gfxyArr)
}

//玩家拿牌 moveToNext:是否移动到下一个位置,杠牌时不移动
func (t *Table) playerFetchPai(moveToNext bool) {

	if moveToNext {
		for i := 0; i < t.TableCfg.PlayerCt; i++ {
			t.Majhong.MoveSenderIndexToNext()
			if t.seats[t.Majhong.CurtSenderIndex].GetState() != consts.SeatStateGameHasHu {
				break
			}
			//最大只要找3次
			if i == t.TableCfg.PlayerCt-1 {
				logs.Info("tableId:%v,*************************************playerFetchPai error!!! ", t.ID)
				return
			}
		}
	}

	t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].CancerGSH() //拿牌取消过手胡
	t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].CancerGSP() //拿牌取消过手碰

	//位置拿牌
	card := t.Majhong.FetchACard(t.Majhong.CurtSenderIndex, !moveToNext) //杠牌时取最后一张(阜阳麻将)
	t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].AddHandPai(card, true)

	//if t.state > consts.TableStateDealCard {
	//	if t.TableCfg.TableType == TableType_BBMJ && t.TableCfg.DaiHua > 0 {
	//		if card.IsBBHuaPai() {
	//			t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].RemoveHandCard(card)
	//			t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].HuaPaiArr = append(t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].HuaPaiArr, card)
	//			t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].BuHua.BuCards = append(t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].BuHua.BuCards, card.GetData())
	//			time.Sleep(time.Second * 2)
	//			//发送出牌动作
	//			t.SendTableInfo()
	//			t.playerFetchPai(false)
	//		}
	//
	//	} else if t.TableCfg.TableType == TableType_HYMJ {
	//		if card.IsHYHuaPai(t.TableCfg.FengLing) {
	//			t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].RemoveHandCard(card)
	//			t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].HuaPaiArr = append(t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].HuaPaiArr, card)
	//			t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].BuHua.BuCards = append(t.Majhong.CMajArr[t.Majhong.CurtSenderIndex].BuHua.BuCards, card.GetData())
	//			time.Sleep(time.Second * 2)
	//			//发送出牌动作
	//			t.SendTableInfo()
	//			t.playerFetchPai(false)
	//		}
	//	}
	//}

	t.ExecOptInfo = NewExecOptInfo()
	t.ExecOptInfo.OptSeatId = t.Majhong.CurtSenderIndex
	t.ExecOptInfo.OptType = OptTypeFetch
	t.ExecOptInfo.OptData = card.GetData()

}

//检测是否要结束游戏,牌墙中没有牌则结束游戏,这个方法只在自己拿完最后一张无思考或者有思考但放弃时调用
func (t *Table) chkOver() bool {
	remainPaiCt := t.Majhong.GetRemainPaiCt()
	if remainPaiCt == 0 { // 最后一张有胡,继续思考,超时则游戏结束
		logs.Info("tableId:%v------------------------>牌墙中已经没有牌,游戏结束", t.ID)
		t.GameOver()
		return true
	}
	return false
}

//游戏结束
func (t *Table) GameOver() {

	logs.Info("tableId:%v------------------------>gameOver()", t.ID)

	t.autoChangeState.Stop() //游戏结束立即停止计时器

	//第一局结束扣除房卡
	if t.Majhong.GameCt == 1 {
		isCharge := true
		if t.DirectOver { //直接解散不扣房卡
			isCharge = false
		}

		logs.Info("tableId:%v--------------------->t.Majhong.GameCt:%v  t.DirectOver:%v   isCharge:%v", t.ID, t.Majhong.GameCt, t.DirectOver, isCharge)
		//扣除房卡
		if isCharge && t.TableCfg.PayWay == 0 {
			if t.handler != nil {

				//赠送房间
				if t.TableCfg.Present == consts.Yes {

					t.GameHappend = true //标记赠送的房间已经发生游戏

				} else { //创建房间
					// for _, v := range t.GetSeats() {
					// 	logs.Info("---------------------------->tableId:%v,扣除玩家房卡:%v,数量:%v",
					// 		t.ID, v.GetPlayer(), t.TableCfg.AvgPerCard)
					// 	t.handler.ConsumeFangka(v.GetPlayer().User.ID, t.TableCfg.AvgPerCard)
					// }
					uids := make([]int, 0, 4)
					for _, v := range t.GetSeats() {
						uids = append(uids, v.GetPlayer().ID())
					}
					logs.Info("---------------------------->tableId:%v,扣除玩家房卡:%v,数量:%v",
						t.ID, uids, t.TableCfg.AvgPerCard)
					err := t.handler.MultiConsumeFangka(uids, t.ID, t.TableCfg.AvgPerCard)
					if err != nil {
						logs.Info("---------------------------->tableId:%v,扣除玩家房卡:%v,数量:%v 失败 %v",
							t.ID, uids, t.TableCfg.AvgPerCard, err)
					}
				}
			}
		}
	}

	//直接解散牌桌 不计算积分
	if t.Majhong.GameCt > 0 { //防止第一局游戏未开始就解散
		if !t.DirectOver { //直接解散也不计算
			//单局游戏结束结算
			t.tableInter.CalSingle()
		}
	}

	//游戏结束,发送牌局总结束,解散房间
	t.setState(consts.TableStateShowResult) //结算状态

	//保存牌局战绩
	t.saveZhanJi()

	if t.IsTotalGameOver() {
		//牌局全部结束,离开所有玩家

		if t.TableCfg.PayWay > 0 {
			IsCharge := true
			if t.Majhong.GameCt == 1 {
				if t.DirectOver { //直接解散不扣房卡
					IsCharge = false
				}
			}
			if IsCharge {
				NeedPerCard := t.TableCfg.FangKa
				uids := make([]int, 0)
				WinScore := make([]int, 0)
				for _, v := range t.Majhong.CMajArr {
					WinScore = append(WinScore, v.TotalScore)
				}
				MaxScore := GetMaxElement(WinScore)
				for _, v := range t.Majhong.CMajArr {
					if v.TotalScore == MaxScore {
						uids = append(uids, t.seats[v.SeatID].GetPlayer().ID())
						t.handler.MultiConsumeFangka(uids, t.ID, NeedPerCard)
						break
					}
				}
			}

		}

		t.ExitAll()
	} else {
		//结算状态,开始计时
		t.autoChangeState.Reset(time.Second * consts.GameShowResultTime)
	}

}

//后台控制解散
func (t *Table) BackstageControlGameOver() {

	logs.Info("tableId:%v------------------------>gameOver()", t.ID)

	t.autoChangeState.Stop() //游戏结束立即停止计时器

	//直接解散牌桌 不计算积分
	if t.Majhong.GameCt > 0 { //防止第一局游戏未开始就解散
		if !t.DirectOver { //直接解散也不计算
			//单局游戏结束结算
			t.tableInter.CalSingle()
		}
	}

	//游戏结束,发送牌局总结束,解散房间
	t.setState(consts.TableStateShowResult) //结算状态

	//保存牌局战绩
	t.saveZhanJi()

	if t.IsTotalGameOver() {
		//牌局全部结束,离开所有玩家
		logs.Info("牌局全部结束,离开所有玩家")
		t.ExitAll()
	} else {
		//结算状态,开始计时
		t.autoChangeState.Reset(time.Second * consts.GameShowResultTime)
	}

}

//房间游戏是否结束(达到总局数 或提前结束)
func (t *Table) IsTotalGameOver() bool {
	//判断是否结算整个牌局
	remainCt := t.TableCfg.GameCt - t.Majhong.GameCt //剩余局数
	logs.Info("tableId:%v------------IsTotalGameOver---------->>>剩余局数:%v,t.TableCfg.GameCt:%v,t.Majhong.GameCt:%v",
		t.ID, remainCt, t.TableCfg.GameCt, t.Majhong.GameCt)
	if remainCt <= 0 || t.DirectOver {
		return true
	}
	return false
}

//游戏总结算  保存战绩
func (t *Table) saveZhanJi() {

	logs.Info("tableId:%v---------------------->>>saveZhanJi", t.ID)

	scores := make([]player.Score, 0, 4)

	for _, v := range t.GetSeats() {
		score := player.Score{
			UID:      v.GetPlayer().ID(),
			NickName: v.GetPlayer().NickName(),
			Score:    t.Majhong.CMajArr[v.GetId()].Score,
		}
		scores = append(scores, score)
	}

	//保存
	t.handler.AddZhanji(t, t.Majhong.GameCt, scores)
}

//计算位置刮风下雨积分
func (t *Table) calGfxy(_gfxySeatID int) {

	//计算胡牌位置 刮风下雨所得积分
	cmaj := t.Majhong.CMajArr[_gfxySeatID]
	gfxyRec := cmaj.GfxyRec
	for _recIndex, v := range gfxyRec {
		logs.Info("tableId:%v ---------------------->计算位置:%v 刮风下雨所得积分,gfxyRec[%v]:%v", t.ID, _gfxySeatID, _recIndex, v)

		addScore := 0

		//输家
		for _seatID, _score := range v {
			if _seatID != _gfxySeatID {
				loseCmaj := t.Majhong.CMajArr[_seatID]
				if _score != 0 {
					loseCmaj.AddPxScore("杠分", -_score)
					t.Majhong.ChangeCmajScore(_seatID, -_score) //扣除积分
					addScore += _score
				}

			}
		}
		//赢家
		if addScore != 0 {
			cmaj.AddPxScore("杠分", addScore)
			t.Majhong.ChangeCmajScore(_gfxySeatID, addScore) //刮风下雨位置玩家 增加积分
		}

	}
}

//获取已胡牌玩家个数
func (t *Table) GetHasHuCt() int {
	hasHuCt := 0
	for _, v := range t.seats {
		if v.GetState() == consts.SeatStateGameHasHu {
			hasHuCt++
		}
	}
	return hasHuCt
}

//自动补花相关操作--------------------------------------------------------------------------------------
func (t *Table) buhuaLogic() {
	//if t.tableInter.ChkOver() {
	//	logs.Info("tableId:%v-------------->玩家拿最后一张别人出的牌,游戏结束!", t.ID)
	//	//t.setState(consts.TableStateIdle)
	//	return
	//}

	if t.GetState() > consts.TableStateDealCard && t.GetState() < consts.TableStateShowResult {
		lastFatchId := t.Majhong.LastFetchSeatID
		_cmaj := t.Majhong.CMajArr[lastFatchId]

		lastFatchCard := t.Majhong.LastFetchMCard
		if t.TableCfg.TableType == TableType_BBMJ && t.TableCfg.DaiHua > 0 {

			if lastFatchCard.IsBBHuaPai() {
				t.Majhong.LastSendCard = lastFatchCard.Clone() //保存这张牌
				t.Majhong.LastSenderSeatID = lastFatchId       //保存这张牌

				_cmaj.RemoveHandCard(lastFatchCard)
				_cmaj.HuaPaiArr = append(_cmaj.HuaPaiArr, lastFatchCard)
				_cmaj.BuHua.BuCards = append(_cmaj.BuHua.BuCards, lastFatchCard.GetData())

				t.ExecOptInfo = NewExecOptInfo()
				t.ExecOptInfo.OptSeatId = lastFatchId
				t.ExecOptInfo.OptType = OptTypeSend
				t.ExecOptInfo.OptData = t.Majhong.LastSendCard.GetData()

				time.Sleep(time.Second * 1)
				t.SendTableInfo()
				//发送出牌动作
				if t.Majhong.GetRemainPaiCt() > 0 {
					t.playerFetchPai(false)
					t.tableInter.FetcherStartThink(HUTYPE_DETAIL_GANG_SHANG_HUA)
					t.SendTableInfo()
				} else {
					t.GameOver()
				}
			}

		} else if t.TableCfg.TableType == TableType_HYMJ {
			if lastFatchCard.IsHYHuaPai(t.TableCfg.FengLing) {
				t.Majhong.LastSendCard = lastFatchCard.Clone() //保存这张牌
				t.Majhong.LastSenderSeatID = lastFatchId       //保存这张牌

				t.ExecOptInfo = NewExecOptInfo()
				t.ExecOptInfo.OptSeatId = lastFatchId
				t.ExecOptInfo.OptType = OptTypeSend
				t.ExecOptInfo.OptData = t.Majhong.LastSendCard.GetData()

				_cmaj.RemoveHandCard(lastFatchCard)
				_cmaj.HuaPaiArr = append(_cmaj.HuaPaiArr, lastFatchCard)
				_cmaj.BuHua.BuCards = append(_cmaj.BuHua.BuCards, lastFatchCard.GetData())
				time.Sleep(time.Second * 1)
				t.SendTableInfo()
				//发送出牌动作
				if t.Majhong.GetRemainPaiCt() > 0 {
					t.playerFetchPai(false)
					t.tableInter.FetcherStartThink(HUTYPE_DETAIL_GANG_SHANG_HUA)
					t.SendTableInfo()
				} else {
					t.GameOver()
				}
			}
		}
	}

}

//牌桌报听后相关操作--------------------------------------------------------------------------------------
func (t *Table) tingLogic() {

	if t.GetState() == consts.TableStateWaiteSend {
		for _, v := range t.seats {
			if v == nil || v.GetPlayer() == nil || !t.Majhong.CMajArr[v.Id].IsTing {
				continue
			}
			if !(t.Majhong.CMajArr[v.Id].IsTing && t.Majhong.CMajArr[v.Id].HowTingCt > 0) {
				continue
			}
			if v.GetId() == t.Majhong.CurtSenderIndex {
				time.Sleep(time.Second * 1)
				t.autoSendBySys()
				return
			}
		}
	}

}

//牌桌机器人相关--------------------------------------------------------------------------------------
func (t *Table) robotLogic() {

	//牌桌空闲状态 ,机器人进入桌子逻辑----------------------------------
	if t.GetState() == consts.TableStateIdle ||
		t.GetState() == consts.TableStateWaiteReady {

		realPlayerCt := t.GetRealPlayerCt()
		totalPlayerCt := t.GetTotalPlayerCt()
		robotCt := totalPlayerCt - realPlayerCt

		if robotCt < int(t.TableCfg.RobotCt) {
			//有真实玩家则进入机器人
			realCt := 1
			if !config.Opts().CloseTest {
				realCt = 0
			}
			if realPlayerCt >= realCt && totalPlayerCt < t.TableCfg.PlayerCt {
				for i := 0; i < t.TableCfg.RobotCt; i++ {
					t.robotEnter()
				}
				return
			}
		}
	}

	//其他状态处理------------------------------------------------------
	for _, v := range t.seats {

		if v == nil || v.GetPlayer() == nil || !v.GetPlayer().IsRobot() {
			continue
		}

		_player := v.GetPlayer()
		//牌桌空闲或等待准备,机器人自动准备
		if t.GetState() == consts.TableStateIdle || t.GetState() == consts.TableStateWaiteReady {
			if v.GetState() != consts.SeatStateGameReady {
				t.SeatReady(_player.ID())
				return
			}

		} else if t.GetState() == consts.TableStateWaiteSend {
			//当前等待机器人操作
			if v.GetId() == t.Majhong.CurtSenderIndex {
				t.autoSendBySys()
				return
			}

		} else if t.GetState() == consts.TableStateWaiteThink {
			//机器人在思考者列表中
			if t.Majhong.HasThinker(v.GetId()) {
				t.autoThinkByRobot(v)
				return
			}

		}
	}
}

// RobotEnter 机器人自动进入桌子
func (t *Table) robotEnter() {

	robot := t.robots.NewRobot()
	_robot := player.NewRobot(robot.robotID, robot.robotName)

	//初始必要游戏数据
	// random := rand.New(rand.NewSource(time.Now().UnixNano()))
	// sex := random.Intn(1) + 1 //性别
	// _robot.User.Sex = sex

	t.sendRobotEnterInfo(_robot) //机器人进入
	_idleSeatID := t.SearchIdleSeat()
	if _idleSeatID == -1 {
		return
	}

	//新建座位坐下并准备
	_seat := seat.NewSeat(_idleSeatID, _robot)
	t.sitdown(_seat)
	t.SeatReady(robot.robotID)
}

// 系统自动出牌 (玩家超时 或 机器人出牌)
func (t *Table) autoSendBySys() {
	speakerSeat := t.getSpeakerSeat()
	sendCard := t.Majhong.CMajArr[speakerSeat.GetId()].LastFetchMCard
	//这里需要优化,由于碰\杠 LastFetchCard 可能已经不在手上
	if sendCard == nil || !t.Majhong.CMajArr[speakerSeat.GetId()].IsHandHasSamePai(sendCard.GetData()) {

		length := len(t.Majhong.CMajArr[speakerSeat.GetId()].GetHandPai())
		sendCard = t.Majhong.CMajArr[speakerSeat.GetId()].GetHandPai()[length-1] //最后一张
	}

	t.SeatOpt(speakerSeat, OptTypeSend, sendCard.GetData(), false)
}

// 系统自动思考 (玩家超时 则取消思考)
func (t *Table) autoThinkBySys() {

	//系统超时,这里 可能有1~3个玩家均超时
	//1个玩家的情况:1)碰杠手动取消,胡玩家超时取消 2),胡手动取消,碰杠玩家超时取消
	//2个玩家的情况:碰杠已经操作且被保,胡手动取消

	for _, _seatID := range t.Majhong.CurtThinkerArr {
		_tkSeat := t.seats[_seatID]

		//因为该位置保存了操作,等待处理,则此处直接返回,不处理该位置,同一时刻只有一个保存操作
		if len(t.SaveOpt) > 0 && t.SaveOpt[0] == _seatID {
			continue
		} else {
			//没有保存操作的位置则放弃
			t.SeatOpt(_tkSeat, OptTypeCancel, -1, false) //取消操作
		}
	}
}

// 机器人自动思考操作
func (t *Table) autoThinkByRobot(_tkSeat *seat.Seat) {
	_seatID := _tkSeat.GetId()
	cmaj := t.Majhong.CMajArr[_seatID]

	if cmaj.CanOpt(OptTypeHu) { //可胡
		if consts.GameWaiteThinkTime-t.GetTimerRemainTime() < consts.GameWaiteRobotOptHuTime {
			return
		}
		t.SeatOpt(_tkSeat, OptTypeHu, -1, false)
	} else if cmaj.CanOpt(OptTypePeng) { //可碰
		if consts.GameWaiteThinkTime-t.GetTimerRemainTime() < consts.GameWaiteRobotOptPengTime {
			return
		}
		t.SeatOpt(_tkSeat, OptTypePeng, -1, false)
	} else if cmaj.CanOpt(OptTypeGang) { //可杠
		if consts.GameWaiteThinkTime-t.GetTimerRemainTime() < consts.GameWaiteRobotOptGangTime {
			return
		}
		_gangData := t.Majhong.CMajArr[_seatID].RobotGetGangPai()
		t.SeatOpt(_tkSeat, OptTypeGang, _gangData, false)
	} else if cmaj.CanOpt(OptTypeBu) { //可补花
		if consts.GameWaiteThinkTime-t.GetTimerRemainTime() < consts.GameWaiteRobotOptBUHUATime {
			return
		}
		t.SeatOpt(_tkSeat, OptTypeBu, -1, false)
	} else if cmaj.CanOpt(OptTypeTing) { //可听牌
		if consts.GameWaiteThinkTime-t.GetTimerRemainTime() < consts.GameWaiteRobotOptTingTime {
			return
		}
		t.SeatOpt(_tkSeat, OptTypeTing, -1, false)
	}
}

//获得桌子创建 到现在的时间 (秒)
func (t *Table) GetTableLiveTime() int {
	liveTime := int((time.Now().Unix()*1000 - t.GetCreateTime()) / 1000)
	return liveTime
}

//请求解散桌子的剩余时间
func (t *Table) GetDismissRemainTime() int {
	dismissRemainTime := int((int64(consts.GameDismissTableTime)*1000 - (time.Now().Unix()*1000 - t.ReqDismissTime)) / 1000)
	return dismissRemainTime
}

// 请求解散桌子
func (t *Table) ReqDismissTable(_player *player.Player) {

	seatID := t.GetSeats().GetSeatByUID(_player.ID()).GetId()
	t.AgreeDismissArr[seatID] = consts.DismissTableStateAgree //位置同意
	//记录请求时间
	t.ReqDismissTime = time.Now().Unix() * 1000 //转成毫秒
	t.ReqDismissSeatID = seatID
	t.autoDismissTable.Reset(time.Second * consts.GameDismissTableTime) //开始计时

	//将2个机器人设置自动同意
	count := 0
	for _, v := range t.GetSeats() {
		if v.GetPlayer().IsRobot() {
			t.AgreeDismissArr[v.GetId()] = consts.DismissTableStateAgree //位置同意
			count++
			if count == 2 {
				break
			}
		}
	}

	t.DisMissTable(seatID)
}

//获得同意解散桌子人数
func (t *Table) GetAgreeDismissCt() int {
	count := 0
	for _, v := range t.AgreeDismissArr {
		if v == consts.DismissTableStateAgree {
			count++
		}
	}
	return count
}

//获得同意解散桌子人数
func (t *Table) GetDismissInfo() []string {

	strArr := make([]string, 0)
	if t.ReqDismissSeatID == consts.DefaultIndex {
		return strArr
	}
	reqSeat := t.GetSeats().GetSeatBySeatId(t.ReqDismissSeatID)
	strTmp := "玩家[" + reqSeat.GetPlayer().NickName() + "]申请解散房间,请等待其他玩家选择(超过5分钟未做选择,则默认同意)"
	strArr = append(strArr, strTmp)

	for k, v := range t.AgreeDismissArr {
		if k != t.ReqDismissSeatID { //不包含申请人信息
			seat := t.GetSeats().GetSeatBySeatId(k)
			if seat != nil {
				if v == consts.DismissTableStateWaite { //等待
					strArr = append(strArr, "["+seat.GetPlayer().NickName()+"]等待选择")
				} else if v == consts.DismissTableStateAgree { //同意
					strArr = append(strArr, "["+seat.GetPlayer().NickName()+"]同意")
				}
			}
		}
	}

	return strArr
}

//清空 解散桌子产生的临时数据
func (t *Table) ClearDismissInfo() {

	//停止计时器
	t.autoDismissTable.Stop()

	//清空数据
	t.ReqDismissSeatID = consts.DefaultIndex
	t.AgreeDismissArr = make([]int, 4)
	for k, _ := range t.AgreeDismissArr {
		t.AgreeDismissArr[k] = consts.DismissTableStateWaite
	}
}

//检测当前是否可以退桌
func (t *Table) ChkExitTable(_player *player.Player) bool {

	//uid := _player.ID()
	if t.Majhong.GameCt == 0 { //第一局未开始,只有非房主玩家可以离桌,房主只能解散(在有其它真实玩家情况下)
		//_seat := t.GetSeats().GetSeatByUID(uid)
		//if _seat.GetId() == t.FangSeatID {
		//	if len(t.GetSeats()) > 1 {
		//		logs.Info("****************game error tableId:%v, 有其他玩家 房主不可发退桌请求!t.GetSeats():%v", t.ID, t.GetSeats())
		//		return false
		//	}
		//}

	} else { //第一局已经开始,所有人均不能离桌
		logs.Info("tableId:%v, 游戏已经开始, 所有人不可发退桌请求!", t.ID)
		return false
	}
	return true
}

// 离开所有玩家
func (t *Table) ExitAll() {

	logs.Info("tableId:%v --------------------> table.ExitAll,  离开所有玩家", t.ID)

	if len(t.seats) == 0 {
		//检测销毁 (t.GetSeats()可能为空,因此这样需要调用)
		t.ChkDestory()
	} else {
		for _, v := range t.GetSeats() {
			t.Exit(v.GetPlayer())
		}
	}

}

// Exit 玩家离开桌子(主动退出 或 被动退出)
func (t *Table) Exit(_player *player.Player) {

	logs.Info("tableId:%v, table.Exit, %v 离开桌子", t.ID, _player.String())

	_uID := _player.ID()
	_seat := t.seats.GetSeatByUID(_uID)

	t.seats.RemoveSeat(_seat.GetId()) //删除牌桌座位
	_player.Reset()                   //删除玩家临时数据

	//检测销毁
	if len(t.seats) == 0 {
		t.ChkDestory()
	}

}

// 检测桌子是否需要销毁
func (t *Table) ChkDestory() {
	if t.TableCfg.Present == consts.Yes {
		//赠送的房间不立刻删除(玩家进了又退出),由game tickTable 检测删除
		if t.GameHappend {
			//发生了游戏,立刻销毁,防止玩家再次进入
			logs.Info("tableId:%v-----------已进行了游戏----------->销毁赠送的桌子", t.ID)
			t.DestroyTable()
		} else {
			//没有发生游戏
			if t.DirectOver { //可能第一局未打完解散的
				logs.Info("tableId:%v-----------第一局未打完解散----------->销毁赠送的桌子", t.ID)
				t.DestroyTable()
			} else {
				if t.GetTableLiveTime() > consts.GameTableMaxLiveTime {
					//超过最大生存时间
					if !t.GameHappend {
						//没有发生游戏,则退回房卡
						logs.Info("tableId:%v------赠送的房间没有发生游戏,退回房卡 PresentUID:%v,fangka:%v", t.ID, t.PresentUID, t.TableCfg.FangKa)
						t.handler.ConsumeFangkaLoss(t.PresentUID, t.ConsumeID)
					}
					t.DestroyTable()
				}
			}
		}
	} else {
		logs.Info("tableId:%v---------------------->销毁创建的桌子", t.ID)
		t.DestroyTable()
	}
}

// SearchIdleSeat 查找空闲座位
func (t *Table) SearchIdleSeat() int {
	for i := 0; i < t.TableCfg.PlayerCt; i++ {
		if t.seats[i] == nil {
			return i
		}
	}
	return -1
}

//非游戏逻辑函数-----------------------------------------------------------------------------------------

//返回牌桌 已经准备的玩家数
func (t *Table) getReadyCt() int {
	count := 0
	for _, v := range t.seats {
		if v.GetState() == consts.SeatStateGameReady {
			count++
		}
	}
	return count
}

//返回牌桌 真实玩家数
func (t *Table) GetRealPlayerCt() int {
	count := 0
	for _, v := range t.seats {
		if v.GetPlayer().ID() > consts.RobotMaxUid {
			count++
		}
	}
	return count
}

//返回牌桌玩家总数
func (t *Table) GetTotalPlayerCt() int {
	if t.seats == nil {
		return 0
	}
	return len(t.seats)
}

// GetSeats 返回所有座位信息
func (t *Table) GetSeats() seat.Seats {
	return t.seats
}

// IsPlayerSit 检测玩家是否坐下
func (t *Table) IsPlayerSit(_uid int) bool {
	for _, v := range t.seats {
		if v.GetPlayer().ID() == _uid {
			return true
		}
	}
	return false
}

// 获取当前说话的位置
func (t *Table) getSpeakerSeat() *seat.Seat {
	seat := t.seats.GetSeatBySeatId(t.Majhong.CurtSenderIndex)
	if seat == nil {
		logs.Debug("tableId:%v,t.CurSpeaker:%v", t.ID, t.Majhong.CurtSenderIndex)
	}
	return seat
}

// KickPlayerInfo 踢出玩家消息
func (t *Table) kickPlayerInfo(_uID int) {
	if t.handler != nil {
		t.handler.KickPlayerInfo(t, _uID)
	}
}

// SendTableInfo 发送TableInfo
func (t *Table) SendTableInfo() {
	if t.handler != nil {
		t.handler.SendTableInfo(t)
		t.clearCltData()
	}
}

//发送TopTip
func (t *Table) SendTopTip(tipStr string) {
	if t.handler != nil {
		t.handler.SendTopTip(t, tipStr)
	}
}

func (t *Table) DisMissTable(_seatID int) {
	if t.handler != nil {
		t.handler.DisMissTable(t, _seatID)
	}
}

func (t *Table) DestroyTable() {
	if t.handler != nil {
		t.handler.DestroyTable(t)
	}
}

func (t *Table) GetHandler() Handler {
	return t.handler
}

// RobotEnterInfo 机器人
func (t *Table) sendRobotEnterInfo(_robot *player.Player) {
	if t.handler != nil {
		t.handler.RobotEnterInfo(t, _robot)
	}
}

// Destroy 销毁桌子
func (t *Table) Destroy() {
	// 桌子管理
	select {
	case <-t.finish:
	default:
		close(t.finish)
	}
}

// 添加一个定制操作
func (t *Table) addTimerAction(opt consts.TimerType, param int) {
	var delay = time.Second * 2
	// 根据不同的 opt ，延迟不同
	time.AfterFunc(delay, func() {
		t.timerAction <- TimerAction{
			timerType: opt,
			param:     param,
		}
	})
}

// 获取牌桌倒计时 剩余时间
func (t *Table) GetTimerRemainTime() int {
	timeUse := 0
	if t.state == consts.TableStateWaiteReady {
		timeUse = consts.GameWaiteReadyTime
	} else if t.state == consts.TableStateChangeCard {
		timeUse = consts.GameChangeCardTime
	} else if t.state == consts.TableStateSelectQue {
		timeUse = consts.GameSelectQueTime
	} else if t.state == consts.TableStateWaiteSend {
		timeUse = consts.GameWaiteSendTime
	} else if t.state == consts.TableStateWaiteThink {
		timeUse = consts.GameWaiteThinkTime
	} else if t.state == consts.TableStateZhuaNiao {
		timeUse = consts.GameZhuaNiaoTime
	} else if t.state == consts.TableStateShowResult {
		timeUse = consts.GameShowResultTime
	}
	time := int((int64(timeUse)*1000 - (time.Now().Unix()*1000 - t.updateTime)) / 1000)
	if time < 0 {
		time = 0
	}
	return time
}

// GetReaminSendTime 获得玩家剩余出牌时间
func (t *Table) GetReaminSendTime() int {
	return int((consts.GameWaiteSendTime*1000 - (time.Now().Unix()*1000 - t.updateTime)) / 1000)
}

// 设置牌桌状态
func (t *Table) setState(_state consts.TableState) {

	//空闲状态时清空游戏临时数据
	logs.Info("tableId:%v-----------------> 修改牌桌状态: %v ==> %v", t.ID, t.GetState(), _state)

	if _state == consts.TableStateIdle {
		if t.Majhong.GameCt >= 1 {
			t.Majhong.GameCt++ //游戏局数+1 (第1局 结算时局数+1)
		}
	} else if _state == consts.TableStateWaiteReady {
		if !config.Opts().CloseTest {
			t.autoChangeState.Reset(time.Second * consts.GameWaiteReadyTime) //等待准备,开始计时
		}
	} else if _state == consts.TableStateDealCard {
		t.autoChangeState.Reset(time.Second * consts.GameWaiteDealCardTime) //发牌状态,开始计时

	} else if _state == consts.TableStateWaiteSend {
		logs.Info("tableId:%v----------------->[等待位置%v出牌]", t.ID, t.Majhong.CurtSenderIndex)
		if !config.Opts().CloseTest {
			t.autoChangeState.Reset(time.Second * consts.GameWaiteSendTime) //等待位置出牌,开始计时
		}
	} else if _state == consts.TableStateWaiteThink {
		logs.Info("tableId:%v----------------->[等待玩家思考] t.Majhong.CurtThinkerArr:%v", t.ID, t.Majhong.CurtThinkerArr)
		if !config.Opts().CloseTest {
			t.autoChangeState.Reset(time.Second * consts.GameWaiteThinkTime) //等待思考,开始计时
		}
	} else if _state == consts.TableStateChangeCard {
		if !config.Opts().CloseTest {
			t.autoChangeState.Reset(time.Second * consts.GameChangeCardTime) //换三张状态,开始计时
		}

	} else if _state == consts.TableStateSelectQue {
		if !config.Opts().CloseTest {
			t.autoChangeState.Reset(time.Second * consts.GameSelectQueTime) //选缺状态,开始计时
		}

	} else if _state == consts.TableStateZhuaNiao {

		t.autoChangeState.Reset(time.Second * consts.GameZhuaNiaoTime) //抓鸟状态

	} else if _state == consts.TableStateShowResult {

		t.autoChangeState.Reset(time.Second * consts.GameShowResultTime)
	}

	t.updateTime = time.Now().Unix() * 1000 //转成毫秒
	t.state = _state
	t.SendTableInfo()
}

//getter and setter------------------------------------------------------------------------------

// 返回牌桌状态
func (t *Table) GetState() consts.TableState {
	return t.state
}

func (t *Table) GetCreateTime() int64 {
	return t.createTime
}
