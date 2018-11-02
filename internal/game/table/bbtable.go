package table

import (
	"qianuuu.com/ahmj/internal/config"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/seat"
	. "qianuuu.com/ahmj/internal/mjcomn"
	"qianuuu.com/lib/logs"
)

// [蚌埠麻将] 牌桌
type BBTable struct {
	*Table
}

func NewBBTable(_tableID int, _robots *Robots, _tableCfg *config.TableCfg) *BBTable {
	table := NewTable(_tableID, _robots, _tableCfg)
	ret := &BBTable{
		Table: table,
	}

	ret.Table.tableInter = ret
	return ret
}

//玩家胡牌
func (t *BBTable) HuPai(_seat *seat.Seat) {

	seatId := _seat.GetId()
	huCmaj := t.Majhong.CMajArr[seatId]
	huType := huCmaj.HuType //胡牌类型
	huCard := NewMCard(huCmaj.OptInfo.HuCard)
	logs.Info("胡牌玩家手牌%v", huCmaj.GetHandPai())
	//自摸
	if huType == HUTYPE_ZIMO {

		//接炮
	} else if huType == HUTYPE_JIEPAO {

		//添加这张牌
		huCmaj.AddHandPai(huCard, false)
		lastFangPaoSeatID := t.Majhong.LastSenderSeatID

		//抢杠胡
		if t.Majhong.CMajArr[seatId].HuTypeDetail == HUTYPE_DETAIL_QIANGGANG {
			if t.Majhong.LastMianGangCard != nil {

				lastFangPaoSeatID = t.Majhong.LastMianGangSeatID
				//从面杠的玩家手中删除这张牌 ,这里可能出现多个人抢杠胡,只能删除一次
				if t.Majhong.CMajArr[t.Majhong.LastMianGangSeatID].IsHandHasSamePai(t.Majhong.LastMianGangCard.GetData()) {
					logs.Info("删除弯杠玩家手牌 t.Majhong.LastMianGangSeatID:%v,t.Majhong.LastMianGangCard:%v", t.Majhong.LastMianGangSeatID, t.Majhong.LastMianGangCard)
					t.Majhong.CMajArr[t.Majhong.LastMianGangSeatID].RemoveHandCard(t.Majhong.LastMianGangCard)
				} else {
					logs.Info("弯杠牌已经被删除,多人抢杠胡!!! t.Majhong.LastMianGangCard:%v", t.Majhong.LastMianGangCard)
				}
			}

		}

		//检测是否一炮多响
		if t.Majhong.YiPaoDXSeatID == consts.DefaultIndex {
			if t.Majhong.LastHuCardData == consts.DefaultIndex {
				t.Majhong.LastHuCardData = huCard.GetData()
			} else {
				if t.Majhong.LastHuCardData == huCard.GetData() { //胡的是同一张,说明是一炮多响
					t.Majhong.YiPaoDXSeatID = lastFangPaoSeatID
				}
			}
		}
	}

	logs.Info("tableId:%v--------------->位置:%v胡牌,所胡牌:%v,类型:%v", t.ID, seatId, huCard, GetHuTypeName(huType))

	t.Majhong.CMajArr[seatId].ResetOptInfo() //清空本位置思考操作
	t.Majhong.RemoveThinker(seatId)          //删除思考位置

	isMultHu := false //是否多人胡牌
	if t.Majhong.GetThinkerCt() > 0 {
		isMultHu = true

		t.SendTableInfo()
		//多人胡牌,第1~2个胡操作不发送
		//logs.Info("tableId:%v-------------->多人胡牌,第1~2个胡操作不发送----------t.Majhong.GetThinkerCt():%v", t.ID, t.Majhong.GetThinkerCt())
		//t.ExecOptInfo = NewExecOptInfo()
	}

	_seat.SetState(consts.SeatStateGameHasHu)

	//结算位置胡牌信息
	t.calHuResult(_seat)

	if !isMultHu {
		if t.tableInter.ChkOver() {
			logs.Info("tableId:%v-------------->玩家胡了最后一张别人出的牌,游戏结束!", t.ID)
			return
		}
	}

	//仍然有胡牌思考者,刷新牌桌,等待其他胡牌玩家思考
	if isMultHu {
		t.SendTableInfo()
		return //等待另外一个胡牌玩家思考
	}

	//最后一人胡牌->检测抓鸟
	//t.ChkZhuaNiao()
	//不抓鸟,进入游戏结算
	t.GameOver()
}

//位置胡牌,计算牌型\积分等
func (t *BBTable) calHuResult(_seat *seat.Seat) {

	huSeatID := _seat.GetId()
	huCmaj := t.Majhong.CMajArr[huSeatID]
	huType := t.Majhong.CMajArr[huSeatID].HuType
	huTypeDetail := huCmaj.HuTypeDetail

	//计算 应扣分的座位id
	loseSeatIDArr := make([]int, 0)

	if huType == HUTYPE_JIEPAO { //接炮
		//抢杠胡
		if huTypeDetail == HUTYPE_DETAIL_QIANGGANG {
			loseSeatIDArr = append(loseSeatIDArr, t.Majhong.LastMianGangSeatID)
			t.Majhong.CMajArr[t.Majhong.LastMianGangSeatID].DianPaoCt++ //点炮次数

			//点杠花 -接炮
		} else if huTypeDetail == HUTYPE_DETAIL_GANG_SHANG_PAO {
			loseSeatIDArr = append(loseSeatIDArr, t.Majhong.LastDianZhiGangSeatID)
			t.Majhong.CMajArr[t.Majhong.LastDianZhiGangSeatID].DianPaoCt++ //点炮次数

			//普通胡牌
		} else {
			loseSeatIDArr = append(loseSeatIDArr, t.Majhong.LastSenderSeatID)
			t.Majhong.CMajArr[t.Majhong.LastSenderSeatID].DianPaoCt++ //点炮次数
		}

		huCmaj.JiePaoCt++ //接炮次数

	} else if huType == HUTYPE_ZIMO { //自摸
		for i := 0; i < t.TableCfg.PlayerCt; i++ {
			if i != huSeatID {
				loseSeat := t.seats[i]
				if loseSeat.GetState() != consts.SeatStateGameHasHu { //不包括已胡牌的玩家
					loseSeatIDArr = append(loseSeatIDArr, i)
				}
			}
		}
		//记录统计数据
		huCmaj.ZimoCt++ //自摸次数
	}

	//积分计算 胡牌基础分数  ----------------------------------------------------------------
	totalScore := 0  //总分
	zhuangScore := 1 // 庄家输赢+-1
	baseScore := 1   //点炮基础分
	if huType == HUTYPE_ZIMO {
		baseScore = 2 //自摸1分
	}

	//胡牌基础分-----------------------------------
	if huType == HUTYPE_JIEPAO {
		if huTypeDetail == HUTYPE_DETAIL_QIANGGANG {
			huCmaj.AddPxScore("抢杠胡", baseScore)
		}
		huCmaj.AddPxScore("接炮", baseScore)
	} else {
		huCmaj.AddPxScore("自摸", baseScore)
	}
	totalScore += baseScore
	//庄分数-----------------------------------
	if t.Majhong.DSeatID == huSeatID {
		totalScore += zhuangScore
		huCmaj.AddPxScore("庄", zhuangScore)

	}
	//杠分分数-----------------------------------
	if huCmaj.GetGangCt() > 0 {
		allCt := huCmaj.GetGangCt() / 4 //杠牌总数（暗杠+明杠）
		anGang := huCmaj.AgCt

		if anGang > 0 {
			anGangScore := anGang * 2
			totalScore += anGangScore
			huCmaj.AddPxScore("暗杠", anGangScore)
		}

		gang := allCt - anGang
		if gang > 0 {
			totalScore += gang
			huCmaj.AddPxScore("明杠", gang)
		}

	}
	//补花分数-----------------------------------
	if len(huCmaj.HuaPaiArr) > 0 {
		totalScore += len(huCmaj.HuaPaiArr)
		huCmaj.AddPxScore("补花", len(huCmaj.HuaPaiArr))
	}
	//抢杠胡分数-----------------------------------
	if t.TableCfg.DaiHua <= 0 { //不带花 *3  带花抢杠胡相当于普通
		if huTypeDetail == HUTYPE_DETAIL_QIANGGANG {
			totalScore = totalScore * 3
			huCmaj.AddPxScore("抢杠 *3", 0)
		}
	}
	if huTypeDetail == HUTYPE_DETAIL_GANG_SHANG_HUA && t.TableCfg.DaiHua <= 0 { //不带花的情况 杠上花所有分数*2
		totalScore = totalScore * 2
		huCmaj.AddPxScore("杠上花 *2", 0)

	}

	logs.Info("tableId:%v--------calHuResult--------->loseSeatIDArr:%v,huCmaj.HuPaiRec[huSeatID]:%v",
		t.TableCfg.TableId, loseSeatIDArr, huCmaj.HuPaiRec[huSeatID])
	logs.Info("huCmaj.GetPxScoreInfo()===========>%v", huCmaj.GetPxScoreInfo())
	//记录输赢分
	for _, _loseSeatID := range loseSeatIDArr {
		zhuang := 0
		if t.Majhong.DSeatID == _loseSeatID {
			loseCmaj := t.Majhong.CMajArr[_loseSeatID]
			loseCmaj.AddPxScore("庄", -1)
			zhuang = zhuangScore
			//if huTypeDetail == HUTYPE_DETAIL_GANG_SHANG_HUA {
			//	totalScore += 2
			//} else {
			//	totalScore += 1
			//}
		}
		loseScore := totalScore + zhuang
		huCmaj.HuPaiRec[_loseSeatID] += loseScore //输家积分
		huCmaj.HuPaiRec[huSeatID] += loseScore    //赢家积分
	}

}

// 确定庄家
func (t *BBTable) MakeDSeat() {

	//测试
	if !config.Opts().CloseTest {

		t.Majhong.DSeatID = 0
		t.Majhong.CurtSenderIndex = t.Majhong.DSeatID
		logs.Info("tableId:%v---------------->makeDSeat    m.DUser:%v", t.ID, t.Majhong.DSeatID)
		return
	}

	//第一局开房玩家坐庄
	//若黄庄 则上一局的庄家为庄
	//上一局 如有玩家胡牌,则胡牌玩家坐庄。
	if t.Majhong.GameCt == 1 {
		t.Majhong.DSeatID = t.FangSeatID
		logs.Info("tableId:%v---------------->makeDSeat   房主:%v", t.ID, t.FangSeatID)
	} else {
		//之后将由上局庄下一座位坐庄
		zhuang := t.Majhong.DSeatID
		logs.Info("tableId:%v---------------->makeDSeat  上局庄 位置:%v", t.TableCfg.TableId, zhuang)
		if t.Majhong.FirstHuSeatID != consts.DefaultIndex {
			if t.Majhong.YiPaoDXSeatID != consts.DefaultIndex {
				//if zhuang == t.Majhong.YiPaoDXSeatID {
				//
				//} else { //庄没有胡下一家坐庄
				if (zhuang + 1) > (t.TableCfg.PlayerCt - 1) {
					zhuang = 0
				} else {
					zhuang += 1
				}
				//}
				t.Majhong.DSeatID = zhuang
				logs.Info("tableId:%v---------------->makeDSeat  上局一炮多响 位置:%v", t.TableCfg.TableId, t.Majhong.YiPaoDXSeatID)
			} else {
				if zhuang == t.Majhong.FirstHuSeatID {

				} else { //庄没有胡下一家坐庄
					if (zhuang + 1) > (t.TableCfg.PlayerCt - 1) {
						zhuang = 0
					} else {
						zhuang += 1
					}
				}
				t.Majhong.DSeatID = zhuang
				logs.Info("tableId:%v---------------->makeDSeat  上局胡牌 位置:%v", t.TableCfg.TableId, t.Majhong.FirstHuSeatID)
			}
		} else {
			//无任何人胡牌则上一局的庄家为庄
			if t.Majhong.LastZhuangSeatID == consts.DefaultIndex {
				t.Majhong.LastZhuangSeatID = zhuang
			}
			t.Majhong.DSeatID = t.Majhong.LastZhuangSeatID
			logs.Info("tableId:%v---------------->makeDSeat   流局 上个庄家位置:%v", t.TableCfg.TableId, t.Majhong.LastZhuangSeatID)
		}

	}

	t.Majhong.CurtSenderIndex = t.Majhong.DSeatID
	t.Majhong.LastZhuangSeatID = t.Majhong.DSeatID

	logs.Info("tableId:%v---------------->makeDSeat   庄家:%v", t.TableCfg.TableId, t.Majhong.DSeatID)
}

//拿牌者拿牌后开始思考  _specialPxId:某些状态下思考 带入牌型
func (t *BBTable) FetcherStartThink(_specialPxId int) {
	//检测自己思考
	isSelfThink := t.ThinkSelfPai(_specialPxId, false)

	//需要思考
	if isSelfThink {
		t.setState(consts.TableStateWaiteThink)
	} else {
		_seatId := t.Majhong.CurtSenderIndex
		_cmaj := t.Majhong.CMajArr[_seatId]

		if _cmaj.Check_BB_BuHua() && t.TableCfg.DaiHua > 0 {
			//系统自动补花
			t.buhuaLogic()
		} else {
			//无须思考,等待出牌
			t.tableInter.SendCheckTing()
		}
	}
}

//拿到牌后思考自己的牌 isPeng:是否碰后思考
func (t *BBTable) ThinkSelfPai(_specialPxId int, isPeng bool) bool {

	m := t.Majhong
	m.ClearThinker()
	seatId := m.CurtSenderIndex
	m.AddThinker(seatId)
	curtCmaj := m.CMajArr[seatId]
	curtCmaj.ResetOptInfo() //清空操作
	logs.Info("tableId:%v---------------------->位置:%v开始思考", m.TableCfg.TableId, m.CurtSenderIndex)

	// 检测补花 --------------------
	//if t.TableCfg.DaiHua > 0 {
	//	isbuhua := false
	//	if m.GetRemainPaiCt() >= 1 {
	//			isbuhua = curtCmaj.Check_BB_BuHua()
	//	}
	//	if isbuhua {
	//		//curtCmaj.SetOpt(OptTypeBu)
	//		logs.Info("tableId:%v,[位置%v]可补花", m.TableCfg.TableId, m.CurtSenderIndex)
	//		logs.Info("tableId:%v,--------->curtCmaj.TempBuHuaArr:%v", m.TableCfg.TableId, curtCmaj.TempBuHuaArr)
	//		return false
	//	}
	//
	//}

	// 检测hu牌 --------------------
	isHu := false
	if !isPeng {
		isHu = t.CheckHuPai(seatId)
	}
	if isHu {
		//带入牌型:杠后开花
		if _specialPxId != consts.DefaultIndex {
			//if _specialPxId == consts.PXID_FYMJ_GANG_SHANG_HUA { //杠上花牌型 判断 扎杠
			//		curtCmaj.AddPxId(_specialPxId)
			//} else {
			curtCmaj.HuTypeDetail = _specialPxId
			//curtCmaj.AddPxId(_specialPxId)
			//}
		}

		curtCmaj.SetOpt(OptTypeHu)    //设置胡操作
		curtCmaj.HuType = HUTYPE_ZIMO //胡牌类型为自摸

		//记录位置所胡的牌(天胡也默认为最后拿到的那张)
		huCard := curtCmaj.LastFetchMCard.Clone()
		curtCmaj.OptInfo.HuCard = huCard.GetData()
		curtCmaj.HuMCard = NewMCard(huCard.GetData())

		logs.Info("tableId:%v---------------------->[位置%v]可胡牌,所胡牌:%v,牌型:%v",
			m.TableCfg.TableId, m.CurtSenderIndex, huCard.String(), curtCmaj.HuPxIdArr)
	}

	// 检测暗杠--------------------
	isAnGang := false
	if m.GetRemainPaiCt() >= 1 { //牌墙中必须大于1张时才能思考暗杠
		isAnGang = t.Check_AnGang(seatId)
	}

	if isAnGang {
		logs.Info("tableId:%v---------------------->[位置%v]可暗杠,杠牌:%v", m.TableCfg.TableId, m.CurtSenderIndex, curtCmaj.TempGangPaiArr)
		curtCmaj.SetOpt(OptTypeGang)
		//将杠牌加入
		for i := 0; i < len(curtCmaj.TempGangPaiArr); i = i + 4 {
			//暗杠可能是刚拿到的牌,发给客户端时要移动到最后,如果刚拿到的牌是杠牌,则设置成刚拿到的
			if m.CMajArr[m.CurtSenderIndex].LastFetchMCard != nil &&
				m.CMajArr[m.CurtSenderIndex].LastFetchMCard.Equal(curtCmaj.TempGangPaiArr[i]) {
				curtCmaj.OptInfo.AddGangCard(m.CMajArr[m.CurtSenderIndex].LastFetchMCard.GetData())
			} else {
				curtCmaj.OptInfo.AddGangCard(curtCmaj.TempGangPaiArr[i].GetData())
			}
		}
	}

	// 检测面杠 --------------------
	isMianGang := false
	remainCt := m.GetRemainPaiCt()
	if remainCt >= 1 {
		isMianGang = curtCmaj.Check_MianGang()
	}
	if isMianGang {
		curtCmaj.SetOpt(OptTypeGang)
		logs.Info("tableId:%v,[位置%v]可面杠", m.TableCfg.TableId, m.CurtSenderIndex)
		logs.Info("tableId:%v,--------->curtCmaj.TempGangPaiArr:%v", m.TableCfg.TableId, curtCmaj.TempGangPaiArr)
		//将面杠牌加入
		for i := 0; i < len(curtCmaj.TempGangPaiArr); i = i + 4 {

			gangCard := curtCmaj.TempGangPaiArr[i]
			//找出手牌上相等的牌并设置
			realGangCard := curtCmaj.GetHandEqualPai(gangCard)
			logs.Info("tableId:%v,--------->设置面杠牌realGangCard:%v", m.TableCfg.TableId, realGangCard)
			if realGangCard != nil {
				curtCmaj.OptInfo.GangCard = append(curtCmaj.OptInfo.GangCard, realGangCard.GetData())
			}
		}
	}

	return isHu || isAnGang || isMianGang
}

//出牌后其他位置玩家思考
func (t *BBTable) ThinkOtherPai(_tkCard *MCard) bool {

	m := t.Majhong

	seatThink := make([]bool, 4)

	logs.Info("tableId:%v,---------------->m.CurtThinkerArr:%v", m.TableCfg.TableId, m.CurtThinkerArr)
	//按位置顺序检测思考
	for i := 0; i < len(m.CurtThinkerArr); i++ {

		_seatID := m.CurtThinkerArr[i]
		_tkCmaj := m.CMajArr[_seatID]
		_tkCmaj.ResetOptInfo() //清空操作

		// 检测胡牌 ------------------------------------------------------------
		isHu := false
		_tkCmaj.AddHandPai(_tkCard, false) //先添加这张牌

		//自摸胡不带点炮
		if t.TableCfg.ZimoHu == consts.Yes {
			logs.Info("tableId:%v,---------------->自摸胡,不带点炮", m.TableCfg.TableId)
		} else {

			////点数限制,杠开没有限制
			//huLimit := false
			//_color := _tkCard.GetColor()
			//_value := _tkCard.GetValue()
			//if _color == Color_Feng || _color == Color_Zfb {
			//	huLimit = true
			//} else { //万筒条
			//	if _value < 4 {
			//		huLimit = true
			//	}
			//}
			//
			//if huLimit {
			//	logs.Info("tableId:%v,---------------->点炮胡牌 点数限制,_tkCard:%v", m.TableCfg.TableId, _tkCard)
			//} else {
			isHu = t.CheckHuPai(_seatID)
			//}
		}

		if isHu {
			_tkCmaj.HuType = HUTYPE_JIEPAO
			//思考他牌时检测过手胡
			curtFanCt := m.GetSeatFanCt(_seatID, false)
			isGuoShouHu := m.CMajArr[_seatID].CheckGSH(curtFanCt)
			if isGuoShouHu && _tkCmaj.GuoShouHuCard.Equal(_tkCard) {
				logs.Info("tableId:%v------------------->[位置%v]过手胡", _seatID, m.TableCfg.TableId)
				isHu = false //将胡置为false
			}
			if isHu {
				//检测地胡 只有庄家拿了一张牌,其他人没拿牌
				if m.TableCfg.TiandiHu == consts.Yes {
					//添加座位号非庄家判断 (庄家打出一张牌被其他人碰,然后碰的出一张又被庄家胡牌,导致判断地胡)
					//最近一个出牌人为庄家,非闲家出牌,这个判断已经包含了上面的情况
					if m.PaiIndex == (13*m.TableCfg.PlayerCt+1) &&
						_seatID != m.DSeatID &&
						m.LastSenderSeatID == m.DSeatID {
						logs.Info("tableId:%v------------------->牌型为 [地胡]", m.TableCfg.TableId)
						_tkCmaj.AddPxId(consts.PXID_HZMJ_DIHU)
					}
				}

				//保存所胡的牌
				_tkCmaj.SetOpt(OptTypeHu)
				huCard := m.GetSeatHuCard(_seatID)
				_tkCmaj.OptInfo.HuCard = huCard.GetData()
				_tkCmaj.HuMCard = NewMCard(huCard.GetData())

				logs.Info("tableId:%v------------------->[位置%v]可胡牌,牌=%v", m.TableCfg.TableId, _seatID, huCard.String())
			}
		}

		//检测完胡后删除这张牌,无论是否胡牌
		_tkCmaj.RemoveHandCard(_tkCard)

		// 检测碰牌 ------------------------------------------------------------
		isPeng := _tkCmaj.Check_Peng(_tkCard)
		if isPeng {
			isGuoShouPeng := m.CMajArr[_seatID].CheckGSP()
			if isGuoShouPeng && _tkCmaj.GuoShouPengCard.Equal(_tkCard) {
				logs.Info("tableId:%v------------------->[位置%v]过手碰", _seatID, m.TableCfg.TableId)
				isPeng = false

			}
			if isPeng {
				_tkCmaj.SetOpt(OptTypePeng)
				_tkCmaj.OptInfo.PengCard = append(_tkCmaj.OptInfo.PengCard, _tkCard.GetData())
				logs.Info("tableId:%v------------------->[位置%v]可碰牌", m.TableCfg.TableId, _seatID)
			}

		}

		// 检测明杠 ------------------------------------------------------------
		isMingGang := false
		if t.Majhong.GetRemainPaiCt() >= 1 { //有牌才能杠
			isMingGang = _tkCmaj.Check_MingGang(_tkCard)
		}

		if isMingGang {
			_tkCmaj.SetOpt(OptTypeGang)
			_tkCmaj.OptInfo.GangCard = append(_tkCmaj.OptInfo.GangCard, _tkCard.GetData())
			logs.Info("tableId:%v,[位置%v]可直杠", m.TableCfg.TableId, _seatID)
		}

		seatThink[_seatID] = isHu || isPeng || isMingGang
		if !seatThink[_seatID] { //如果该位置无思考,则删除
			logs.Info("tableId:%v------------------->位置:%v无思考", m.TableCfg.TableId, _seatID)
			m.RemoveThinker(_seatID)
			i-- // 如有元素删除,下标 -1
		}

	}

	logs.Info("tableId:%v------------------->seatThink:seatThink", seatThink)
	//只要有一个位置思考,则返回true
	for _, v := range seatThink {
		if v {
			return true
		}
	}

	return false
}

// 胡牌检测
func (t *BBTable) CheckHuPai(_seatId int) bool {

	if !config.Opts().OpenHu {
		return false
	}

	_cmaj := t.Majhong.CMajArr[_seatId]

	//清空牌型记录
	_cmaj.ClearPxId()
	_cmaj.ClearExtPxId()

	//检测胡牌
	_chuobj := t.GetCalHuObject(_seatId)

	_calHuInfo := CalHuPai(_chuobj.MjCfg, _chuobj.HandPaiArr)
	_cmaj.SaveCalHuInfo = _calHuInfo

	return t.Check_hupaiRule(_seatId, _calHuInfo)

}
func (t *BBTable) Check_hupaiRule(_seatId int, _calHuInfo *CalHuInfo) bool {
	_cmaj := t.Majhong.CMajArr[_seatId]
	logs.Info("tableId:%v--------CheckHuPai---------->_calHuInfo.PxType:%v", t.TableCfg.TableId, _calHuInfo.PxType)

	if t.TableCfg.DaiHua > 0 {
		handPai := _cmaj.GetHandPai()
		for _, v := range handPai {
			if v.IsAQHuaPai() {
				return false
			}
		}
	}
	if _calHuInfo.PxType > PXTYPE_UNKNOW {
		//t.Check_PxId(_seatId, _calHuInfo)
		if _calHuInfo.PxType == PXTYPE_PINGHU {
			_cmaj.AddPxId(consts.PXID_HYMJ_PINGHU)
			logs.Info("tableId:%v--------CheckHuPai---------->平胡胡牌", t.TableCfg.TableId)

		}
		//if len(_cmaj.HuPxIdArr) > 0  {
		return true
		//}
	}

	return false
}

//获取胡牌运算对象
func (t *BBTable) GetCalHuObject(_seatId int) *CalHuObject {

	_mjcfg := NewMjCfg()
	_mjcfg.TableType = t.TableCfg.TableType
	_mjcfg.MaxColorCt = 5

	_cmaj := t.Majhong.CMajArr[_seatId]
	if t.TableCfg.KehuQidui == consts.Yes {
		if _cmaj.GetPengCt() == 0 && _cmaj.GetGangCt() == 0 { //未碰\未杠
			_mjcfg.Check7Dui = true
		}
	}

	_calHuObject := NewCalHuObject(_mjcfg)
	//手牌
	for _, v := range _cmaj.HandPaiArr {
		for _, n := range v {
			_color := n.GetColor()
			_calHuObject.HandPaiArr[_color] = append(_calHuObject.HandPaiArr[_color], n)
		}
	}

	return _calHuObject
}

// 胡牌牌型检测
func (t *BBTable) Check_PxId(_seatId int, _calHuInfo *CalHuInfo) bool {

	huCmaj := t.Majhong.CMajArr[_seatId]

	//清一色
	if huCmaj.Check_QYS() {
		huCmaj.AddPxId(consts.PXID_HZMJ_QYS)
		logs.Info("tableId:%v--------Check_PxId---------->清一色", t.TableCfg.TableId)
	}

	//对对胡
	if huCmaj.Check_DDHu(_calHuInfo) {
		huCmaj.AddPxId(consts.PXID_HZMJ_DDH)
		logs.Info("tableId:%v--------Check_PxId---------->对对胡", t.TableCfg.TableId)
	}

	return false
}

//// 其他牌型检测
//func (t *BBTable) Check_ExtPxId(_seatId int) bool {
//
//	huCmaj := t.Majhong.CMajArr[_seatId]
//
//	if t.TableCfg.MenQing == consts.Yes {
//		if huCmaj.Check_MQ() {
//			huCmaj.AddExtPxId(consts.EXTPXID_HZMJ_MENQING)
//		}
//	}
//
//	if t.Check_SUPAI(_seatId) {
//		huCmaj.AddExtPxId(consts.EXTPXID_HZMJ_SUPAI)
//	}
//
//	return false
//}

////检测 素牌:没有红中
//func (t *BBTable) Check_SUPAI(_seatId int) bool {
//
//	_cmaj := t.Majhong.CMajArr[_seatId]
//	for _, v := range _cmaj.GetHandPai() {
//		if t.Majhong.IsLaizi(v) {
//			return false
//		}
//	}
//
//	return true
//}

// 暗杠检测
func (t *BBTable) Check_AnGang(_seatId int) bool {

	if !config.Opts().OpenGang {
		return false
	}

	_cmaj := t.Majhong.CMajArr[_seatId]
	_cmaj.SortHandPai()
	_cmaj.TempGangPaiArr = make([]*MCard, 0)

	ct := 0
	for _, v := range _cmaj.HandPaiArr {
		iSize := len(v)
		if iSize >= 4 {
			count := 0
			value := -1
			for m, n := range v {
				////红中不参与暗杠检测
				//if n.IsHongZhong() {
				//	continue
				//}
				val := n.GetValue()
				if val != value {
					count = 1
					value = n.GetValue()
				} else {
					count++
				}

				if count == 4 {

					//非选缺或不是选缺牌则加入
					for k := 0; k < 4; k++ {
						_cmaj.TempGangPaiArr = append(_cmaj.TempGangPaiArr, v[m-k])
						ct++
					}
				}
			}
		}
	}
	if ct >= 4 {
		return true
	}
	return false
}

// 最后一个可胡玩家选择取消胡
func (t *BBTable) LastHuThinkerCancer() {
	logs.Info("tableId:%v-------------------> 最后一个可胡玩家选择取消胡牌")
t.GameOver()
	//检测抓鸟
	//t.ChkZhuaNiao()
}

//庄家第一次思考
func (t *BBTable) DPlayerFirstThink() {

	//庄家开始思考 //检测天胡
	_pxId := consts.DefaultIndex
	//_pxId :=  HUTYPE_DETAIL_GANG_SHANG_HUA
	if t.TableCfg.TiandiHu == consts.Yes {
		if t.Majhong.PaiIndex == (13*t.TableCfg.PlayerCt + 1) { //庄家起手第一张胡牌
			_pxId = consts.PXID_HYMJ_TIANHU
		}
	}
	t.FetcherStartThink(_pxId)
}

// 等待出牌->听牌检测
func (t *BBTable) SendCheckTing() {

	_seatId := t.Majhong.CurtSenderIndex
	_cmaj := t.Majhong.CMajArr[_seatId]

	//检测听牌
	_calHuObject := t.GetCalHuObject(_seatId)
	_sendTipArr := _calHuObject.CalTingPai()

	logs.Info("tableId:%v------------------SendCheckTing()-------------->len(_sendTipArr):%v", t.TableCfg.TableId, len(_sendTipArr))

	_cmaj.SendTipArr = make([]*SendTip, 0)
	for _, v := range _sendTipArr {
		sendCard := NewMCard(v.SendCard)
		huCard := make([]*MCard, 0)
		hudata := make([]int, 0)
		sendtip := NewSendTip()
		for i := 0; i < len(v.HuCards); i++ {
			//听牌检测
			_calHuInfo := v.HuInfos[i]
			_cmaj.ClearPxId()
			_cmaj.SaveCalHuInfo = &_calHuInfo
			logs.Info("tableId:%v------------------打出的牌------------->打出%v可胡:%v ,牌型:%v", t.TableCfg.TableId, sendCard, v.HuCards, _calHuInfo.PxList)
			ishu := t.Check_hupaiRule(_seatId, &_calHuInfo) //检测是否符合当前牌桌的胡牌
			if ishu {
				hudata = append(hudata, v.HuCards[i])
				huCard = append(huCard, NewMCard(v.HuCards[i]))
				sendtip.SendCard = v.SendCard
				sendtip.HuInfos = append(sendtip.HuInfos, _calHuInfo)
				sendtip.HuCards = hudata

			}
		}
		if len(huCard) > 0 {
			logs.Info("tableId:%v------------------SendCheckTing()-------------->打出%v可胡:%v", t.TableCfg.TableId, sendCard, huCard)
			_cmaj.SendTipArr = append(_cmaj.SendTipArr, sendtip)
			//_cmaj.SendTipArr = append(_cmaj.SendTipArr, v)
		}

	}

	t.setState(consts.TableStateWaiteSend)
}

// 更新位置当前听牌
func (t *BBTable) UpdateTingCards(_seatId int) {

	_cmaj := t.Majhong.CMajArr[_seatId]
	//先清空听牌,有可能打完牌后不听牌
	_cmaj.TingCards = make([]int, 0)

	//logs.Info("----------------------UpdateTingCards()------------------len(_cmaj.SendTipArr):%v", len(_cmaj.SendTipArr))
	if len(_cmaj.SendTipArr) > 0 { //之前检测有听牌
		lastSendCard := _cmaj.LastSendMCard
		//检测上一次打出的牌是否在SendTipArr中
		//logs.Info("----------------------UpdateTingCards()------------------lastSendCard:%v", lastSendCard)
		if lastSendCard != nil {
			for _, v := range _cmaj.SendTipArr {
				if NewMCard(v.SendCard).Equal(lastSendCard) { //判断相等,不是相同
					for _, n := range v.HuCards {
						_cmaj.TingCards = append(_cmaj.TingCards, n)
					}
					//logs.Info("----------------------UpdateTingCards()------------------_cmaj.TingCards:%v", _cmaj.TingCards)
					break
				}
			}
		}
	}
}

//玩家杠牌
func (t *BBTable) Gang(_seatId int, _data int, _isSave bool) {

	_cmaj := t.Majhong.CMajArr[_seatId]
	_gangCard := NewMCard(_data)
	_cmaj.CancerGSH() //杠牌取消过手胡
	_cmaj.CancerGSP() //杠牌取消过手碰

	//logs.Info("tableId:%v,--------------->_gangType:%v,_isSave:%v ", t.ID, _gangType, _isSave)
	_gangType := t.getGangType(_seatId, _data)
	t.ExecOptInfo.OptDetail = _gangType

	//面杠时判断抢杠胡
	if _gangType == GANGTYPE_MIAN {
		if !_isSave {
			if t.TableCfg.QiangGang == consts.Yes { //可抢杠胡
				//isReGang := false
				//if _cmaj.InReGangState {
				//	if _gangCard.EqualByData(_cmaj.ReGangCardData) {
				//		logs.Info("tableId:%v,--------------->抢杠操作,位置 %v 热杠,不允许抢杠操作! ",
				//			t.ID, _seatId)
				//	}
				//}
				//if isReGang {
				//
				//} else {
				t.addOtherThinker(_seatId)
				isQiangGang := t.ThinkQiangGangHu(_seatId, _gangCard)
				if isQiangGang {
					//暂存面杠操作
					t.SaveOpt = []int{_seatId, OptTypeGang, _data}
					logs.Info("tableId:%v,--------------->抢杠操作,位置 %v 面杠操作被保存,操作数组:%v ",
						t.ID, _seatId, _cmaj.OptInfo)

					//删除位置操作
					_cmaj.ResetOptInfo() //清空本位置操作 ,不删除思考者,用于第二个玩家操作时判断
					t.ExecOptInfo = NewExecOptInfo()
					t.SendTableInfo() //刷新操作按钮
					return            //等待抢杠玩家思考
					//}
				}
			}
		}
	}

	//执行杠操作
	logs.Info("tableId:%v-------------------->位置:%v杠牌,杠类型:%v", t.ID, _seatId, GetGangTypeName(_gangType))
	if _gangType == GANGTYPE_AN {

		_cmaj.DoAnGang(_gangCard)
		t.GuaFXiay_AM(_seatId, 2)

		_cmaj.AgCt++
		_cmaj.AnGangCt++ //暗杠次数

	} else if _gangType == GANGTYPE_ZHI { //直杠,只扣点杠玩家的积分

		_cmaj.DoMingGang(_gangCard)
		//立即计算刮风下雨所得积分
		//_dianGangSeatID := t.Majhong.LastSenderSeatID

		_score := 1 //点杠
		//t.GuaFXiay_Z(_seatId, _score, _dianGangSeatID)
		t.GuaFXiay_AM(_seatId, _score) //和面杠一样处理
		_cmaj.ZhiGCt++
		_cmaj.MingGangCt++ //明杠次数

		//记录点杠玩家
		dianGangSeatID := t.Majhong.LastSenderSeatID
		t.Majhong.CMajArr[dianGangSeatID].DgCt++
		t.Majhong.LastDianZhiGangSeatID = dianGangSeatID

		//2.点杠包赔：在每局剩7摞牌提示弹出后，若牌局中有玩家出牌，且该牌被其他玩家杠走，并且直接杠上开花，
		//则本局需要包赔，若开杠玩家没有直接胡牌，则本局不需要包赔
		if t.Majhong.GetRemainPaiCt() <= 14 {
			_cmaj.BaoGangSeatId = dianGangSeatID //如果多个,只记录最后一个
			logs.Info("tableId:%v-------------------->位置:%v 杠牌,最后14张,点杠包赔 dianGangSeatID:%v",
				t.ID, _seatId, dianGangSeatID)
		}

	} else if _gangType == GANGTYPE_MIAN { //面杠

		_cmaj.DoMianGang(_gangCard)
		_cmaj.MianGCt++
		_cmaj.MingGangCt++ //明杠次数
		t.GuaFXiay_AM(_seatId, 1)
	}

	_cmaj.ResetOptInfo() //清空本位置思考操作

	//发送杠牌动作
	t.SendTableInfo()

	//杠完后继续拿牌
	t.Majhong.SetCurtSenderIndex(_seatId)
	t.playerFetchPai(false)

	//杠后拿的一张牌设置牌型 ［杠上花］
	_cmaj.LastGangType = _gangType //记录杠类型,判断 直杠杠后花是放炮 还是自摸

	//杠后等待出牌 标记玩家进入 [杠上炮] 状态
	t.Majhong.LastDoGangSeatID = _seatId           //记录最近一个杠的位置
	_cmaj.GangShangPao = true                      //记录进入杠上炮状态
	t.Majhong.LastDoGangType = _gangType           //记录杠牌类型
	t.Majhong.LastDoGangData = _gangCard.GetData() //记录杠牌值

	_pxId := HUTYPE_DETAIL_GANG_SHANG_HUA //带入杠上花牌型
	t.FetcherStartThink(_pxId)

}

//玩家面杠操作后,查看其他位置是否抢杠胡 _gangSeatID:面杠操作的位置id
func (t *BBTable) ThinkQiangGangHu(_gangSeatID int, _tkCard *MCard) bool {

	m := t.Majhong
	seatThink := make([]bool, 4)

	//按位置顺序检测思考
	for i := 0; i < len(m.CurtThinkerArr); i++ {
		_seatID := m.CurtThinkerArr[i]
		if _seatID == _gangSeatID { //IMT!! 此时面杠玩家也在 思考列表中,不能再次参与检测
			continue
		}

		_tkCmaj := m.CMajArr[_seatID]
		_tkCmaj.ResetOptInfo() //清空操作

		// 检测hu牌 -------------------------------------
		isHu := false

		_tkCmaj.AddHandPai(_tkCard, false) //先添加这张牌
		isHu = t.tableInter.CheckHuPai(_seatID)
		_tkCmaj.RemoveHandCard(_tkCard) //检测完胡后删除这张牌

		//点数限制,杠开没有限制
		//huLimit := false
		//_color := _tkCard.GetColor()
		//_value := _tkCard.GetValue()
		//if _color == Color_Feng || _color == Color_Zfb {
		//	huLimit = true
		//} else { //万筒条
		//	if _value < 5 {
		//		huLimit = true
		//	}
		//}
		//
		//if huLimit {
		//	isHu = false
		//	logs.Info("tableId:%v,---------------->抢杠胡牌 点数限制,_tkCard:%v", m.TableCfg.TableId, _tkCard)
		//}

		if isHu {
			_tkCmaj.SetOpt(OptTypeHu)
			_tkCmaj.HuType = HUTYPE_JIEPAO
			_tkCmaj.HuTypeDetail = HUTYPE_DETAIL_QIANGGANG
			_tkCmaj.OptInfo.HuCard = _tkCard.GetData()
			_tkCmaj.HuMCard = NewMCard(_tkCard.GetData())
			logs.Info("tableId:%v ------------------>[位置%v]可 抢杠胡,_tkCard:%v,_gangSeatID:%v ",
				m.TableCfg.TableId, _seatID, _tkCard, _gangSeatID)

			//记录这张 可能被抢杠胡的 牌和位置
			m.LastMianGangCard = _tkCard.Clone()
			m.LastMianGangSeatID = _gangSeatID
		}

		seatThink[_seatID] = isHu
		if !seatThink[_seatID] { //如果该位置无思考,则删除
			logs.Info("tableId:%v------------------>位置:%v无思考", m.TableCfg.TableId, _seatID)
			m.RemoveThinker(_seatID)
			i-- // 如有元素删除,下标 -1
		}
	}

	logs.Info("tableId:%v-------ThinkQiangGangHu----->seatThink:%v", m.TableCfg.TableId, seatThink)
	//只要有一个位置思考,则返回true
	for _, v := range seatThink {
		if v {
			return true
		}
	}

	return false
}

//单局计算
func (t *BBTable) CalSingle() {

	//无人胡牌,不计算分数,流局 -----------------------------------------------------------------------------
	if t.GetHasHuCt() == 0 {
		t.Majhong.Flow = consts.Yes
	}

	//积分结算--------------------------------------------------------------------------------------------
	for seatID, v := range t.seats {

		if v.GetState() == consts.SeatStateGameHasHu { //胡牌玩家计算胡牌积分
			huPaiRec := t.Majhong.CMajArr[seatID].HuPaiRec
			logs.Info("tableId:%v---------------------积分结算---------------------->HuPaiRec:%v", t.ID, huPaiRec)
			addScore := 0
			for _sID, score := range huPaiRec {
				if _sID != seatID { //被扣积分玩家
					t.Majhong.ChangeCmajScore(_sID, -score)
					addScore += score
				}
			}
			t.Majhong.ChangeCmajScore(seatID, addScore) //胡牌玩家 增加积分

			if t.TableCfg.KaiHuSuanGang == consts.Yes { //开胡算杠
				t.calGfxy(seatID)
			}
		}

		if t.TableCfg.YouGangYouFen == consts.Yes { //有杠有分
			if t.Majhong.Flow != consts.Yes { //非流局
				t.calGfxy(seatID)
			}
		}

	}

}

//玩家碰牌
func (t *BBTable) Peng(_seatId int) {

	//碰牌操作检测过手胡
	t.chkGSH(_seatId)

	_pengCard := t.Majhong.LastSendCard
	logs.Info("tableId:%v,位置:%v碰牌,牌:%v", t.ID, _seatId, _pengCard)
	t.Majhong.CMajArr[_seatId].DoPeng(_pengCard)
	t.Majhong.CMajArr[_seatId].ResetOptInfo() //清空本位置思考操作

	t.Majhong.CMajArr[_seatId].LastFetchMCard = nil // 碰完牌 将LastFetchCard设置为空,防止发送tableinfo 时被移动,导致显示错误

	//碰完等待出牌
	t.Majhong.SetCurtSenderIndex(_seatId)
	t.tableInter.SendCheckTing()
}

//一局游戏开始
func (t *BBTable) GameStart() {

	if t.Majhong.GameCt == 0 {
		t.Majhong.GameCt++ //游戏局数+1 (第一局 游戏开始时增加)
	}

	logs.Info("------------->gameStart(%v %v)<------------- 游戏开始!", t.ID, t.Majhong.GameCt)

	//确定庄家
	t.tableInter.MakeDSeat()

	//确定完庄家后再 清除部分上一局用于下局判断庄家位置的的临时数据
	t.Majhong.ClearLastTmpData()

	//发牌
	t.Majhong.DealCard()

	//发牌状态
	t.setState(consts.TableStateDealCard)
	//
	//if t.TableCfg.DaiHua > 0 {
	//	//time.Sleep(time.Second * 5)
	//	logs.Info("tableId:%v ========================>所有玩家开始补花", t.ID)
	//	//	所有玩家开始轮流补花
	//	t.LunLiuBuHua()
	//	logs.Info("tableId:%v ========================>玩家开始思考", t.ID)
	//}

}

//检测单局游戏结束,牌墙中没有牌则结束游戏,这个方法只在自己拿完最后一张无思考或者有思考但放弃时调用
func (t *BBTable) ChkOver() bool {
	remainPaiCt := t.Majhong.GetRemainPaiCt()
	res := 0
	if t.TableCfg.DaiHua <= 0 {
		res = 14
	}
	if remainPaiCt <= res { // 最后一张有胡,继续思考,超时则游戏结束
		logs.Info("tableId:%v------------------------>[蚌埠仅剩%v张],游戏结束", t.ID, res)
		t.GameOver()
		return true
	}
	return false
}
