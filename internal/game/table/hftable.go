package table

import (
	"qianuuu.com/ahmj/internal/config"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/seat"
	"qianuuu.com/lib/logs"
	. "qianuuu.com/ahmj/internal/mjcomn"
)

// [合肥麻将] 牌桌
type HFTable struct {
	*Table
}

func NewHFTable(_tableID int, _robots *Robots, _tableCfg *config.TableCfg) *HFTable {
	table := NewTable(_tableID, _robots, _tableCfg)
	ret := &HFTable{
		Table: table,
	}
	ret.Table.tableInter = ret
	return ret
}

//玩家胡牌
func (t *HFTable) HuPai(_seat *seat.Seat) {

	seatId := _seat.GetId()
	huCmaj := t.Majhong.CMajArr[seatId]
	huType := huCmaj.HuType //胡牌类型

	//保存所胡的牌
	huCard := NewMCard(huCmaj.OptInfo.HuCard)

	if huType == HUTYPE_JIEPAO { //接炮

		//添加这张牌
		huCmaj.AddHandPai(huCard, false)

		//记录统计数据
		huCmaj.JiePaoCt++                                         //接炮次数
		t.Majhong.CMajArr[t.Majhong.LastSenderSeatID].DianPaoCt++ //点炮次数

		//检测是否一炮多响
		if t.Majhong.YiPaoDXSeatID == consts.DefaultIndex {
			if t.Majhong.LastHuCardData == consts.DefaultIndex {
				t.Majhong.LastHuCardData = huCard.GetData()
			} else {
				if t.Majhong.LastHuCardData == huCard.GetData() { //胡的是同一张,说明是一炮多响
					t.Majhong.YiPaoDXSeatID = t.Majhong.LastSenderSeatID
				}
			}
		}

	} else if huType == HUTYPE_ZIMO { //自摸
		huCmaj.ZimoCt++ //自摸次数
	}

	logs.Info("tableId:%v----------------->位置:%v胡牌,所胡牌:%v,胡牌类型:%v", t.ID, seatId, huCard, GetHuTypeName(huType))

	t.Majhong.CMajArr[seatId].ResetOptInfo() //清空本位置思考操作
	t.Majhong.RemoveThinker(seatId)          //删除思考位置
	_seat.SetState(consts.SeatStateGameHasHu)

	isMultHu := false //是否多人胡牌
	if t.Majhong.GetThinkerCt() > 0 {
		isMultHu = true
	}

	//结算位置胡牌信息
	t.calHuResult(_seat)

	//仍然有胡牌思考者,刷新牌桌,等待其他胡牌玩家思考
	if isMultHu {
		t.SendTableInfo()
		return
	}

	//胡牌,则游戏结束
	logs.Info("tableId:%v----------------->(最后)玩家胡牌,游戏结束!", t.ID)
	t.GameOver()

}

//位置胡牌,计算牌型\积分等
func (t *HFTable) calHuResult(_seat *seat.Seat) {

	huSeatID := _seat.GetId()
	huCmaj := t.Majhong.CMajArr[huSeatID]
	huType := t.Majhong.CMajArr[huSeatID].HuType

	//胡牌类型\牌型描述
	if huType == HUTYPE_JIEPAO {
		huCmaj.AddPxScore("接炮", 0)
	} else {
		huCmaj.AddPxScore("自摸", 0)
	}

	totalZui := 0 //总嘴数

	//胡牌牌型嘴数
	if huCmaj.HasPxId(consts.PXID_HFMJ_TIANHU) { //天胡,不再算其他牌型
		huCmaj.ClearPxId()
		huCmaj.AddPxId(consts.PXID_HFMJ_TIANHU)
	}
	if huCmaj.HasPxId(consts.PXID_HFMJ_DIHU) { //地胡,不再算其他牌型
		huCmaj.ClearPxId()
		huCmaj.AddPxId(consts.PXID_HFMJ_DIHU)
	}

	//牌型嘴数
	for _, v := range huCmaj.HuPxIdArr {
		if v == consts.PXID_HFMJ_GANG_SHANG_HUA {
			continue
		}
		pxName := consts.GetHuPxName_HFMJ(v)
		pxZui := consts.GetHuPxScore_HFMJ(v)
		if v == consts.PXID_HFMJ_TIANHU ||
			v == consts.PXID_HFMJ_DIHU ||
			v == consts.PXID_HFMJ_QYS { //天地胡\清一色嘴数限制
			pxZui = t.TableCfg.TdqZuiCt
		}

		if v == consts.PXID_HFMJ_PINGHU { //牌开
			pxZui = t.TableCfg.BaseScore
			pxName = "牌开"
		}
		huCmaj.AddPxScore(pxName, pxZui)
		totalZui += pxZui
	}

	//另加嘴 牌型分数-----------------------------------TODO ADD
	//t.Check_ExtPxId(huSeatID)
	//for _, v := range huCmaj.ExtPxIdArr {
	//	zuiNum := consts.GetExtPxScore_HFMJ(v)
	//	count := 1
	//	if v == consts.EXTPXID_HFMJ_ZHI {
	//		count = huCmaj.Check_ZHI()
	//	} else if v == consts.EXTPXID_HFMJ_TONG {
	//		count = huCmaj.Check_TONG()
	//	} else if v == consts.EXTPXID_HFMJ_KAN {
	//		count = huCmaj.Check_KAN()
	//	} else if v == consts.EXTPXID_HFMJ_ANGANG {
	//		count = huCmaj.AgCt
	//	} else if v == consts.EXTPXID_HFMJ_SI_HUO {
	//		count = huCmaj.Check_SIHUO()
	//	}
	//	//TODO 还有个别牌型有多个
	//	zuiNum *= count
	//	extPxName := consts.GetExtPxName_HFMJ(v)
	//	huCmaj.AddPxScore(extPxName, zuiNum)
	//	totalZui += zuiNum
	//}

	//杠上开花,总嘴数*2 TODO 测试
	if huCmaj.HasPxId(consts.PXID_HFMJ_GANG_SHANG_HUA) {
		totalZui *= 2
		huCmaj.AddPxScore("杠上开花 x2", 0)
	}

	//应扣分座位id数组
	loseSeatIDArr := make([]int, 0)
	if huType == HUTYPE_ZIMO {
		for i := 0; i < t.TableCfg.PlayerCt; i++ {
			if i != huSeatID {
				loseSeat := t.seats[i]
				if loseSeat.GetState() != consts.SeatStateGameHasHu { //不包括已胡牌的玩家
					loseSeatIDArr = append(loseSeatIDArr, i)
				}
			}
		}
	} else {
		loseSeatIDArr = append(loseSeatIDArr, t.Majhong.LastSenderSeatID)
	}

	scoreAdd := 0 //胡牌玩家赢分

	//庄闲加嘴 ----------------------------------------------------------------------------
	for _, _loseSeatID := range loseSeatIDArr {

		loseScore := totalZui
		loseCmaj := t.Majhong.CMajArr[_loseSeatID]

		if huType == HUTYPE_JIEPAO { //接炮
			if huSeatID == t.Majhong.DSeatID { //庄家接炮,庄数*2嘴
				_zhuangXianZui := huCmaj.LianZhuangCt * 2
				loseScore += _zhuangXianZui
				huCmaj.AddPxScore("庄家接炮", _zhuangXianZui)

			} else { //闲家接炮
				if _loseSeatID == t.Majhong.DSeatID { //庄家放炮,庄数*2嘴
					_zhuangXianZui := loseCmaj.LianZhuangCt * 2
					loseScore += _zhuangXianZui
					huCmaj.AddPxScore("庄家放炮", loseScore)

				}
			}
		} else if huType == HUTYPE_ZIMO { //自摸
			if huSeatID == t.Majhong.DSeatID { //庄家自摸
				_zhuangXianZui := huCmaj.LianZhuangCt * 4
				loseScore += _zhuangXianZui
				huCmaj.AddPxScore("庄家自摸", _zhuangXianZui)
			} else {
				if _loseSeatID == t.Majhong.DSeatID { //闲家自摸
					_zhuangXianZui := loseCmaj.LianZhuangCt * 4
					loseScore += _zhuangXianZui
					huCmaj.AddPxScore("闲家自摸", loseScore)
				}
			}
		}

		t.Majhong.CMajArr[huSeatID].HuPaiRec[_loseSeatID] = loseScore
		scoreAdd += loseScore
	}

	//胡牌玩家总分
	t.Majhong.CMajArr[huSeatID].HuPaiRec[huSeatID] = scoreAdd
}

// 胡牌牌型检测
func (t *HFTable) Check_PxId(_seatId int, _calHuInfo *CalHuInfo) bool {

	_cmaj := t.Majhong.CMajArr[_seatId]

	//清一色
	if _cmaj.Check_QYS() {
		_cmaj.AddPxId(consts.PXID_HFMJ_QYS)
	}

	return false
}

// 胡牌另加嘴检测
func (t *HFTable) Check_ExtPxId(_seatId int) bool {

	//huCmaj := t.Majhong.CMajArr[_seatId]
	//
	//if huCmaj.Check_ZHI() > 0 { //[支]
	//	huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_ZHI)
	//}
	//
	//if huCmaj.Check_KA() { //[卡]
	//	huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_KA)
	//}
	//
	//if huCmaj.Check_QUE_MEN() { //[缺门]
	//	huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_QUE_MENG)
	//}
	//
	//if huCmaj.Check_TONG() > 0 { //[同]
	//	if huCmaj.Check_10TONG() { //［10同］
	//		huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_SHI_TONG)
	//	} else {
	//		huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_TONG)
	//	}
	//}
	//
	//kanCt := huCmaj.Check_KAN()
	//if kanCt > 0 { //[坎]
	//	if kanCt >= 3 {
	//		if huCmaj.Check_3LIAN_KAN() { //［3连坎］
	//			huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_3LIAN_KAN)
	//		}
	//		if huCmaj.Check_4ANKAN() { //［四暗刻］
	//			huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_SI_AN_KE)
	//		}
	//	} else {
	//		huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_KAN)
	//	}
	//}
	//
	//if huCmaj.AgCt > 0 { //[暗杠]
	//	huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_ANGANG)
	//}
	//
	//if huCmaj.Check_SIHUO() > 0 { //[四活]
	//	huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_SI_HUO)
	//}
	//
	//// 双铺子
	//shuanPuArr := huCmaj.Check_SHUANG_PUZI()
	//if len(shuanPuArr) > 0 {
	//	mingCt := 0
	//	anCt := 0
	//	for i := 0; i < len(shuanPuArr); i = i + 6 { //6个为一组
	//		isMing := false
	//		for j := i; j < 6; j++ {
	//			if shuanPuArr[j] == huCmaj.HuMCard.GetData() {
	//				isMing = true
	//				break
	//			}
	//		}
	//		if isMing {
	//			mingCt++
	//		} else {
	//			anCt++
	//		}
	//	}
	//
	//	if anCt == 2 {
	//		huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_2SHUANG_PU_ZI) //［双暗双铺］
	//	}
	//	if anCt == 1 {
	//		huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_AN_SHUANG_PU) //［暗双铺］
	//	}
	//	if mingCt == 1 {
	//		huCmaj.AddExtPxId(consts.EXTPXID_HFMJ_MING_SHUANG_PU) //［明双铺］
	//	}
	//}
	//
	//if len(huCmaj.ExtPxIdArr) > 0 {
	//	return true
	//}
	return false
}

// 确定庄家
func (t *HFTable) MakeDSeat() {

	//第一局为创建房间的玩家为庄
	//之后将若庄家胡牌(或流局)则继续坐庄,庄家不胡则由庄家下家坐庄

	if t.Majhong.GameCt == 1 {
		t.Majhong.DSeatID = 0
		logs.Info("tableId:%v---------------->makeDSeat   房主:%v", t.ID, 0)

	} else {

		if len(t.Majhong.HasHuArr) > 0 { //上局有人胡牌
			//检测是否包含庄家,包含则庄家继续坐庄
			if HasElement(t.Majhong.HasHuArr, t.Majhong.LastZhuangSeatID) {
				logs.Info("tableId:%v------------------->makeDSeat 上局庄家胡牌,庄家继续坐庄 LastZhuangSeatID:%v,t.Majhong.HasHuArr:%v",
					t.ID, t.Majhong.LastZhuangSeatID, t.Majhong.HasHuArr)
				t.Majhong.DSeatID = t.Majhong.LastZhuangSeatID
			} else {
				nextSeatId := t.Majhong.GetNextSeatID(t.Majhong.LastZhuangSeatID)
				logs.Info("tableId:%v------------------->makeDSeat  上局庄家没有胡牌,庄家下家坐庄 位置:%v", t.ID, nextSeatId)
				t.Majhong.DSeatID = nextSeatId
			}
		} else {
			//若荒庄则庄家继续坐庄
			t.Majhong.DSeatID = t.Majhong.LastZhuangSeatID
			logs.Info("tableId:%v------------------->makeDSeat  荒庄 上局庄家:%v", t.ID, t.Majhong.LastZhuangSeatID)
		}

	}

	if t.Majhong.DSeatID < 0 || t.Majhong.DSeatID > t.TableCfg.PlayerCt {
		logs.Info("tableId:%v------------------->***************error makeDSeat t.Majhong.DSeatID:%v", t.ID, t.Majhong.DSeatID)
		t.Majhong.DSeatID = 0
	}

	t.Majhong.CurtSenderIndex = t.Majhong.DSeatID
	t.Majhong.LastZhuangSeatID = t.Majhong.DSeatID

	//庄家连庄+1,并清零其他位置连庄数据
	for _, v := range t.seats {
		_cmaj := t.Majhong.CMajArr[v.Id]
		if v.Id == t.Majhong.DSeatID {
			_cmaj.LianZhuangCt++
		} else {
			_cmaj.LianZhuangCt = 0
		}
	}

	logs.Info("tableId:%v-------------------> makeDSeat    m.DUser:%v", t.ID, t.Majhong.DSeatID)
}

//拿牌者拿牌后开始思考 _specialPxId:某些状态下思考 带入牌型
func (t *HFTable) FetcherStartThink(_specialPxId int) {

	//检测自己思考
	isSelfThink := t.ThinkSelfPai(_specialPxId, false)

	//需要思考
	if isSelfThink {
		t.setState(consts.TableStateWaiteThink)

	} else {
		//如果是最后4张牌,则连续发出,不用等待出牌
		remainCt := t.Majhong.GetRemainPaiCt()

		if remainCt == 0 { //最后一个人不胡牌 游戏结束,不再出牌
			t.GameOver()
			return
		}

		if remainCt < 4 {
			logs.Info("tableId%v---------------->海捞区拿牌无思考,下一个拿牌 remainCt:%v", t.ID, remainCt)
			t.playerFetchPai(true)
			t.FetcherStartThink(consts.PXID_HFMJ_HAI_DI_LAO_YUE) //海底捞月 (合肥麻将 最后4张为海捞区)
			return
		}

		//没有思考,等待出牌
		t.tableInter.SendCheckTing()
	}
}

//拿到牌后思考自己的牌
func (t *HFTable) ThinkSelfPai(_specialPxId int, ispeng bool) bool {

	m := t.Majhong
	m.ClearThinker()
	seatId := m.CurtSenderIndex
	m.AddThinker(seatId)
	curtCmaj := m.CMajArr[seatId]
	curtCmaj.ResetOptInfo() //清空操作
	logs.Info("tableId:%v---------------------->位置:%v开始思考", m.TableCfg.TableId, m.CurtSenderIndex)

	// 检测hu牌 --------------------
	isHu := t.CheckHuPai(seatId)
	if isHu {

		//天胡\杠后开花
		if _specialPxId != consts.DefaultIndex {
			curtCmaj.AddPxId(_specialPxId)
		}

		curtCmaj.SetOpt(OptTypeHu)    //设置胡操作
		curtCmaj.HuType = HUTYPE_ZIMO //胡牌类型为自摸

		//记录位置所胡的牌(天胡也默认为最后拿到的那张)
		huCard := curtCmaj.LastFetchMCard.Clone()
		curtCmaj.OptInfo.HuCard = huCard.GetData()
		curtCmaj.HuMCard = NewMCard(huCard.GetData())

		logs.Info("tableId:%v---------------------->[位置%v]可胡牌,所胡牌:%v,牌型:%v ",
			m.TableCfg.TableId, m.CurtSenderIndex, huCard.String(), curtCmaj.HuPxIdArr)
	}

	// 检测暗杠--------------------
	isAnGang := false
	if m.GetRemainPaiCt() > 4 { //合肥麻将 牌墙中必须大于4张时才能思考暗杠(最后4张是海捞牌)
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

	return isHu || isAnGang
}

//出牌后其他位置玩家思考
func (t *HFTable) ThinkOtherPai(_tkCard *MCard) bool {

	m := t.Majhong

	seatThink := make([]bool, 4)

	logs.Info("tableId:%v,---------------->m.CurtThinkerArr:%v", m.TableCfg.TableId, m.CurtThinkerArr)
	//按位置顺序检测思考
	for i := 0; i < len(m.CurtThinkerArr); i++ {

		_seatID := m.CurtThinkerArr[i]
		_tkCmaj := m.CMajArr[_seatID]
		_tkCmaj.ResetOptInfo() //清空操作

		// 检测hu牌 --------------------
		isHu := false

		//自摸胡不带点炮
		if t.TableCfg.ZimoHu == consts.Yes {
			logs.Info("tableId:%v,---------------->自摸胡,不带点炮", m.TableCfg.TableId)
		} else {
			_tkCmaj.AddHandPai(_tkCard, false) //先添加这张牌
			isHu = t.CheckHuPai(_seatID)

			if isHu {
				_tkCmaj.HuType = HUTYPE_JIEPAO

				//思考他牌时检测过手胡
				curtFanCt := m.GetSeatFanCt(_seatID, false)
				isGuoShouHu := m.CMajArr[_seatID].CheckGSH(curtFanCt)
				if isGuoShouHu {
					logs.Info("tableId:%v------------------->[位置%v]过手胡", _seatID, m.TableCfg.TableId)
					isHu = false //将胡置为false
				} else {
					//检测地胡 只有庄家拿了一张牌,其他人没拿牌
					if m.TableCfg.TiandiHu == consts.Yes {
						//添加座位号非庄家判断 (庄家打出一张牌被其他人碰,然后碰的出一张又被庄家胡牌,导致判断地胡)
						//最近一个出牌人为庄家,非闲家出牌,这个判断已经包含了上面的情况
						if m.PaiIndex == (13*m.TableCfg.PlayerCt+1) &&
							_seatID != m.DSeatID && m.LastSenderSeatID == m.DSeatID {
							logs.Info("tableId:%v------------------->牌型为 [地胡]", m.TableCfg.TableId)
							_tkCmaj.AddPxId(consts.PXID_HFMJ_DIHU)
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
		}

		seatThink[_seatID] = isHu
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
func (t *HFTable) CheckHuPai(_seatId int) bool {

	if !config.Opts().OpenHu {
		return false
	}

	_cmaj := t.Majhong.CMajArr[_seatId]
	//清空牌型记录
	_cmaj.ClearPxId()
	_cmaj.ClearExtPxId()

	//检测前先排序手牌
	_cmaj.SortHandPai()

	//合肥麻将检测8支
	if !_cmaj.Check_8Zhi() {
		return false
	}

	//检测胡牌
	_chuobj := t.GetCalHuObject(_seatId)
	_calHuInfo := CalHuPai(_chuobj.MjCfg, _chuobj.HandPaiArr)

	if _calHuInfo.PxType > PXTYPE_UNKNOW {

		if _calHuInfo.PxType == PXTYPE_PINGHU {
			_cmaj.AddPxId(consts.PXID_HZMJ_PINGHU)
			logs.Info("tableId:%v--------CheckHuPai---------->平胡胡牌", t.TableCfg.TableId)

		} else if _calHuInfo.PxType == PXTYPE_7DUI { //七对

			logs.Info("tableId:%v--------CheckHuPai---------->七对胡牌", t.TableCfg.TableId)
			if _calHuInfo.Check_HAOHUA7DUI() { //豪华七对
				if _calHuInfo.Check_SUPERHAOHUA7DUI() { //超豪华七对
					_cmaj.AddPxId(consts.PXID_HFMJ_SUPERHAOHUA7DUI)
				} else {
					_cmaj.AddPxId(consts.PXID_HFMJ_HAOHUA7DUI)
				}
			} else {
				_cmaj.AddPxId(consts.PXID_HZMJ_7DUI) //普通七对
			}
		}

		if len(_cmaj.HuPxIdArr) > 0 {
			//检测胡牌牌型
			t.Check_PxId(_seatId, _calHuInfo)

			//其他牌型检测
			t.Check_ExtPxId(_seatId)

			return true
		}
	}

	return false
}

//获取胡牌运算对象
func (t *HFTable) GetCalHuObject(_seatId int) *CalHuObject {

	_mjcfg := NewMjCfg()
	_mjcfg.TableType = t.TableCfg.TableType
	_mjcfg.MaxColorCt = 3

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

// 暗杠检测
func (t *HFTable) Check_AnGang(_seatId int) bool {

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
func (t *HFTable) LastHuThinkerCancer() {

	logs.Info("tableId:%v-------------------> 最后一个可胡玩家选择取消胡牌")

	//游戏结束
	t.GameOver()
}

//庄家第一次思考
func (t *HFTable) DPlayerFirstThink() {

	//庄家开始思考 //检测天胡
	_pxId := consts.DefaultIndex
	if t.TableCfg.TiandiHu == consts.Yes {
		if t.Majhong.PaiIndex == (13*t.TableCfg.PlayerCt + 1) { //庄家起手第一张胡牌
			_pxId = consts.PXID_HFMJ_TIANHU
		}
	}
	t.FetcherStartThink(_pxId)
}

//等待出牌->听牌检测
func (t *HFTable) SendCheckTing() {
	//_seatId := t.Majhong.CurtSenderIndex
	//chkCmaj := t.Majhong.CMajArr[_seatId]
	//
	//chkCmaj.SendTipArr = make([]*maj.SendTip, 0)
	////获得手牌
	//allHandPai := chkCmaj.GetHandPai()
	//for _, sendCard := range allHandPai {
	//	//手牌中有红中 当检测到红中时 跳过
	//	if sendCard.IsHongZhong() {
	//		continue
	//	}
	//	//先删除这张牌
	//	chkCmaj.RemoveHandCard(sendCard)
	//	//能胡的牌
	//	huCards := make([]int, 0)
	//	dataArr := []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	//	for _, data := range dataArr {
	//		_saveCard := NewMCard(data)
	//		//将牌添加到手牌中
	//		chkCmaj.AddHandPai(_saveCard, false)
	//		//检测胡牌
	//		isHu := t.CheckHuPai(_seatId)
	//		//删除这张牌
	//		chkCmaj.RemoveHandCard(_saveCard)
	//		if isHu {
	//			huCards = append(huCards, data) //记录这张牌
	//			//TODO 记录胡牌对应的分数
	//		}
	//	}
	//	if len(huCards) > 0 {
	//		//过滤打出去重复的牌
	//		if sendCard.GetData()%4 == 0 {
	//			sendTip := maj.NewSendTip()
	//			//打出去能听得牌
	//			sendTip.SendCard = sendCard.GetData()
	//			for i := 0; i < len(huCards); i++ {
	//				sendTip.HuCards = append(sendTip.HuCards, huCards[i])
	//			}
	//			//添加一个听牌提示
	//			chkCmaj.SendTipArr = append(chkCmaj.SendTipArr, sendTip)
	//		}
	//	}
	//	//添加上检测删除的牌
	//	chkCmaj.AddHandPai(sendCard, false)
	//}
	//for _, v := range chkCmaj.SendTipArr {
	//	sendCard := NewMCard(v.SendCard)
	//	huCard := make([]*MCard, 0)
	//	for _, vv := range v.HuCards {
	//		huCard = append(huCard, NewMCard(vv))
	//	}
	//	logs.Info("tableId:%v------------------SendCheckTing()-------------->打出%v可胡:%v", t.ID, sendCard, huCard)
	//}

	t.setState(consts.TableStateWaiteSend)
}

// 更新位置当前听牌
func (t *HFTable) UpdateTingCards(_seatId int) {

	_cmaj := t.Majhong.CMajArr[_seatId]
	if len(_cmaj.SendTipArr) > 0 { //之前检测有听牌
		lastSendCard := _cmaj.LastSendMCard
		//检测上一次打出的牌是否在SendTipArr中
		if lastSendCard != nil {
			for _, v := range _cmaj.SendTipArr {
				if v.SendCard == lastSendCard.GetData() {
					//更新听牌
					_cmaj.TingCards = make([]int, 0)
					for _, n := range v.HuCards {
						_cmaj.TingCards = append(_cmaj.TingCards, n)
					}
					break
				}
			}
		}
	}
}

//玩家杠牌
func (t *HFTable) Gang(_seatId int, _data int, _isSave bool) {

	_cmaj := t.Majhong.CMajArr[_seatId]
	_gangCard := NewMCard(_data)
	_cmaj.CancerGSH() //杠牌取消过手胡

	//logs.Info("tableId:%v,--------------->_gangType:%v,_isSave:%v ", t.ID, _gangType, _isSave)
	_gangType := t.getGangType(_seatId, _data)
	t.ExecOptInfo.OptDetail = _gangType

	//面杠时判断抢杠胡
	if _gangType == GANGTYPE_MIAN {
		if !_isSave {
			if t.TableCfg.QiangGang == consts.Yes { //可抢杠胡
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
		_dianGangSeatID := t.Majhong.LastSenderSeatID

		_score := 3 //点杠3分
		t.GuaFXiay_Z(_seatId, _score, _dianGangSeatID)
		_cmaj.ZhiGCt++
		_cmaj.MingGangCt++ //明杠次数

		//记录点杠玩家
		dianGangSeatID := t.Majhong.LastSenderSeatID
		t.Majhong.CMajArr[dianGangSeatID].DgCt++

		t.Majhong.LastDianZhiGangSeatID = dianGangSeatID

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

	_extPxId := consts.DefaultIndex
	if t.TableCfg.TableType == TableType_HFMJ {
		_extPxId = consts.PXID_HFMJ_GANG_SHANG_HUA //带入杠上花
	}
	t.tableInter.FetcherStartThink(_extPxId)

}

//玩家面杠操作后,查看其他位置是否抢杠胡 _gangSeatID:面杠操作的位置id
func (t *HFTable) ThinkQiangGangHu(_gangSeatID int, _tkCard *MCard) bool {

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
		if isHu {
			_tkCmaj.SetOpt(OptTypeHu)
			_tkCmaj.HuType = HUTYPE_JIEPAO
			_tkCmaj.HuTypeDetail = HUTYPE_DETAIL_QIANGGANG
			_tkCmaj.OptInfo.HuCard = _tkCard.GetData()
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
func (t *HFTable) CalSingle() {

	//无人胡牌,不计算分数,流局 -----------------------------------------------------------------------------
	if t.GetHasHuCt() == 0 {
		t.Majhong.Flow = consts.Yes
	}

	//积分结算------------------------------------------------------------------------------------------
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
		}

		t.calGfxy(seatID) //计算位置刮风下雨积分
	}
}

//玩家碰牌
func (t *HFTable) Peng(_seatId int) {

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
func (t *HFTable) GameStart() {

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

	////发完牌后翻开最后一张牌(从牌墙中删除)
	//t.Majhong.LastFanMCard = t.Majhong.MCards[len(t.Majhong.MCards)-1].Clone()

	//_tempMCards := make([]*MCard, 0)
	//for i := 0; i < len(t.Majhong.MCards)-1; i++ { //删除最后一张
	//	_tempMCards = append(_tempMCards, t.Majhong.MCards[i])
	//}
	//
	//t.Majhong.MCards = make([]*MCard, 0) //重新添加
	//for i := 0; i < len(_tempMCards); i++ {
	//	t.Majhong.MCards = append(t.Majhong.MCards, _tempMCards[i])
	//}
	//
	//logs.Info("tableId:%v------------->gameStart 翻开最后一张牌:%v,len(t.Majhong.MCards):%v",
	//	t.ID, t.Majhong.LastFanMCard, len(t.Majhong.MCards))

	//发牌状态
	t.setState(consts.TableStateDealCard)
}

//检测单局游戏结束,牌墙中没有牌则结束游戏,这个方法只在自己拿完最后一张无思考或者有思考但放弃时调用
func (t *HFTable) ChkOver() bool {
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
