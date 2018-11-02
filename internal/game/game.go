package game

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"golang.org/x/protobuf/proto"
	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/protoapi"
	"qianuuu.com/xgmj/internal/game/seat"
	"qianuuu.com/xgmj/internal/game/table"
	"qianuuu.com/xgmj/internal/game/timer"
	"qianuuu.com/xgmj/internal/mjcomn"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/xgmj/internal/redig"
	"qianuuu.com/xgmj/qo"
	"qianuuu.com/hall/lib/gath"
	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/util"
	"qianuuu.com/player"
)

// MsgHandler 消息处理接口
type MsgHandler interface {
	ReadMessage() *protobuf.RequestCmd
	WriteMessage(*protobuf.ResponseCmd)

	TestRequestCmd(*protobuf.RequestCmd)
}

// Game 游戏
type Game struct {
	gath *gath.Gath
	//所有玩家
	playerSvr player.Server
	//游戏桌子
	TableMap *table.Tables

	TableInfoRec *util.Map //TableInfo 数据记录

	msgHanler MsgHandler
}

// NewGame 创建一个游戏服务
func NewGame() *Game {
	ret := &Game{}

	ret.playerSvr = player.NewServer()
	ret.TableMap = table.InitTables(ret)
	ret.TableInfoRec = &util.Map{}

	return ret
}

//游戏世界循环
func (g *Game) circle() {
	// 世界任务，创建世界任务队列 (ql.Go 所调用的方法是有先后次序的)
	ql := qo.New()

	// 检测玩家掉线任务
	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				ql.Go(func() {
					// logs.Info("玩家掉线检测任务")
					g.tickPlayer()
				})
			}
		}
	}()

	//打印在线人数
	go func() {
		for {
			select {
			case <-time.After(time.Second * 60 * 5):
				ql.Go(func() {
					logs.Info("--------------------------------------------------桌子数量:%v-----------> 当前在线玩家数:%v", g.TableMap.TableCount(), g.playerSvr.PlayerCount())
				})
			}
		}
	}()

	// 删除 牌桌任务
	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				ql.Go(func() {
					g.tickTable()
				})
			}
		}
	}()
}

//牌桌状态检测 (60分钟内,牌桌 1~4 人,且4人都已掉线,则强制删除牌桌)
func (g *Game) tickTable() {
	if g.playerSvr == nil {
		return
	}

	tableIds := []int{}

	g.TableMap.ReadRange(func(tid int, table *table.Table) {

		//桌子创建时间超过1小时
		if table.GetTableLiveTime() > consts.GameTableMaxLiveTime {
			deadLineCt := 0
			for _, n := range table.GetSeats() {
				if g.playerSvr.IsTimeout(n.GetPlayer().ID(), consts.GameMaxSavePlayerDataTime) {
					deadLineCt++
				}
			}
			logs.Info("-------------------------------tickTable, tableId:%v 桌子创建时间超过1小时,离线玩家数:%v,真实玩家总数:%v",
				tid, deadLineCt, table.GetRealPlayerCt())
			//所有人都超时离线
			if deadLineCt == table.GetRealPlayerCt() {
				tableIds = append(tableIds, tid)
			} else {
				logs.Info("warn-------------------------------tickTable, tableId:%v 桌子创建时间超过1小时,离线玩家数:%v,真实玩家总数:%v",
					tid, deadLineCt, table.GetRealPlayerCt())
			}
		}
	})

	// 在读循环之外执行登出操作
	if len(tableIds) > 0 {
		logs.Custom(logs.AnalyseTag, " detroy table  tableIds:%v", tableIds)
		for _, _talbeId := range tableIds {
			_table := g.TableMap.GetTable(_talbeId)
			_table.Go(func() {
				_table.ExitAll()
			})
		}
	}
}

//玩家状态检测(掉线检测)
func (g *Game) tickPlayer() {
	if g.playerSvr == nil {
		return
	}

	if !config.Opts().OpenHeart {
		return
	}

	// 记录所有可能超时的 uid
	uids := []int{}
	g.playerSvr.ReadRange(func(_uid int, _player *player.Player) {
		if !_player.IsRobot() {

			if g.playerSvr.IsTimeout(_uid, consts.GameCheckOffLineTime) {
				if !_player.IsOffline() {
					logs.Info("玩家短暂离线(30秒) Player:%v,_player.TableId:%v", _player, _player.GetTableID())
					_player.SetOffline(true)
					if _player.GetTableID() > 0 {
						_table := g.TableMap.GetTable(_player.GetTableID())
						_table.Go(func() {
							_table.SendTableInfo()
						})
					}
				}
			}

			if g.playerSvr.IsTimeout(_uid, consts.GameMaxSavePlayerDataTime) {
				logs.Info("系统检测到玩家掉线 TableId:%v , Player:%v,_player ", _player.GetTableID(), _player)
				exitSucc := g.TableMap.ExitTable(_player)
				if exitSucc {
					uids = append(uids, _uid)
				}
			}
		}
	})
	// 在读循环之外执行登出操作
	if len(uids) > 0 {
		logs.Info("-----------------------> Logout player uids:%v", uids)
		for _, uid := range uids {
			g.playerSvr.Logout(uid)
		}
	}
}

func (g *Game) handlerRequestCmd() {
	protoapi.Init(func(pcmd *protobuf.ResponseCmd) {
		//logs.Info("protoapi response %v", pcmd)
		g.msgHanler.WriteMessage(pcmd)
	})
	// 请求消息处理
	for {
		cmd := g.msgHanler.ReadMessage()

		// 获取当前的玩家
		uid := int(cmd.Head.GetUid())
		player := g.playerSvr.GetPlayer(uid)

		// 防止消息读取阻塞，使用 go (如果 player 不存在使用 go 异步，如果存在，使用 player 顺序调用)
		player.Go(func() {
			// 捕获 API 调用异常
			defer func() {
				if err := recover(); err != nil {
					trace := make([]byte, 1<<16)
					n := runtime.Stack(trace, true)
					logs.Info("%v", fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
						err, n, trace[:n]))
				}
			}()

			if cmd.Simple != nil {

				if cmd.Simple.GetTag() != protobuf.ReqSimple_HEART_BEAT {
					logs.Info("read test: %v", cmd.Simple.GetTag())
				}

				switch cmd.Simple.GetTag() {

				case protobuf.ReqSimple_LOGIN_GAME: // 登陆游戏
					protoapi.LoginGame(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_HEART_BEAT: // 游戏心跳
					protoapi.HeartBeat(cmd, g.playerSvr)
					break
				case protobuf.ReqSimple_EXIT_GAME: // 退出游戏
					protoapi.ExitGame(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_PLAYER_INFO: // 请求玩家信息
					protoapi.GetPlayerInfo(cmd, g.playerSvr)
					break
				case protobuf.ReqSimple_ENTER_TABLE: // 进入桌子
					protoapi.EnterTable(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_EXIT_TABLE: // 离开桌子
					protoapi.ExitTable(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_READY: // 点击准备
					protoapi.PlayerReady(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_PLAYER_OPT: // 游戏操作
					protoapi.PlayerOpt(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_TABLE_CHAT: // 牌桌聊天
					protoapi.TableChat(cmd, g.TableMap, g.playerSvr)
					break
					//case protobuf.ReqSimple_WORLD_CHAT: // 世界聊天
					//	protoapi.WorldChat(cmd, g.TableMap, g.playerSvr)
					//	break
					//case protobuf.ReqSimple_GAME_CHARGE: // 充值
					//	protoapi.Charge(cmd, g.TableMap, g.playerSvr)
					//	break
					//case protobuf.ReqSimple_GAME_BUY: // 购买金币
					//	protoapi.BuyCoin(cmd, g.TableMap, g.playerSvr)
					//	break
					//case protobuf.ReqSimple_GET_RANK_LIST: //请求排行
					//	protoapi.GetRankList(cmd, g.TableMap, g.playerSvr)
					//	break
				case protobuf.ReqSimple_SELECT_QUE: //选缺
					//protoapi.SelectQue(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_REQ_DISMISS: //请求解散房间
					protoapi.ReqDismiss(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_RESP_DISMISS: //响应解散房间
					protoapi.RespDismiss(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_OFFLINE: //玩家断线
					protoapi.OffLine(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_REQ_STOP: //程序后台运行
					protoapi.ReqStop(cmd, g.TableMap, g.playerSvr)
					break
				case protobuf.ReqSimple_TABLEINFO_REC: //牌局回放
					protoapi.TableInfoRec(cmd, g.playerSvr)
					break
					//case protobuf.ReqSimple_FANGKA_LIST: //房卡数据列表
					//	protoapi.FangKaList(cmd, g.playerSvr)
					//	break

				}
			}

			if cmd.CrateTable != nil { //创建桌子
				protoapi.CreateTable(cmd, g.TableMap, g.playerSvr)
			}

			//更新接收消息时间
			_player := g.playerSvr.GetPlayer(int(cmd.Head.GetUid()))
			if _player != nil {
				_player.UpdateMsgRecevTime()
				check := true
				if cmd.Simple != nil {
					if cmd.Simple.GetTag() == protobuf.ReqSimple_REQ_STOP {
						check = false // 去除该请求
					}
				}
				if check && _player.IsOffline() { //离线返回
					_player.SetOffline(false)
					if _player.GetTableID() > 0 {
						_table := g.TableMap.GetTable(_player.GetTableID())
						if _table != nil {
							_table.Go(func() {
								_table.SendTableInfo()
							})

						}
					}
				}
			}
		})
	}
}

//销毁桌子回调
func (g *Game) DestroyTable(_table *table.Table) {

	logs.Info("game DestroyTable  _tableId:%v", _table.ID)

	//删除牌桌
	g.TableMap.RemoveTable(_table.ID)

	//如果有牌局记录,则删除
	g.RemoveTableInfoRec(_table.ID)

}

//解散桌子消息
func (g *Game) DisMissTable(_table *table.Table, _seatID int) {

	dismissRemainTime := _table.GetDismissRemainTime()
	logs.Info("game DisMissTable  _tableId:%v ,dismissRemainTime:%v ,_table.GetAgreeDismissCt():%v ,_seatID:%v,seats:%v",
		_table.ID, dismissRemainTime, _table.GetAgreeDismissCt(), _seatID, _table.GetSeats())

	////先判断是否是房主在第一局未开始时的计算
	//if _seatID >= 0 && _seatID < _table.TableCfg.PlayerCt {
	//	if _seatID == _table.FangSeatID { //房主
	//		if _table.Majhong.GameCt == 0 && _table.GetState() <= consts.TableStateWaiteReady {
	//
	//			strTmp := "房主[" + _table.GetSeats().GetSeatBySeatId(_seatID).GetPlayer().NickName() + "]解散了房间!"
	//			g.DoDisMiss(_table, strTmp)
	//			return
	//		}
	//	}
	//}

	//MaxPlayerCt-1 个人同意 或 超时解散房间
	doDismiss := false
	if _table.TableCfg.PlayerCt > 2 {
		if _table.GetAgreeDismissCt() >= _table.TableCfg.PlayerCt-1 {
			doDismiss = true
		}
	}
	if _table.TableCfg.PlayerCt == 2 { //2个人时,必须都同意
		if _table.GetAgreeDismissCt() == 2 {
			doDismiss = true
		}
	}
	if doDismiss || dismissRemainTime <= 0 {

		strTmp := ""
		if dismissRemainTime <= 0 {
			strTmp = "时间到,解散房间成功!"
		} else {
			nameArr := make([]string, 0)
			for _, v := range _table.GetSeats() {
				if _table.AgreeDismissArr[v.GetId()] == consts.DismissTableStateAgree {
					nameArr = append(nameArr, v.GetPlayer().NickName())
				}
			}
			if len(nameArr) >= _table.TableCfg.PlayerCt-1 {
				if len(nameArr) == 1 {
					strTmp = "经玩家[" + nameArr[0] + "]同意,房间解散成功"
				} else if len(nameArr) == 2 {
					strTmp = "经玩家[" + nameArr[0] + "],[" + nameArr[1] + "]同意,房间解散成功"
				} else if len(nameArr) == 3 {
					strTmp = "经玩家[" + nameArr[0] + "],[" + nameArr[1] + "],[" + nameArr[2] + "]同意,房间解散成功"
				}

			}
		}
		g.DoDisMiss(_table, strTmp, false)

	} else {
		//申请人位置,及各个位置当前是否同意
		intArr := make([]int32, 5)
		intArr[4] = int32(_table.ReqDismissSeatID)
		for _, v := range _table.GetSeats() {
			intArr[v.GetId()] = int32(_table.AgreeDismissArr[v.GetId()])
		}

		//刷新当前解散牌桌界面信息
		for _, v := range _table.GetSeats() {
			if !v.GetPlayer().IsRobot() {
				pcmd := &protobuf.ResponseCmd{
					Head: &protobuf.RespHead{
						Uid:    proto.Int32(int32(v.GetPlayer().ID())),
						MsgID:  proto.Int32(0),
						Result: proto.Int32(0),
					},
					Simple: &protobuf.RespSimple{
						Tag:      protobuf.RespSimple_DISMISS_TABLE.Enum(),
						IntValue: proto.Int32(int32(_table.GetDismissRemainTime())),
						StrArr:   _table.GetDismissInfo(),
						IntArr:   intArr,
					},
				}
				g.msgHanler.WriteMessage(pcmd)
			}
		}
	}
}

//执行解散
func (g *Game) DoDisMiss(_table *table.Table, _tipStr string, flag bool) {

	_table.ClearDismissInfo() //清空 解散数据

	//发送同意解散房间消息
	for _, v := range _table.GetSeats() {
		if !v.GetPlayer().IsRobot() {
			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(v.GetPlayer().ID())),
					MsgID:  proto.Int32(0),
					Result: proto.Int32(0),
				},
				Simple: &protobuf.RespSimple{
					Tag:      protobuf.RespSimple_DISMISS_RESULT.Enum(),
					IntValue: proto.Int32(int32(consts.Yes)),
					StrValue: proto.String(_tipStr),
				},
			}
			g.msgHanler.WriteMessage(pcmd)
		}
	}

	if flag {
		//销毁桌子
		logs.Info("tableId:%v, 强制解散牌桌!", _table.TableCfg.TableId)

		if _table.Majhong.GameCt == 0 && _table.GetState() <= consts.TableStateWaiteReady {
			//销毁桌子
			logs.Info("tableId:%v, 第一局游戏还未开始房主解散了房间!", _table.ID)
			_table.ExitAll()
			return
		} else {
			_table.DirectOver = true          //游戏直接结束
			_table.BackstageControlGameOver() //后台控制解散
			//	_table.ExitAll()
		}

		return
	}

	//如果是第一局游戏还未开始房主解散了房间
	if _table.Majhong.GameCt == 0 && _table.GetState() <= consts.TableStateWaiteReady {
		//销毁桌子
		logs.Info("tableId:%v, 第一局游戏还未开始房主解散了房间!", _table.ID)
		_table.ExitAll()
		return
	}
	_table.DirectOver = true //游戏直接结束
	_table.GameOver()

}

// AddZhanji 添加战绩
func (g *Game) AddZhanji(_table *table.Table, _seq int, score []player.Score) {
	//logs.Info("战绩---- > :%#v", g.TableInfoRec.Get(_table.ID))
	_tableId := _table.ID
	if _seq != 0 && g.TableInfoRec.Get(_tableId) == nil {
		logs.Info("*******************AddZhanji error!!! tableId:%v, TableInfoRec not found", _tableId)
		return
	}

	tablInfoArr := make([]protobuf.TableInfo, 0)
	if _seq != 0 {
		tinforec := g.TableInfoRec.Get(_tableId).(*protobuf.TableInfoRec)
		tablInfoArr = tinforec.TableInfoArr
	}
	// 玩家 id 信息构建
	uids := make([]int, 0, 4)
	for _, v := range _table.GetSeats() {
		uids = append(uids, v.GetPlayer().ID())
	}
	error := g.AddTableRecord(_tableId, _seq, tablInfoArr, score)
	if error != nil {
		logs.Info("tableId:%v---------------AddTableRecord-->error:%v", _tableId, error)
	}
	g.RemoveTableInfoRec(_tableId)

}

// AddTableRecord 添加牌桌记录 (桌子 id， 牌局， tableinfo 序列)
func (g *Game) AddTableRecord(tableid, inning int, tableinfo []protobuf.TableInfo, score []player.Score) error {
	//将结构序列化成JSON
	ret, err := json.Marshal(tableinfo)
	if err != nil {
		return err
	}

	_, err = g.playerSvr.AddTableRecord(tableid, inning, score, ret)
	return err
}

//踢出玩家消息
func (g *Game) KickPlayerInfo(_table *table.Table, _uID int) {

	pcmd := &protobuf.ResponseCmd{
		Head: &protobuf.RespHead{
			Uid:    proto.Int32(int32(_uID)),
			MsgID:  proto.Int32(int32(0)),
			Result: proto.Int32(int32(0)),
		},
		Simple: &protobuf.RespSimple{
			Tag: protobuf.RespSimple_KICK_OUT.Enum(),
		},
	}
	g.msgHanler.WriteMessage(pcmd)

	pcmd = &protobuf.ResponseCmd{
		Head: &protobuf.RespHead{
			Uid:    proto.Int32(int32(_uID)),
			MsgID:  proto.Int32(int32(0)),
			Result: proto.Int32(int32(0)),
		},
		Simple: &protobuf.RespSimple{
			Tag:      protobuf.RespSimple_TOP_TIP.Enum(),
			StrValue: proto.String("由于你长时间未准备,被系统踢出桌子"),
		},
	}
	g.msgHanler.WriteMessage(pcmd)

}

// 桌子机器人进入回调消息
func (g *Game) RobotEnterInfo(_table *table.Table, _player *player.Player) {
	// _user := &domain.User{
	// 	ID:       _player.ID(),
	// 	NickName: _player.NickName(),
	// }
	// _ = _user
}

//牌桌中主动推送回调消息
func (g *Game) SendTableInfo(_table *table.Table) {

	_nameArr := make([]string, 0)

	_table.GetSeats().Foreach(func(key int, seat *seat.Seat) {

		player := seat.GetPlayer()
		_nameArr = append(_nameArr, player.String())
		if !player.IsRobot() { //机器人消息 不发送
			_tableInfo := protobuf.Helper.GetTableInfo(_table, player.ID())

			for i := 0; i < len(_tableInfo.SeatInfo); i++ { //不发送其他玩家的手牌
				if int(_tableInfo.SeatInfo[i].PlayerInfo.GetUid()) != player.ID() {
					handcardLen := len(_tableInfo.SeatInfo[i].GetHandCards())
					if (handcardLen > 0 && _table.GetState() != consts.TableStateShowResult) {
						//logs.Info("修改前的手牌  id:%v   手牌:%v", i, _tableInfo.SeatInfo[i].GetHandCards())
						handCard := _tableInfo.SeatInfo[i].GetHandCards()
						for i := 0; i < handcardLen; i++ {
							handCard[i] = int32(0 + i)
						}
						_tableInfo.SeatInfo[i].HandCards = handCard
						//logs.Info("修改后的手牌  id:%v   手牌:%v", i, _tableInfo.SeatInfo[i].GetHandCards())
					}

				}
			}

			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(player.ID())),
					MsgID:  proto.Int32(int32(0)),
					Result: proto.Int32(int32(0)),
				},
				TableInfo: _tableInfo,
			}
			g.msgHanler.WriteMessage(pcmd)

			if seat.GetId() == 0 { //只记录一次 tableinfo ==> 只记录位置0的

				gameCt := _table.Majhong.GameCt
				if gameCt == 0 {
					gameCt = 1
				}
				if g.TableInfoRec.Get(_table.ID) == nil { //新建一个记录
					tinforec := protobuf.NewTableInfoRec(_table.ID, gameCt)
					g.TableInfoRec.Set(_table.ID, tinforec)
				}
				_tableInfoSave := protobuf.Helper.GetTableInfo(_table, player.ID())

				tinforec := g.TableInfoRec.Get(_table.ID).(*protobuf.TableInfoRec)
				//tinforec.AddInfoRec(*_tableInfo)
				tinforec.AddInfoRec(*_tableInfoSave)
			}
		}

	})

	logs.Info("----------------->SendTableInfo, tableId:%v, _nameArr:%v", _table.ID, _nameArr)
}

//删除保存的单局回放数据
func (g *Game) RemoveTableInfoRec(_tableId int) {
	if g.TableInfoRec.Get(_tableId) == nil {
		return
	}
	//删除数据
	len1 := g.TableInfoRec.Len()
	g.TableInfoRec.Del(_tableId)
	len2 := g.TableInfoRec.Len()
	logs.Info("----------------->g.TableInfoRec  Len1():%v , Len2():%v", len1, len2)
}

//牌桌中主动推送回调消息
func (g *Game) SendTopTip(_table *table.Table, tipStr string) {

	_playerArr := protobuf.Helper.GetTablePlayerArr(_table)

	for _, player := range _playerArr {
		if !player.IsRobot() {
			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(player.ID())),
					MsgID:  proto.Int32(int32(0)),
					Result: proto.Int32(int32(0)),
				},
				Simple: &protobuf.RespSimple{
					Tag:      protobuf.RespSimple_TOP_TIP.Enum(),
					StrValue: proto.String(tipStr),
				},
			}
			g.msgHanler.WriteMessage(pcmd)
		}
	}
}

// ConsumeFangkaLoss 返还房卡
func (g *Game) ConsumeFangkaLoss(_uid int, consumeid int) {
	g.playerSvr.ConsumeFangkaLoss(_uid, consumeid)
	_player := g.playerSvr.GetPlayer(_uid)

	if _player != nil {
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(_uid)),
				MsgID:  proto.Int32(0),
				Result: proto.Int32(0),
			},
			PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
		}
		g.msgHanler.WriteMessage(pcmd)
	}

}

// MultiConsumeFangka 多人消耗房卡
func (g *Game) MultiConsumeFangka(uids []int, tid, count int) error {
	_, err := g.playerSvr.MultiConsumeGoods(uids, tid, count)
	if err != nil {
		return err
	}
	for _, _uid := range uids {
		_player := g.playerSvr.GetPlayer(_uid)
		if _player != nil {
			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(_uid)),
					MsgID:  proto.Int32(0),
					Result: proto.Int32(0),
				},
				PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
			}
			g.msgHanler.WriteMessage(pcmd)
		}
	}
	return nil
}

// ConsumeFangka 消耗房卡
func (g *Game) ConsumeFangka(_uid int, _tid int, count int) (int, error) {

	ret, err := g.playerSvr.ConsumeFangka(_uid, _tid, count)
	if err == nil && ret > 0 {
		_player := g.playerSvr.GetPlayer(_uid)
		if _player != nil {
			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(_uid)),
					MsgID:  proto.Int32(0),
					Result: proto.Int32(0),
				},
				PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
			}
			g.msgHanler.WriteMessage(pcmd)
		}
	} else {
		logs.Error("消耗房卡失败 %v %v %v", _uid, _tid, count)
	}
	return ret, err
}

// Serve 游戏服务
func (g *Game) Serve() {
	// 首先启动服务器数据库链接
	cfg := config.Opts()

	_ = g.playerSvr.Init(config.Opts().RPCUrl, config.Opts().AppName)

	server := redig.NewGate()
	server.AddHallAddr(cfg.HallRdds, "_list:hall")
	g.msgHanler = server
	go server.Serve(cfg.ReceiveRdds, "_list:"+cfg.AppName)

	// 注意启动顺序 (goroutine 可能是无顺序的)
	go g.handlerRequestCmd()
	go g.circle()

	// 接受游戏命令
	gh := gath.New(cfg.AppName, cfg.HallRdds)
	g.gath = gh
	//logs.Info("接受游戏命令:%v ---------------", g.gath)
	go func() {
		gh.ListeningGameCommand(&GameCmdHandler{g: g})
	}()

	// 服务启动完毕之后再测试
	if !config.Opts().CloseTest {
		//go g.testMessage()
		go g.testLogic()
	}

	//游戏定时器启动
	timer.Serve(g.playerSvr)

}

//游戏逻辑测试
func (g *Game) testLogic() {

	logs.Info("------------------------------------------>testLogic()")

	////[合肥麻将]牌桌
	//_tableCfg := config.NewTableCfg()
	//_tableCfg.TableType = mjcomn.TableType_HFMJ
	//_tableCfg.RobotCt = 4
	//_tableCfg.PlayerCt = 4
	//_tableCfg.BaseScore = 1
	//_tableCfg.DianpaoHu = consts.Yes
	//_tableCfg.TiandiHu = consts.Yes
	//g.TableMap.CreateHFTable(_tableCfg)

	//[红中麻将]牌桌
	//_tableCfg := config.NewTableCfg()
	//_tableCfg.TableType = mjcomn.TableType_HZMJ
	//_tableCfg.KehuQidui = consts.Yes
	//_tableCfg.GameCt = 8
	//_tableCfg.RobotCt = 4
	//_tableCfg.PlayerCt = 4
	//
	//tid, err := g.playerSvr.GetTableID()
	//if err != nil || tid <= 0 {
	//}
	//_tableCfg.TableId = tid
	//g.TableMap.CreateHZTable(_tableCfg)

	//[怀远麻将]牌桌
	_tableCfg := config.NewTableCfg()
	_tableCfg.TableType = mjcomn.TableType_HYMJ
	_tableCfg.RobotCt = 2
	_tableCfg.MaxCardColorIndex = 6
	_tableCfg.PlayerCt = 2
	_tableCfg.GameCt = 1
	_tableCfg.BaseScore = 1
	_tableCfg.WuHuaGuo = 0
	_tableCfg.FengLing = consts.Yes
	_tableCfg.DianpaoHu = consts.Yes
	_tableCfg.KehuQidui = consts.Yes
	_tableCfg.TiandiHu = consts.Yes
	g.TableMap.CreateHYTable(_tableCfg)

	////[蚌埠麻将]牌桌
	//	_tableCfg := config.NewTableCfg()
	//	_tableCfg.TableType = mjcomn.TableType_BBMJ
	//	_tableCfg.RobotCt = 4
	//	_tableCfg.PlayerCt = 1
	//	_tableCfg.BaseScore = 1
	//	_tableCfg.DaiHua = consts.Yes
	//	_tableCfg.DianpaoHu = consts.Yes
	//	_tableCfg.TiandiHu = consts.Yes
	//	g.TableMap.CreateBBTable(_tableCfg)

}

//消息测试
func (g *Game) testMessage() {
	// 模拟发送消息
	logs.Info("write msg ..")

	//_testUid := int(100509)
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestLoginGame(_testUid)) //登陆游戏
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestHeartBeat(_testUid))        //游戏心跳
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestGetPlayerInfo(_testUid))    //获取个人信息
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestExitGame(_testUid))        //退出游戏
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestRoomList(_testUid))        //房间列表
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestCreateTable(_testUid)) //创建桌子
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestEnterTable(_testUid))      //进入桌子
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestExitTable(_testUid))       //退桌
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestChangeTable(_testUid))     //换桌
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestReady(_testUid)) //准备
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestQuickStart(_testUid))      //快速开始
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestTableChat(_testUid))       //牌桌聊天
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestWorldChat(_testUid)) //世界聊天
	//g.msgHanler.TestRequestCmd(protobuf.Helper.SystemWorldChat()) //世界聊天
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestGameCharge(_testUid))      //充值
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestGameBuy(_testUid))         //金币购买
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestWatchPlayerInfo(_testUid))    //查看玩家信息
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestGetRankList(_testUid))      //请求排行
	//g.msgHanler.TestRequestCmd(protobuf.Helper.TestZhanJiList(_testUid)) //请求战绩

	logs.Info("test end ...")
}
