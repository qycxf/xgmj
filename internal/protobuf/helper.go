//
// Author: leafsoar
// Date: 2016-05-12 17:55:33
//

package protobuf

import (
	"time"

	"golang.org/x/protobuf/proto"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/seat"
	"qianuuu.com/xgmj/internal/game/table"
	"qianuuu.com/xgmj/internal/mjcomn"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

// Helper Protobuf 帮助
var Helper = &helper{}

// MsgHelper 消息帮助类
type helper struct {
}

//------------------------------------------------------------> TEST

// 请求排行
func (h *helper) TestGetRankList(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_GET_RANK_LIST.Enum(),
			IntValue: proto.Int32(1), //1:金币排行 2:今日赢取
		},
	}
	return cmd
}

// 查看个人信息
func (h *helper) TestWatchPlayerInfo(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_WATCH_PLAYER_INFO.Enum(),
			IntValue: proto.Int32(100),
		},
	}
	return cmd
}

// 购买
func (h *helper) TestGameBuy(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_GAME_BUY.Enum(),
			IntValue: proto.Int32(2001),
		},
	}
	return cmd
}

// 充值
func (h *helper) TestGameCharge(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_GAME_CHARGE.Enum(),
			IntValue: proto.Int32(1002),
		},
	}
	return cmd
}

// 世界聊天 - 个人
func (h *helper) TestWorldChat(uid int) *RequestCmd {
	content := "你好,这是世界聊天测试文字"
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_WORLD_CHAT.Enum(),
			StrValue: &content,
		},
	}
	return cmd
}

// 世界聊天 - 系统
func (h *helper) SystemWorldChat(title, content string) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(0)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_WORLD_CHAT.Enum(),
			StrValue: &content,
		},
	}
	return cmd
}

// 牌桌聊天
func (h *helper) TestTableChat(uid int) *RequestCmd {
	content := "你好,这是测试文字"
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_TABLE_CHAT.Enum(),
			StrValue: &content,
		},
	}
	return cmd
}

// 游戏心跳
func (h *helper) TestHeartBeat(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag: ReqSimple_HEART_BEAT.Enum(),
		},
	}
	return cmd
}

// 登陆游戏
func (h *helper) TestLoginGame(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag: ReqSimple_LOGIN_GAME.Enum(),
		},
	}
	return cmd
}

// 获取个人信息
func (h *helper) TestGetPlayerInfo(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag: ReqSimple_PLAYER_INFO.Enum(),
		},
	}
	return cmd
}

// 退出游戏
func (h *helper) TestExitGame(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag: ReqSimple_EXIT_GAME.Enum(),
		},
	}
	return cmd
}

// 请求房间列表
func (h *helper) TestRoomList(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag: ReqSimple_ROOT_LIST.Enum(),
		},
	}
	return cmd
}

////创建桌子消息
//func (h *helper) TestCreateTable(uid int) *RequestCmd {
//	cmd := &RequestCmd{
//		Head: &ReqHead{
//			Uid:   proto.Int32(int32(uid)),
//			MsgID: proto.Int32(1),
//		},
//		CrateTable: &ReqCreateTalble{
//			TableType: proto.Int32(TableType_HFMJ),
//			RobotCt:   proto.Int32(3),
//		},
//	}
//	return cmd
//}

//进入桌子消息
func (h *helper) TestEnterTable(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_ENTER_TABLE.Enum(), //进入桌子
			IntValue: proto.Int32(1),
		},
	}
	return cmd
}

//退桌
func (h *helper) TestExitTable(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_EXIT_TABLE.Enum(),
			IntValue: proto.Int32(1),
		},
	}
	return cmd
}

//换桌
func (h *helper) TestChangeTable(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_CHANGE_TABLE.Enum(),
			IntValue: proto.Int32(1),
		},
	}
	return cmd
}

// 准备消息
func (h *helper) TestReady(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_READY.Enum(), //准备
			IntValue: proto.Int32(1),
		},
	}
	return cmd
}

//快速开始
func (h *helper) TestQuickStart(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag:      ReqSimple_QUICK_START.Enum(),
			IntValue: proto.Int32(1),
		},
	}
	return cmd
}

func (h *helper) TestZhanJiList(uid int) *RequestCmd {
	cmd := &RequestCmd{
		Head: &ReqHead{
			Uid:   proto.Int32(int32(uid)),
			MsgID: proto.Int32(1),
		},
		Simple: &ReqSimple{
			Tag: ReqSimple_ZHANJI_LIST.Enum(),
		},
	}
	return cmd
}

//----------------------------------------------------------------------------------------> TEST OVER

// GetTablePlayerArr 获取牌桌所有玩家
func (h *helper) GetTablePlayerArr(_table *table.Table) []*player.Player {
	seats := _table.GetSeats()
	_playerArr := make([]*player.Player, 0)

	//logs.Info("------------------------------->len(seats):%v", len(seats))
	for _, v := range seats {
		_playerArr = append(_playerArr, v.GetPlayer())
	}
	return _playerArr
}

//返回 protobuf.TableInfo 对象
func (h *helper) GetTableInfo(_table *table.Table, _senduid int) *TableInfo {

	//获取每个位置信息
	_seatInfo := make([]*SeatInfo, 0)
	_sendSeat := _table.GetSeats().GetSeatByUID(_senduid)
	seats := _table.GetSeats()

	seats.Foreach(func(key int, seat *seat.Seat) {
		_sinfo := h.GetSeatInfo(_table, seat, _sendSeat)
		_seatInfo = append(_seatInfo, _sinfo)
	})

	//位置操作动作信息(客户端使用)
	execOpt := &ExecOpt{}
	if _table.ExecOptInfo != nil {
		_hupxIdArr := make([]int32, 0)
		for _, v := range _table.ExecOptInfo.HupxIdArr {
			_hupxIdArr = append(_hupxIdArr, int32(v))
		}

		execOpt = &ExecOpt{
			OptSeatId:     proto.Int32(int32(_table.ExecOptInfo.OptSeatId)),     //操作位置
			OptType:       proto.Int32(int32(_table.ExecOptInfo.OptType)),       //操作类型
			OptData:       proto.Int32(int32(_table.ExecOptInfo.OptData)),       //操作对应的牌值
			OptDetail:     proto.Int32(int32(_table.ExecOptInfo.OptDetail)),     //操作详细类型
			DianPaoSeatId: proto.Int32(int32(_table.ExecOptInfo.DianPaoSeatId)), //点炮位置
			HupxIdArr:     _hupxIdArr,
		}
	}

	_curtSpeaker := 0                        //当前等待出牌者
	_lastSendCardData := consts.DefaultIndex //最近一个玩家打出的牌
	if _table.GetState() >= consts.TableStateWaiteSend && _table.GetState() <= consts.TableStateShowResult {
		_curtSpeaker = _table.Majhong.CurtSenderIndex
		if _table.Majhong.LastSendCard != nil {
			_lastSendCardData = _table.Majhong.LastSendCard.GetData()
		}
	}

	//等待位置出牌,发送出牌提示信息  -----------------------------------------------------
	sendTipArr := make([]*SendTip, 0)
	if _table.GetState() == consts.TableStateWaiteSend {
		senderSeatId := _table.Majhong.CurtSenderIndex
		senderCmaj := _table.Majhong.CMajArr[senderSeatId]
		logs.Info("=======================>len(senderCmaj.SendTipArr):%v", len(senderCmaj.SendTipArr))
		for _, v := range senderCmaj.SendTipArr {

			huCards := make([]int32, 0)
			for _, n := range v.HuCards {
				huCards = append(huCards, int32(n))
			}

			sendTip := &SendTip{
				SendCard: proto.Int32(int32(v.SendCard)),
				HuCards:  huCards,
			}
			sendTipArr = append(sendTipArr, sendTip)
		}
	}

	_lianZhuangCt := 0 //合肥麻将连庄数
	if _table.Majhong.DSeatID != consts.DefaultIndex {
		//_lianZhuangCt = _table.Majhong.CMajArr[_table.Majhong.DSeatID].LianZhuangCt
	}

	//抓鸟数据
	_zhuaNiaoInfo := &ZhuaNiaoInfo{
		NiaoCardArr:  make([]int32, 0),
		ZhongNiaoArr: make([]int32, 0),
	}

	if _table.GetState() >= consts.TableStateZhuaNiao {
		znInfo := _table.Majhong.ZhuaNiaoInfo
		for i := 0; i < len(znInfo.NiaoCardArr); i++ {
			_zhuaNiaoInfo.NiaoCardArr = append(_zhuaNiaoInfo.NiaoCardArr, int32(znInfo.NiaoCardArr[i]))
		}
		for i := 0; i < len(znInfo.ZhongNiaoArr); i++ {
			_zhuaNiaoInfo.ZhongNiaoArr = append(_zhuaNiaoInfo.ZhongNiaoArr, int32(znInfo.ZhongNiaoArr[i]))
		}
	}

	//结算信息
	_seatResult := make([]*SeatResult, 0)
	_gameResult := make([]*SeatResult, 0)
	_totalResult := make([]*GameResult, 0)
	if _table.GetState() == consts.TableStateShowResult {

		//计算最大赢家和最佳炮手
		winScArr := make([]int, 4)
		PaoArr := make([]int, 4)
		for _, v := range _table.GetSeats() {
			winScArr[v.GetId()] = _table.Majhong.CMajArr[v.GetId()].TotalScore
			PaoArr[v.GetId()] = _table.Majhong.CMajArr[v.GetId()].DianPaoCt
		}

		maxWinner := make([]int, 0)
		maxWinSc := mjcomn.GetMaxElement(winScArr)
		if maxWinSc > 0 {
			for k, v := range winScArr {
				if v == maxWinSc {
					maxWinner = append(maxWinner, int(k))
					break //只取第一个
				}
			}
		}

		maxPaoshou := make([]int, 0)
		maxPao := mjcomn.GetMaxElement(PaoArr)
		if maxPao > 0 {
			for k, v := range PaoArr {
				if v == maxPao {
					maxPaoshou = append(maxPaoshou, int(k))
					break //只取第一个
				}
			}
		}

		for _, v := range _table.GetSeats() {

			//是否大赢家\最佳炮手
			isMaxWinner := consts.No
			isMaxPaoshou := consts.No
			for _, seatID := range maxWinner {
				if v.GetId() == seatID {
					isMaxWinner = consts.Yes
				}
			}
			for _, seatID := range maxPaoshou {
				if v.GetId() == seatID {
					isMaxPaoshou = consts.Yes
				}
			}

			//fanCt := _table.Majhong.CMajArr[v.GetId()].FanCt
			score := _table.Majhong.CMajArr[v.GetId()].Score
			logs.Info("--------------------------------------------------------->score:%v", score)
			isWin := consts.Yes
			if score < 0 {
				score *= -1
				isWin = consts.No
			}

			_result := &SeatResult{
				SeatID:   proto.Int32(int32(v.GetId())),
				IsWinner: proto.Int32(int32(isWin)),
				//FanCt:    proto.Int32(int32(fanCt)),
				Score:  proto.Int32(int32(score)),
				PxInfo: _table.Majhong.CMajArr[v.GetId()].GetPxScoreInfo(),
				HuSeq:  proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].HuSeq)),
			}
			_seatResult = append(_seatResult, _result)

			//牌局总结算
			if _table.TableCfg.GameCt-_table.Majhong.GameCt <= 0 || _table.DirectOver {
				totalScore := _table.Majhong.CMajArr[v.GetId()].TotalScore
				totalIsWin := consts.Yes
				if totalScore < 0 {
					totalScore *= -1
					totalIsWin = consts.No
				}
				_seatRt := &SeatResult{
					SeatID:   proto.Int32(int32(v.GetId())),
					IsWinner: proto.Int32(int32(totalIsWin)),
					Score:    proto.Int32(int32(totalScore)),
				}
				_gameResult = append(_gameResult, _seatRt)

				//第二个总结算界面
				_totalRt := &GameResult{
					SeatID:     proto.Int32(int32(v.GetId())),
					IsWinner:   proto.Int32(int32(totalIsWin)),
					Zimo:       proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].ZimoCt)),
					Jiepao:     proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].JiePaoCt)),
					Dianpao:    proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].DianPaoCt)),
					Angang:     proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].AnGangCt)),
					Minggang:   proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].MingGangCt)),
					Chadajiao:  proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].ChaDaJiao)),
					Chahuazhu:  proto.Int32(int32(_table.Majhong.CMajArr[v.GetId()].ChaHuaZhu)),
					Score:      proto.Int32(int32(totalScore)),
					MaxWinner:  proto.Int32(int32(isMaxWinner)),
					MaxPaoshou: proto.Int32(int32(isMaxPaoshou)),
				}
				_totalResult = append(_totalResult, _totalRt)
			}
		}
	}

	_tableCfg := &TableCfg{
		TableType:     proto.Int32(int32(_table.TableCfg.TableType)),
		PlayerCt:      proto.Int32(int32(_table.TableCfg.PlayerCt)),
		GameCt:        proto.Int32(int32(_table.TableCfg.GameCt)),
		BaseScore:     proto.Int32(int32(_table.TableCfg.BaseScore)),
		DianpaoHu:     proto.Int32(int32(_table.TableCfg.DianpaoHu)),
		ZimoHu:        proto.Int32(int32(_table.TableCfg.ZimoHu)),
		TiandiHu:      proto.Int32(int32(_table.TableCfg.TiandiHu)),
		KehuQidui:     proto.Int32(int32(_table.TableCfg.KehuQidui)),
		QiangGang:     proto.Int32(int32(_table.TableCfg.QiangGang)),
		ZhuaNiaoCt:    proto.Int32(int32(_table.TableCfg.ZhuaNiaoCt)),
		YiMaQuanZh:    proto.Int32(int32(_table.TableCfg.YiMaQuanZh)),
		MenQing:       proto.Int32(int32(_table.TableCfg.MenQing)),
		Present:       proto.Int32(int32(_table.TableCfg.Present)),
		TdqZuiCt:      proto.Int32(int32(_table.TableCfg.TdqZuiCt)),
		KePengGang:    proto.Int32(int32(_table.TableCfg.KePengGang)),
		KaiHuSuanGang: proto.Int32(int32(_table.TableCfg.KaiHuSuanGang)),
		YouGangYouFen: proto.Int32(int32(_table.TableCfg.YouGangYouFen)),
		CreaterId:     proto.Int32(int32(_table.TableCfg.CreaterId)),
		DaiHua:        proto.Int32(int32(_table.TableCfg.DaiHua)),
		FengLing:      proto.Int32(int32(_table.TableCfg.FengLing)),
		BaoTing:       proto.Int32(int32(_table.TableCfg.BaoTing)),
		WuHuaGuo:      proto.Int32(int32(_table.TableCfg.WuHuaGuo)),
		PayWay: proto.Int32(int32(_table.TableCfg.PayWay)),
	}

	//logs.Info("------------------------->RemainCt:%v", _table.Majhong.GetRemainPaiCt())
	_lastCardData := consts.DefaultIndex //牌墙最后一张
	if _tableCfg.GetTableType() == mjcomn.TableType_FYMJ {
		if _table.GetState() > consts.TableStateDealCard {
			_lastCardData = _table.Majhong.MCards[len(_table.Majhong.MCards)-1].GetData()
		}

	}
	_tableInfo := TableInfo{
		TableID:     proto.Int32(int32(_table.ID)),
		TableName:   proto.String(_table.Name),
		State:       proto.Int32(int32(_table.GetState())),
		ExecOpt:     execOpt,
		SeatInfo:    _seatInfo,
		ReaminTime:  proto.Int32(int32(_table.GetTimerRemainTime())),
		CurtSpeaker: proto.Int32(int32(_curtSpeaker)),
		RemainCt:    proto.Int32(int32(_table.Majhong.GetRemainPaiCt())),
		Flow:        proto.Int32(int32(_table.Majhong.Flow)),
		ResultList:  _seatResult,
		GameResult:  _gameResult,
		DSeatID:     proto.Int32(int32(_table.Majhong.DSeatID)),
		GameCt:      proto.Int32(int32(_table.Majhong.GameCt)),
		//FangSeatID:    proto.Int32(int32(_table.FangSeatID)),
		LastSendCard:  proto.Int32(int32(_lastSendCardData)),
		TotalResult:   _totalResult,
		ZhuaNiaoInfo:  _zhuaNiaoInfo,
		TableCfg:      _tableCfg,
		TableTime:     proto.Int32(int32(time.Now().Unix())),
		LianZhuangCt:  proto.Int32(int32(_lianZhuangCt)),
		SendTipArr:    sendTipArr,
		PresenterId:   proto.Int32(int32(_table.PresentUID)),
		PresenterName: proto.String(_table.PresentName),
		ValiNum:       proto.Int32(int32(mjcomn.GetRanExceptX(100000))),
		LastCardData:  proto.Int32(int32(_lastCardData)),
	}
	return &_tableInfo
}

//返回 protobuf.SeatInfo 对象 seat:当前需要获取的位置信息, _sendSeat:当前发送消息的对象
func (h *helper) GetSeatInfo(_table *table.Table, seat *seat.Seat, _sendSeat *seat.Seat) *SeatInfo {

	//位置操作数组 ------------------------------------------------------------
	_seatOpts := &SeatOpts{}

	if _table.GetState() == consts.TableStateWaiteThink {

		optInfo := _table.Majhong.CMajArr[seat.GetId()].OptInfo
		_seatOpts.Peng = proto.Bool(optInfo.Peng)
		_seatOpts.Gang = proto.Bool(optInfo.Gang)
		_seatOpts.Hu = proto.Bool(optInfo.Hu)
		_seatOpts.Chi = proto.Bool(optInfo.Chi)
		_seatOpts.Bu = proto.Bool(optInfo.Bu)
		_seatOpts.Cancer = proto.Bool(optInfo.Cancer)
		_seatOpts.Ting = proto.Bool(optInfo.Ting)
		_seatOpts.PengCard = make([]int32, 0)
		for _, v := range optInfo.PengCard {
			_seatOpts.PengCard = append(_seatOpts.PengCard, int32(v))
		}

		_seatOpts.GangCard = make([]int32, 0)
		for _, v := range optInfo.GangCard {
			_seatOpts.GangCard = append(_seatOpts.GangCard, int32(v))
		}

		_seatOpts.HuCard = proto.Int32(int32(optInfo.HuCard))

		_seatOpts.ChiCard = make([]int32, 0)
		for _, v := range optInfo.ChiCard {
			_seatOpts.ChiCard = append(_seatOpts.ChiCard, int32(v))
		}

		_seatOpts.BuCard = make([]int32, 0)
		for _, v := range optInfo.BuCard {
			_seatOpts.BuCard = append(_seatOpts.BuCard, int32(v))
		}
	}

	//手牌数据 打出的牌数据-----------------------------------------------------------
	_handCards := make([]int32, 0)
	_outCards := make([]int32, 0)
	_huaCards := make([]int32, 0)
	_pgInfos := make([]*PGInfo, 0)
	_scores := make([]int32, 2)    //位置积分数,第一位为符号位,第二位为值
	_tingCards := make([]int32, 0) //位置听牌数据

	if _table.GetState() >= consts.TableStateIdle &&
		_table.GetState() <= consts.TableStateShowResult {
		//位置游戏实时积分
		_totalScore := _table.Majhong.CMajArr[seat.GetId()].TotalScore
		_scores[0] = int32(consts.Yes)
		if _totalScore < 0 {
			_scores[0] = int32(consts.No)
			_totalScore = -_totalScore
		}
		_scores[1] = int32(_totalScore)
	}

	if _table.GetState() >= consts.TableStateDealCard &&
		_table.GetState() <= consts.TableStateShowResult {

		// 牌数据,只发送自己位置的,其他位置不用考虑
		if seat.GetId() == _sendSeat.GetId() {
			tingCards := _table.Majhong.CMajArr[seat.GetId()].TingCards
			for _, v := range tingCards {
				_tingCards = append(_tingCards, int32(v))
			}
		}

		//手牌处理
		_table.Majhong.CMajArr[seat.GetId()].SortHandPai() //排序
		handPai := _table.Majhong.CMajArr[seat.GetId()].GetHandPai()
		for _, v := range handPai {
			_handCards = append(_handCards, int32(v.GetData()))
		}

		//移动中发白到左边
		//if _table.GetState() > consts.TableStateDealCard {
		//	if _table.TableCfg.TableType != mjcomn.TableType_FYMJ ||
		//		_table.TableCfg.TableType != mjcomn.TableType_BBMJ ||
		//		_table.TableCfg.TableType != mjcomn.TableType_HYMJ { //阜阳、蚌埠、怀远 麻将除外
		//		tempArr := make([]int32, 0)
		//		zfbArr := make([]int32, 0)
		//		for i := 0; i < len(_handCards); i++ {
		//			_card := mjcomn.NewMCard(int(_handCards[i]))
		//			if _card.GetColor() == mjcomn.Color_Zfb {
		//				zfbArr = append(zfbArr, int32(_card.GetData()))
		//			} else {
		//				tempArr = append(tempArr, int32(_card.GetData()))
		//			}
		//
		//		}
		//		_handCards = make([]int32, 0) //重新赋值手牌
		//		for i := 0; i < len(zfbArr); i++ {
		//			_handCards = append(_handCards, zfbArr[i])
		//		}
		//		for i := 0; i < len(tempArr); i++ {
		//			_handCards = append(_handCards, tempArr[i])
		//		}
		//	}
		//}

		//判断移动一张手牌到最后一张 的情况 所有情况------------------------
		if _table.GetState() >= consts.TableStateWaiteSend {

			_moveCardData := consts.DefaultIndex //需要移动的牌值
			//如果位置已经胡牌,则移动所胡的牌
			if seat.GetState() == consts.SeatStateGameHasHu {
				_moveCardData = _table.Majhong.CMajArr[seat.GetId()].HuMCard.GetData()
			} else {
				if _table.GetState() == consts.TableStateShowResult {
					//如果位置有查大叫的牌,则移动到最后
					if _table.Majhong.CMajArr[seat.GetId()].ChaDaJiaoData != consts.DefaultIndex {
						_moveCardData = _table.Majhong.CMajArr[seat.GetId()].ChaDaJiaoData
					}
				}
			}
			// 其他情况,只考虑 自己位置信息,其他位置不用考虑
			if seat.GetId() == _sendSeat.GetId() {

				if _table.GetState() == consts.TableStateWaiteSend { //等待出牌
					//如果是当前是拿牌人并等待出牌,则将最后一张拿到的牌移到手牌最后
					//移动条件:是当前出牌人,且是牌桌最后一个拿牌的玩家 TODO 这里处理会有显示问题,可能自己出牌被别人碰,然后别人出牌自己又碰回来,这样满足条件,但不能移动
					if seat.GetId() == _table.Majhong.CurtSenderIndex && seat.GetId() == _table.Majhong.LastFetchSeatID {
						lastFetchCard := _table.Majhong.CMajArr[seat.GetId()].LastFetchMCard
						if lastFetchCard != nil { //lastFetchCard可能为空(未拿牌\碰牌后等待出牌)
							//手中可能已经没有这张牌,例如玩家拿到一张后打出,被碰,然后自己又碰,轮到出牌时则已经没有这张牌
							if _table.Majhong.CMajArr[seat.GetId()].IsHandHasSamePai(lastFetchCard.GetData()) {
								_moveCardData = lastFetchCard.GetData()
							}
						}
					}
				} else if _table.GetState() == consts.TableStateWaiteThink { //等待思考
					if _table.Majhong.HasThinker(seat.GetId()) { //玩家在当前思考列表中
						//分别讨论,是思考自己还是思考他人
						if _table.Majhong.CurtSenderIndex == seat.GetId() { //思考自己
							//胡 和明杠 (暗杠 包含刚拿到的牌)
							//canHu := _table.Majhong.CMajArr[seat.ID].CanOpt(consts.OptTypeHu)
							//if canHu {
							lastFetchCard := _table.Majhong.CMajArr[seat.GetId()].LastFetchMCard
							if lastFetchCard != nil { //lastFetchCard可能为空(未拿牌\碰牌后等待出牌)
								_moveCardData = lastFetchCard.GetData()
							}

							//}

							//canGang := _table.Majhong.CMajArr[seat.ID].CanOpt(consts.OptTypeGang)
							//logs.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~> test canGang:%v", canGang)
							//if canGang {
							//	//须检测 lastFetchCard 是否是可杠的牌
							//	lastFetchCard := _table.Majhong.CMajArr[seat.ID].LastFetchCard
							//	gangArr := _table.Majhong.CMajArr[seat.ID].OptInfo.GangCard
							//	logs.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~> lastFetchCard:%v,TempGangPaiArr:%v", lastFetchCard, gangArr)
							//	for _, v := range gangArr {
							//		//需判断 lastFetchCard 为空,庄家第一次暗杠思考时,第一发牌直接发给庄家14张,没有设置lastFetchCard值
							//		if lastFetchCard != nil && lastFetchCard.GetData() == v {
							//			_moveCardData = lastFetchCard.GetData()
							//		}
							//	}
							//}

						}
					}
				}
			}

			//移动牌到最后
			if _moveCardData != consts.DefaultIndex {
				_handCards = h.moveCardToLast(_handCards, _moveCardData)
			}
		}

		//打出的牌
		outPai := _table.Majhong.CMajArr[seat.GetId()].OutPaiArr
		for _, v := range outPai {
			_outCards = append(_outCards, int32(v.GetData()))
		}
		//打出的花牌
		huaPai := _table.Majhong.CMajArr[seat.GetId()].HuaPaiArr
		for _, v := range huaPai {
			_huaCards = append(_huaCards, int32(v.GetData()))
		}

		//碰杠吃的牌
		pgcArr := _table.Majhong.CMajArr[seat.GetId()].PGCArr
		for _, v := range pgcArr {
			_cards := make([]int32, 0)
			for i := 1; i < len(v); i++ {
				_cards = append(_cards, int32(v[i]))
			}
			_pgInfo := &PGInfo{
				Type:  proto.Int32(int32(v[0])),
				Cards: _cards,
			}

			_pgInfos = append(_pgInfos, _pgInfo)
		}

		//打印手牌值
		str := ""
		for i := 0; i < len(_handCards); i++ {
			str += mjcomn.NewMCard(int(_handCards[i])).String()
		}
		//发送给 0 号位置是打印一次
		if _sendSeat.GetId() == 0 {
			logs.Info("tableId:%v====>  playerId:%v,  seat.GetId():%v, _handCards:%v,:%v", _table.ID, seat.GetPlayer().ID(), seat.GetId(), _handCards, str)
		}

	}

	isOffline := consts.No
	if seat.GetPlayer().IsOffline() {
		isOffline = consts.Yes
	}
	_seatInfo := SeatInfo{
		PlayerInfo: h.GetPlayerInfo(seat.GetPlayer()),
		SeatId:     proto.Int32(int32(seat.GetId())),
		State:      proto.Int32(int32(seat.GetState())),
		SeatOpts:   _seatOpts,
		HandCards:  _handCards,
		OutCards:   _outCards,
		TingCards:  _tingCards,
		PgInfos:    _pgInfos,
		Score:      _scores,
		HuType:     proto.Int32(int32(_table.Majhong.CMajArr[seat.GetId()].HuType)),
		Offline:    proto.Int32(int32(isOffline)),
		Huseq:      proto.Int32(int32(_table.Majhong.CMajArr[seat.GetId()].HuSeq)),
		BuHuaCt:    proto.Int32(int32(len(_huaCards))),
		HuaCards:   _huaCards,
		TingPai:    proto.Bool(_table.Majhong.CMajArr[seat.GetId()].IsTing),
	}

	return &_seatInfo
}

//移动手牌数值中的一张到最后,客户端显示使用
func (h *helper) moveCardToLast(_handCards []int32, _cardData int) []int32 {

	//先查看手牌中是否有该张牌
	isHave := false
	tempArr := make([]int32, 0)
	for _, v := range _handCards {
		if !isHave && v == int32(_cardData) { //这里只能查找一个,因为查大叫的时候 手牌中的牌值可能相同
			isHave = true
		} else {
			tempArr = append(tempArr, v)
		}
	}

	if isHave {
		_handCards = tempArr
		_handCards = append(_handCards, int32(_cardData))
	} else {
		logs.Info("****************************find lastFetchCard error!!!!,_cardData:%v", _cardData)
	}

	return _handCards
}

//返回 protobuf.GetPlayer()Info 对象
func (h *helper) GetPlayerInfo(_player *player.Player) *PlayerInfo {
	_playerInfo := PlayerInfo{
		Uid:      proto.Int32(int32(_player.ID())),
		NickName: proto.String(_player.NickName()),
	}
	return &_playerInfo
}
