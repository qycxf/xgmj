// 一局牌

package maj

import (
	"strings"

	"qianuuu.com/lib/logs"
	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/consts"
	. "qianuuu.com/xgmj/internal/mjcomn"
)

// Majong 一局牌管理
type Majong struct {
	TableCfg *config.TableCfg

	MCards          MCards  //一副麻将牌
	DSeatID         int     //庄家位置
	PaiIndex        int     //发牌下标
	CurtSenderIndex int     //当前等待出牌人下标
	CurtThinkerArr  []int   //当前思考人下标数组,1张牌最大3个玩家同时思考
	CMajArr         []*CMaj //座位麻将牌管理

	LastSenderSeatID int    //最近一个出牌玩家位置
	LastSendCard     *MCard //最近一个玩家打出的牌
	LastFetchSeatID  int    //最近一个拿牌玩家位置
	LastFetchMCard   *MCard //最近一个拿牌玩家拿到的牌

	NiaoCards     []*MCard  //抓鸟牌
	ZhuaNiaoInfo  *ZhuaNiao //抓鸟对象数据
	LaiZiCardData int       //癞子牌

	LastMianGangCard      *MCard //最近一个面杠玩家所杠的牌(暂只保存面杠的牌,用于抢杠胡计算)
	LastMianGangSeatID    int    //最近一个面杠牌玩家位置
	LastDianZhiGangSeatID int    //最近一个点(直)杠牌玩家位置

	LastDoGangSeatID int //最近一个杠牌玩家位置 (暗杠\直杠\弯杠)
	LastDoGangType   int //最近一个杠牌玩家杠牌类型 (暗杠\直杠\弯杠)
	LastDoGangData   int //最近一个杠牌值

	FirstHuSeatID    int   //第一个胡牌玩家位置
	HasHuArr         []int //记录所有胡牌玩家位置(自摸或(多个)接炮)
	LastZhuangSeatID int   //上局庄家位置

	YiPaoDXSeatID   int //如果本局有一炮多响,标记一炮多响的放炮位置
	LastHuCardData  int //最近一个胡牌牌值,用于判断是否 一炮多响
	FatchLastSeatID int //拿最后一张牌的位置id
	HaiDi           bool

	Flow     int //是否流局
	GameCt   int //已玩局数
	QuanFeng int //圈风
	HuSeq    int //记录胡牌顺序(包括接炮\自摸)

}

// 创建一局牌
func NewMajong(_tableCfg *config.TableCfg) *Majong {
	logs.Info("tableId:%v--------------------------->NewMajong() ", _tableCfg.TableId)
	ret := Majong{
		TableCfg: _tableCfg,
	}

	ret.MCards = NewMCards(_tableCfg.TableType)
	ret.CMajArr = make([]*CMaj, 0)

	for seatId := 0; seatId < _tableCfg.PlayerCt; seatId++ {
		cmaj := NewCMaj(seatId, ret.TableCfg)
		ret.CMajArr = append(ret.CMajArr, cmaj)
	}
	if !config.Opts().CloseTestPx {
		ret.ResetTestData()
	}

	ret.Init()

	return &ret
}

// 牌局开始,初始化数据
func (m *Majong) Init() {

	m.PaiIndex = 0
	m.CurtSenderIndex = consts.DefaultIndex
	m.DSeatID = consts.DefaultIndex
	m.CurtThinkerArr = make([]int, 0)
	m.Flow = consts.No
	m.LastSenderSeatID = consts.DefaultIndex
	m.LastSendCard = NewMCard(consts.DefaultIndex)
	m.LastFetchSeatID = consts.DefaultIndex
	m.LastFetchMCard = NewMCard(consts.DefaultIndex)
	m.LastMianGangSeatID = consts.DefaultIndex
	m.LastMianGangCard = NewMCard(consts.DefaultIndex)
	m.LastDianZhiGangSeatID = consts.DefaultIndex
	m.LastDoGangSeatID = consts.DefaultIndex
	m.HasHuArr = make([]int, 0)
	m.LastDoGangType = consts.DefaultIndex
	m.LastDoGangData = consts.DefaultIndex
	m.GameCt = 0
	m.HuSeq = 0
	m.ZhuaNiaoInfo = NewZhuaNiao()
	m.LaiZiCardData = consts.DefaultIndex
	m.HaiDi = false
	logs.Info("-------------------11111------->m.LaiZiCardData:%v", m.LaiZiCardData)

}

// 一局游戏结束后(单局),重置相关游戏数据
func (m *Majong) Reset() {

	//重置本局牌局信息
	m.MCards = NewMCards(m.TableCfg.TableType) //重要,换三张会改变 m.MCards中的牌

	m.PaiIndex = 0
	m.CurtSenderIndex = consts.DefaultIndex
	m.CurtThinkerArr = make([]int, 0)
	m.Flow = consts.No

	m.LastSenderSeatID = consts.DefaultIndex
	m.LastSendCard = NewMCard(consts.DefaultIndex)
	m.LastFetchSeatID = consts.DefaultIndex
	m.LastFetchMCard = NewMCard(consts.DefaultIndex)
	m.LastMianGangSeatID = consts.DefaultIndex
	m.LastMianGangCard = NewMCard(consts.DefaultIndex)
	m.LastDianZhiGangSeatID = consts.DefaultIndex
	m.LastDoGangSeatID = consts.DefaultIndex
	m.LastDoGangType = consts.DefaultIndex
	m.LastDoGangData = consts.DefaultIndex
	m.ZhuaNiaoInfo = NewZhuaNiao()
	m.LaiZiCardData = consts.DefaultIndex

	if !config.Opts().CloseTestPx {
		m.ResetTestData()
	}

	//重置位置游戏信息: 番型\积分数据
	for _, v := range m.CMajArr {
		v.Reset()
	}
	m.HuSeq = 0

}

// 确定完庄家之后清除 上局保存的庄家数据
func (m *Majong) ClearLastTmpData() {

	logs.Info("tableId:%v-------------------------------------------------> ClearLastTmpData() ", m.TableCfg.TableId)

	m.HasHuArr = make([]int, 0)
	m.FirstHuSeatID = consts.DefaultIndex
	m.YiPaoDXSeatID = consts.DefaultIndex
	m.LastHuCardData = consts.DefaultIndex
	m.FatchLastSeatID = consts.DefaultIndex
}

//执行抓鸟
func (m *Majong) DoZhuaNiao(_zhuaSeatId int) {

	logs.Info("tableId:%v 执行抓鸟-------------------->_zhuaSeatId:%v", m.TableCfg.TableId, _zhuaSeatId)
	m.NiaoCards = make([]*MCard, 0)

	niaoCt := m.TableCfg.ZhuaNiaoCt

	if m.TableCfg.YiMaQuanZh == consts.Yes {
		niaoCt = 1 //一码全中
	}

	//先从牌墙中取
	_actAddCt := 0
	for i := 0; i < niaoCt; i++ {
		if m.PaiIndex < len(m.MCards) {
			card := NewMCard(-1)
			if !config.Opts().CloseTestPx {
				card = m.GetFromReamin()
			} else {
				card = m.MCards[m.PaiIndex]
			}
			m.PaiIndex++
			m.NiaoCards = append(m.NiaoCards, card)
			logs.Info("tableId:%v 获取鸟牌-------------------->:%v 剩余:%v张", m.TableCfg.TableId, card.String(), m.GetRemainPaiCt())
			_actAddCt++
		}
	}

	//牌墙中鸟牌不够
	if niaoCt > _actAddCt {
		//用最后一张补足
		reAddCt := niaoCt - _actAddCt
		for i := 0; i < reAddCt; i++ {
			lastCard := m.MCards[len(m.MCards)-1]
			m.NiaoCards = append(m.NiaoCards, lastCard)
			logs.Info("tableId:%v 获取鸟牌-------------------->:%v 剩余:%v张", m.TableCfg.TableId, lastCard.String(), m.GetRemainPaiCt())
		}
	}

	m.ZhuaNiaoInfo.ZhuaSeatId = _zhuaSeatId //保存抓鸟位置
	//存储抓鸟对象信息
	for _, v := range m.NiaoCards {
		m.ZhuaNiaoInfo.NiaoCardArr = append(m.ZhuaNiaoInfo.NiaoCardArr, v.GetData())
		//value := v.GetValue()
		//s1 := _zhuaSeatId
		//s2 := m.GetNextSeatID(s1)
		//s3 := m.GetNextSeatID(s2)
		//s4 := m.GetNextSeatID(s3)
		//seatId := consts.DefaultIndex
		//if value == 1 || value == 5 || value == 9 || v.IsHongZhong() { //红中算庄家
		//	seatId = s1
		//} else if value == 2 || value == 6 {
		//	seatId = s2
		//} else if value == 3 || value == 7 {
		//	seatId = s3
		//} else if value == 4 || value == 8 {
		//	seatId = s4
		//}
		//if seatId != consts.DefaultIndex {
		//	m.ZhuaNiaoInfo.ZhongNiaoArr[seatId]++
		//}
	}
	////打印中鸟信息
	//for i := 0; i < len(m.ZhuaNiaoInfo.ZhongNiaoArr); i++ {
	//	logs.Info("tableId:%v -----------> 位置%v中鸟数:%v ", m.TableCfg.TableId, i, m.ZhuaNiaoInfo.ZhongNiaoArr[i])
	//}

}

//得到位置胡牌后 失败玩家位置Id数组
func (m *Majong) GetLoseSeatIdArr(huSeatID int) []int {

	huCmaj := m.CMajArr[huSeatID]
	huType := huCmaj.HuType //胡牌类型

	loseSeatIdArr := make([]int, 0)

	if huType == HUTYPE_JIEPAO { //接炮

		actLoseSeatId := consts.DefaultIndex //计算实际放炮位置Id

		if huCmaj.HuTypeDetail == HUTYPE_DETAIL_QIANGGANG { //抢杠胡
			actLoseSeatId = m.LastMianGangSeatID

		} else { //普通胡牌
			actLoseSeatId = m.LastSenderSeatID
		}

		loseSeatIdArr = append(loseSeatIdArr, actLoseSeatId)

	} else if huType == HUTYPE_ZIMO { //自摸
		for i := 0; i < m.TableCfg.PlayerCt; i++ {
			if i != huSeatID {
				loseSeatIdArr = append(loseSeatIdArr, i)
			}
		}
	}

	return loseSeatIdArr

}

//计算位置成牌总番数 (牌型番+额外加番)
func (m *Majong) GetSeatFanCt(_seatID int, addpxstr bool) int {

	//计算胡牌牌型番数
	//huPxID := m.CMajArr[_seatID].HuPxId
	//fanCt := m.GetPxFanCt(huPxID)

	//if huPxID == consts.PXID_DDHU && m.TableCfg.TableType == consts.TableTypeFour { //内江麻将  对对胡 2番
	//	fanCt = 2
	//}

	//计算额外牌型番数
	//extFanCt := m.GetExtPx(_seatID, addpxstr)
	//fanCt += extFanCt //添加到总番数
	//
	////是否自摸加番
	////if m.CMajArr[_seatID].HuType == consts.HUTYPE_ZIMO && m.TableCfg.ZimoFan == consts.Yes {
	////	fanCt += 1 //加一番
	////}
	//////最大番数限制
	////if fanCt > m.TableCfg.MaxFanCt {
	////	fanCt = m.TableCfg.MaxFanCt
	////}
	//logs.Info("tableId:%v, ------------------->>> GetSeatFanCt:%v", m.TableCfg.TableId, fanCt)
	//return fanCt

	return 0
}

//得到位置所胡的牌对象
func (m *Majong) GetSeatHuCard(huSeatID int) *MCard {

	huCard := NewMCard(-1)

	//自摸
	if m.CMajArr[huSeatID].HuType == HUTYPE_ZIMO {

		huCard = m.CMajArr[huSeatID].LastFetchMCard.Clone()

		//接炮
	} else if m.CMajArr[huSeatID].HuType == HUTYPE_JIEPAO {

		if m.CMajArr[huSeatID].HuTypeDetail == HUTYPE_DETAIL_QIANGGANG { //抢杠胡
			huCard = m.LastMianGangCard.Clone()

		} else if m.CMajArr[huSeatID].HuTypeDetail == HUTYPE_DETAIL_GANG_SHANG_PAO { //点杠花-接炮
			huCard = m.CMajArr[huSeatID].LastFetchMCard.Clone()

		} else { //普通胡牌
			huCard = m.LastSendCard.Clone()
		}
	}

	if huCard.GetData() == -1 {
		logs.Custom(logs.TableTag, "tableId:%v,********* error! GetSeatHuCard,huSeatID:%v,HuType:%v",
			m.TableCfg.TableId, huSeatID, m.CMajArr[huSeatID].HuType)
	}

	return huCard
}

// 发牌
func (m *Majong) DealCard() {

	for i := 0; i < 5000; i++ {
		m.MCards.Shuffle()
	}

	logs.Info("tableId:%v------->洗牌,开始发牌--m.TableCfg.TestSeatId:%v------>m.MCards:%v", m.TableCfg.TableId, m.TableCfg.TestSeatId, m.MCards)

	seatIdArr := make([]int, 0)
	for i := 0; i < m.TableCfg.PlayerCt; i++ {
		seatIdArr = append(seatIdArr, i)
	}
	testSeatId := m.TableCfg.TestSeatId
	if testSeatId != consts.DefaultIndex {

		//将testSeatId 放到最前面
		tmpIdArr := make([]int, 0)
		tmpIdArr = append(tmpIdArr, testSeatId)
		for i := 0; i < len(seatIdArr); i++ {
			if seatIdArr[i] != testSeatId {
				tmpIdArr = append(tmpIdArr, seatIdArr[i])
			}
		}

		logs.Info("tableId:%v~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>tmpIdArr:%v ", m.TableCfg.TableId, tmpIdArr)
		//重新赋值
		seatIdArr = make([]int, 0)
		for i := 0; i < len(tmpIdArr); i++ {
			seatIdArr = append(seatIdArr, tmpIdArr[i])
		}
	}

	logs.Info("tableId:%v~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>seatIdArr:%v ", m.TableCfg.TableId, seatIdArr)
	for i := 0; i < len(seatIdArr); i++ {
		seatID := seatIdArr[i]
		cards := make(MCards, 0)
		//每人13张牌
		if testSeatId != consts.DefaultIndex &&
			seatID == testSeatId &&
			config.TestData().HasGameCt(m.GameCt) &&
			m.TableCfg.TableType == TableType_HZMJ {

			testc := m.FetchTest()
			for _, v := range testc {
				cards = append(cards, v)
			}
		} else {
			for j := 0; j < 13; j++ {
				cards = append(cards, m.FetchACard(seatID, false))
			}
		}
		//logs.Info("庄家：%v", m.DSeatID)
		//庄家14张
		if seatID == m.DSeatID {
			cards = append(cards, m.FetchACard(seatID, false))
			logs.Info("庄家牌：%v", cards)
		}
		m.CMajArr[seatID].SetHandPai(cards)
	}
	for i := 0; i < m.TableCfg.PlayerCt; i++ {
		handpai := m.CMajArr[i].HandPaiArr
		logs.Info("座位%v的手牌为：%v", i, handpai)
	}
}

func (m *Majong) GetMenFeng(seatId int) int {

	if seatId == m.DSeatID {
		return m.QuanFeng
	}

	nextid1 := m.GetNextSeatID(m.DSeatID)
	nextid2 := m.GetNextSeatID(nextid1)
	nextid3 := m.GetNextSeatID(nextid2)

	if nextid1 == seatId {
		qf := m.GetNextFeng(m.QuanFeng)
		qf1 := m.GetNextFeng(qf)
		return m.GetNextFeng(qf1)

	}
	if nextid2 == seatId {
		qf := m.GetNextFeng(m.QuanFeng)
		return m.GetNextFeng(qf)
	}
	if nextid3 == seatId {
		return m.GetNextFeng(m.QuanFeng)
	}
	return 1
}

func (m *Majong) GetNextFeng(qf int) int {

	feng := qf - 1

	if feng == 0 {
		return 4
	}
	logs.Info("Xxxxxx")
	return feng
}

func (m *Majong) FetchTest() []*MCard {

	cards := make([]*MCard, 0)
	pxstr := config.TestData().GetUseStr()

	pxArr := strings.Split(pxstr, ",")
	pxArr = RandStrArr(pxArr) //乱序 pxstr
	logs.Info("tableId:%v=======================FetchTest============pxstr:%v,pxArr>%v", m.TableCfg.TableId, pxArr, pxstr)

	for _, v := range pxArr {
		//从 m.Cards 中删除
		card := m.DelFromMCardsByName(strings.TrimSpace(v))
		m.PaiIndex++ //换三张需要用到这个值
		cards = append(cards, card)

		logs.Info("tableId:%v,---------->> FetchTest 牌:%v <<----------", m.TableCfg.TableId, card.String())
	}
	//将删去的牌放到 m.cards 前面

	//m.Cards = make([]*Card, 0)
	tmp := make([]*MCard, 0)
	for _, v := range cards {
		tmp = append(tmp, v)
	}
	for _, v := range m.MCards {
		tmp = append(tmp, v)
	}

	m.MCards = make([]*MCard, 0)
	for _, v := range tmp {
		m.MCards = append(m.MCards, v)
	}

	return cards
}

func (m *Majong) DelFromMCardsByName(cardName string) *MCard {
	c := NewMCard(-1)
	tmp := make([]*MCard, 0)
	hasDel := false
	for _, v := range m.MCards {
		if !hasDel {
			if v.String() == cardName {
				hasDel = true
				c = v.Clone()
			} else {
				tmp = append(tmp, v)
			}
		} else {
			tmp = append(tmp, v)
		}
	}
	//重新赋值
	m.MCards = make([]*MCard, 0)
	for _, v := range tmp {
		m.MCards = append(m.MCards, v)
	}
	if c.GetData() < 0 {
		logs.Info("tableId:%v,********** error!!!!! DelFromMCardsByName:[%v]", m.TableCfg.TableId, cardName)
	} else {
		logs.Info("tableId:%v, DelFromMCardsByName:%v", m.TableCfg.TableId, cardName)
	}
	return c
}

//发(从牌墙取出)一张牌,fetchFromLast:是否取牌墙的最后一张
func (m *Majong) FetchACard(seatID int, fetchFromLast bool) *MCard {

	if m.GetRemainPaiCt() == 0 {
		m.FatchLastSeatID = seatID
		logs.Info("tableId ：%v-----------------牌组已经没牌了", m.TableCfg.TableId)
		return nil
	}
	//m.LastSendCard = NewMCard(DefaultIndex)
	if !config.Opts().CloseTestPx {

		card := m.TestFetchACard(seatID)

		m.LastFetchSeatID = seatID
		m.LastFetchMCard = card.Clone()
		m.CMajArr[seatID].LastFetchMCard = card.Clone()
		if m.GetRemainPaiCt() == 0 {
			m.FatchLastSeatID = seatID
		}

		m.PaiIndex++ //换三张需要用到这个值
		logs.Info("tableId:%v,---------->> FetchACard1 牌:%v,剩余%v张 <<----------", m.TableCfg.TableId, card.String(), m.GetRemainPaiCt())
		return card
	}

	logs.Info("tableId:%v,---------->> FetchACard  m.PaiIndex:%v, m.MCards=>%v ", m.TableCfg.TableId, m.PaiIndex, m.MCards)
	card := m.MCards[m.PaiIndex]
	if fetchFromLast { //取牌墙最后一张
		saveCard := card.Clone()
		card = m.MCards[len(m.MCards)-1]     //牌墙最后一张
		m.MCards[len(m.MCards)-1] = saveCard //将当期这张牌放到最后
	}

	m.LastFetchSeatID = seatID
	m.LastFetchMCard = card.Clone()
	m.CMajArr[seatID].LastFetchMCard = card.Clone()
	m.PaiIndex++

	if m.GetRemainPaiCt() == 0 {
		m.FatchLastSeatID = seatID
	}
	logs.Info("tableId:%v,---------->> FetchACard2  牌:%v,剩余%v张 <<----------", m.TableCfg.TableId, card.String(), m.GetRemainPaiCt())
	return card
}

////和牌墙中剩下的牌换
//func (m *Majong) ChangeCard(_intArr []int) []*MCard {
//
//	ranArr := make([]int, 0)    // 随机的三张牌的所在牌墙的下标值
//	rtnArr := make([]*MCard, 0) // 随机的三张牌
//
//	random := rand.New(rand.NewSource(time.Now().UnixNano()))
//	for {
//		ranIndex := random.Intn(56) + 52 //从52~107中取牌   m.PaiIndex当前 = 51
//		has := false
//		for _, v := range ranArr {
//			if v == ranIndex {
//				has = true
//				break
//			}
//		}
//		if !has {
//			ranArr = append(ranArr, ranIndex)
//		}
//		if len(ranArr) == 3 {
//			break
//		}
//	}
//	logs.Info("tableId:%v,-------------------->> ChangeCard,抽取的三张下标为:%v", m.TableCfg.TableId, ranArr)
//	for i := 0; i < 3; i++ {
//		if !config.Opts().CloseTestPx {
//			card := m.FetchACard(3) //这里直接从第三个位置以及剩下的牌取,有问题,测试情况下屏蔽换三张
//			rtnArr = append(rtnArr, card.Clone())
//			TestMCards = append(TestMCards, NewMCard(_intArr[i]))
//
//		} else {
//			index := ranArr[i]
//			rtnArr = append(rtnArr, m.MCards[index].Clone())
//			m.MCards[index] = NewMCard(_intArr[i])
//		}
//	}
//
//	logs.Info("tableId:%v,---------->> ChangeCard,抽取的三张牌为:%v", m.TableCfg.TableId, rtnArr)
//	return rtnArr
//}

//返回牌墙剩余张数
func (m *Majong) GetRemainPaiCt() int {
	return len(m.MCards) - m.PaiIndex
}

//移动当前出牌牌人到下一个
func (m *Majong) MoveSenderIndexToNext() {
	nextSenderIndex := m.GetNextSeatID(m.CurtSenderIndex)
	m.SetCurtSenderIndex(nextSenderIndex) //默认当前拿牌人为当前出牌人
}

//返回 _curtSeatID 下一个位置,无论状态
func (m *Majong) GetNextSeatID(_curtSeatID int) int {
	nextSeatID := _curtSeatID + 1
	if nextSeatID >= m.TableCfg.PlayerCt {
		nextSeatID = 0
	}
	return nextSeatID
}

//设置出牌人位置
func (m *Majong) SetCurtSenderIndex(_senderIndex int) {
	m.CurtSenderIndex = _senderIndex
}

// 是否含有思考者位置
func (m *Majong) HasThinker(_seatID int) bool {
	for _, v := range m.CurtThinkerArr {
		if v == _seatID {
			return true
		}
	}
	return false
}

// 添加思考者位置
func (m *Majong) AddThinker(_thinkSeatID int) {
	//不能重复添加
	for _, v := range m.CurtThinkerArr {
		if v == _thinkSeatID {
			return
		}
	}
	m.CurtThinkerArr = append(m.CurtThinkerArr, _thinkSeatID)
}

//删除思考位置
func (m *Majong) RemoveThinker(_thinkSeatID int) {

	tmpArr := make([]int, 0)
	len1 := len(m.CurtThinkerArr)
	for _, v := range m.CurtThinkerArr {
		if v == _thinkSeatID {

		} else {
			tmpArr = append(tmpArr, v)
		}
	}
	m.CurtThinkerArr = tmpArr
	len2 := len(m.CurtThinkerArr)
	if len2 >= len1 {
		logs.Info("tableId:%v,*********************************** RemoveThinker err! 未找到 _thinkSeatID:%v ", m.TableCfg.TableId, _thinkSeatID)
	}
}

//位置积分变化
func (m *Majong) ChangeCmajScore(_seatID int, _value int) {

	if _value == 0 {
		return
	}

	cmaj := m.CMajArr[_seatID]
	cmaj.Score += _value
	cmaj.TotalScore += _value
	logs.Info("tableId:%v,---------------------------->位置%v积分变化:%v,本局积分:%v,总积分:%v", cmaj.TableCfg.TableId, cmaj.SeatID, _value, cmaj.Score, cmaj.TotalScore)

}

// 清空思考者位置
func (m *Majong) ClearThinker() {
	m.CurtThinkerArr = make([]int, 0)
}

// 返回当前思考者数量
func (m *Majong) GetThinkerCt() int {
	return len(m.CurtThinkerArr)
}

//从 _seatId 右手边查找所有未放弃的玩家(逆时针,不包括 _seatID 本身)
func (m *Majong) SearchArrRight(_seatID int) []int {
	tmpArr := make([]int, 0)
	//从右手 逆时针 开始查找
	for i := _seatID + 1; i < m.TableCfg.PlayerCt; i++ {
		tmpArr = append(tmpArr, i)
	}
	for i := 0; i < _seatID; i++ {
		tmpArr = append(tmpArr, i)
	}
	return tmpArr
}

//从 _seatId 左手边查找所有未放弃的玩家(顺时针)
func (m *Majong) SearchArrLeft(_seatID int) []int {
	tmpArr := make([]int, 0)
	//从左手 顺时针 开始查找
	for i := _seatID - 1; i >= 0; i-- {
		tmpArr = append(tmpArr, i)
	}
	for i := m.TableCfg.PlayerCt - 1; i > _seatID; i-- {
		tmpArr = append(tmpArr, i)
	}
	return tmpArr
}

//牌型名称
func (m *Majong) GetPxName(_pxIndex int) string {
	return [13]string{"平胡", "对对胡", "清一色", "带幺九", "七对", "龙七对", "清对",
		"将对", "将七对", "清七对", "清龙七对", "天胡", "地胡"}[_pxIndex]
}

//牌型番数
func (m *Majong) GetPxFanCt(_pxIndex int) int {
	return []int{0, 1, 2, 3, 2, 2, 3, 3, 4, 4, 4, 5, 5}[_pxIndex]
}

//抓鸟数据对象
type ZhuaNiao struct {
	ZhuaSeatId   int   //抓鸟位置
	NiaoCardArr  []int //n 张鸟牌
	ZhongNiaoArr []int //按座位序号 存储中鸟数

	HuSeatIdArr   []int //胡牌位置 (1~3个)
	LoseSeatIdArr []int //放炮或被自摸位置 (1~3个)
}

func NewZhuaNiao() *ZhuaNiao {
	ret := ZhuaNiao{
		ZhuaSeatId:    consts.DefaultIndex,
		NiaoCardArr:   make([]int, 0),
		ZhongNiaoArr:  make([]int, 4),
		HuSeatIdArr:   make([]int, 0),
		LoseSeatIdArr: make([]int, 0),
	}
	return &ret
}

//用于测试牌型使用  ------------------------------------------------------------------------------------

var TestMCards = NewMCards(3)
var TestpxArr = [][]string{}

func (m *Majong) ResetTestData() {
	//用于测试牌型使用
	TestMCards = NewMCards(m.TableCfg.TableType)
	TestpxArr = [][]string{

		//合肥麻将========================================================================================
		////庄家0 天胡-平胡
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "5万", "6筒", "6筒"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "7万", "7万"},
		//{"2万", "3万", "4万", "5万", "2筒", "3筒", "4筒", "5筒", "2条", "3条", "4条", "5条", "8万", "8万"}}

		////庄家0 暗杠
		//{"2万", "2万", "3万", "3万", "3万", "3万", "4万", "4万", "4万", "4万", "6筒", "6筒", "6筒", "6筒"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "7万", "7万"},
		//{"2万", "3万", "4万", "5万", "2筒", "3筒", "4筒", "5筒", "2条", "3条", "4条", "5条", "8万", "8万"}}

		//庄家0 连续杠
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "6筒", "6筒", "6筒", "6筒", "4万"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "7万", "7万"},
		//{"7筒", "7筒", "7筒", "8筒", "8筒", "8筒", "7条", "7条", "7条", "8条", "8条", "8条", "8万", "8万"}}

		//位置1放炮 位置0胡
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "6筒", "6筒", "6条", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6条"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "8条", "7万", "7万"},
		//{"7筒", "7筒", "7筒", "8筒", "5筒", "8筒", "7条", "7条", "7条", "8条", "8条", "8条", "8万", "8万"}}

		//庄家0 暗杠\胡
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "6筒", "6筒", "6筒"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "7万", "7万"},
		//{"2万", "3万", "4万", "5万", "2筒", "3筒", "4筒", "5筒", "2条", "3条", "4条", "5条", "8万", "8万"}}

		////庄家0 清一色
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "5万", "6万", "8万"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6筒", "6万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万", "7万"}}

		////庄家0 杠上开花
		//{"2万", "2万", "2万", "2万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "5万", "6万", "8万", "7万"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6筒", "6万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万", "7万"}}

		////庄家0 胡牌 支 检测
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "6筒", "6筒", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "6筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万"}}

		////庄家0 胡牌 卡 检测
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "6筒", "8筒", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "7筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万"}}

		////庄家0 胡牌 同 检测
		//{"2万", "2万", "2万","2筒", "2筒", "2筒" , "4万", "4万", "4万", "5万", "5万", "6筒", "8筒", "6条"},
		//{"3万", "3万", "3万", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "7筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万"}}

		////庄家0 胡牌 十同 检测
		//{"2万", "2万", "2万", "2筒", "2筒", "2筒", "2条", "2条", "2条", "2万", "3万", "4万", "8万", "8万"},
		//{"3万", "3万", "3万", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "8筒"},
		//{"4万", "4万", "4万", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "7万"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "7万", "7万"}}

		//庄家0 胡牌 四暗刻 检测
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "5万", "8筒", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "8筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万"}}

		////庄家0 胡牌 暗杠 检测
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "5万", "5万", "8筒", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "8筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万"}}

		////庄家0 胡牌 四活 检测
		//{"2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "4万", "5万", "6万", "8筒", "7万"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "8筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "5条"}}

		////庄家0 胡牌 双暗双铺 检测
		//{"2万", "3万", "4万", "2万", "3万", "4万", "5万", "6万", "7万", "5万", "6万", "7万", "8筒", "7万"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "8筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "5条"}}

		////庄家0 胡牌 暗双铺 检测
		//{"2万", "3万", "4万", "2万", "3万", "4万", "5万", "5万", "5万", "6万", "6万", "6万", "8筒", "7万"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "8筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "5条"}}

		////庄家0 胡牌 明双铺 检测
		//{"2万", "3万", "4万", "2万", "3万", "7筒", "5万", "5万", "5万", "6万", "6万", "6万", "8筒", "8筒"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "4万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "5条"}}

		//庄家0 胡牌 暗双铺+明双铺 检测
		//{"2万", "3万", "4万", "2万", "3万", "7筒", "5万", "6万", "7万", "5万", "6万", "7万", "8筒", "8筒"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5万", "5筒", "6筒", "4万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "5条"}}

		////庄家0 放炮 位置1胡牌
		//{"2万", "3万", "4万", "2万", "3万", "7筒", "5万", "6万", "7万", "5万", "6万", "7万", "8筒"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6筒", "7筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "5条"}}

		//庄家0 七对
		//{"2万", "2万", "3万", "3万", "4万", "4万", "5万", "5万", "6万", "6万", "7万", "7万", "6筒", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万", "7万"}}

		////庄家0 豪华七对
		//{"2万", "2万", "2万", "2万", "4万", "4万", "5万", "5万", "6万", "6万", "7万", "7万", "6筒", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万", "7万"}}

		////庄家0 超豪华七对
		//{"2万", "2万", "2万", "2万", "4万", "4万", "4万", "4万", "6万", "6万", "7万", "7万", "6筒", "6条"},
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万", "6筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "8筒", "8筒"},
		//{"8条", "8条", "8条", "7条", "7条", "7条", "6条", "6条", "6条", "8万", "8万", "8万", "7万", "7万"}}

		//红中麻将========================================================================================

		////庄家0 普通胡牌
		//{"6万", "6万", "6万", "2筒", "2筒", "2筒", "3万", "3万", "3万", "4筒", "4筒", "4筒", "8万", "9万"},
		//{"1筒", "1筒", "1筒", "2万", "2万", "2万", "3筒", "3筒", "3筒", "4万", "4万", "4万", "5筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "6筒", "7筒", "7筒"}}

		////庄家0 七对
		//{"6万", "6万", "7万", "7万", "2万", "2万", "3筒", "3筒", "4筒", "4筒", "6筒", "6筒", "5万", "5万"},
		//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "6筒", "7筒", "7筒"}}

		////庄家0 暗杠\胡
		//{"6万", "6万", "6万", "2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "1万", "6万"},
		//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "6筒", "7筒", "7筒"}}

		////庄家0 清一色
		//{"6万", "6万", "6万", "2万", "2万", "2万", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5万", "5万"},
		//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "6筒", "7筒", "7筒"}}

		//庄家0 暗杠\明杠
		//{"6万", "6万", "6万", "2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "2万"},
		//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "6万"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "6筒", "7筒", "7筒"}}

		////位置1 胡
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万"},
		//{"1万", "1万", "1万", "3万", "3万", "3万", "4万", "4万", "4万", "7万", "8万", "6筒", "6筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "1条", "5条", "5条", "7万"},
		//{}}

		//位置位置1\2胡
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "6万"},
		//{"9万", "9万", "9万", "3万", "3万", "3万", "4万", "4万", "4万", "7万", "8万", "6筒", "6筒", "6万"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "6万"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9条", "9条", "9条", "1条", "1条", "4万", "5万"}}

		////位置1 碰杠胡
		//{"2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒", "5筒", "1万", "6万"},
		//{"1万", "1万", "1万", "3万", "3万", "3万", "4万", "4万", "4万", "7万", "8万", "6筒", "6筒"},
		//{"2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条", "5条", "7万"},
		//{}}

		////位置2碰3胡
		//{"6万", "6万", "6万", "2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "9条"},
		//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "9条", "9条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "7条", "8条"}}

		//位置0 起手4红中
		//{"中", "中", "6万", "2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "中", "中"},
		//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "4条", "5条", "5条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "6筒", "7筒", "7筒"}}

		////位置0 抢杠胡
		//{"6万", "6万", "6万", "2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "3筒", "4筒"},
		//{"9筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "5筒", "5筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "9条", "9条"},
		//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "7条", "8条"}}

		////听牌检测
		//{"9条", "9条", "3筒", "3筒", "4筒","4筒", "5筒", "5筒", "2筒", "3筒", "4筒", "1筒", "1筒","1筒"},
		//{"中", "中", "中", "3万", "4万", "4筒", "6筒", "7筒", "7筒", "1条", "2条", "4条", "5条", "8条", "9条"},
		//{"9筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "5筒", "5筒"},
		//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "南", "1筒","9条"},
		//{ "1条", "1条", "2条", "2条","3条", "3条", "4条", "5条", "6条", "6条", "6条", "9条","9条"}}

		//测试牌,连续 碰-取消  胡取消
		{"中", "中", "中", "2万", "2万", "2万", "3万", "3万", "3万", "4万", "4万", "4万", "5万", "9条", "5万"},
		{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒"},
		{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "8条", "9条"},
		{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "7条", "8条", "4万"}}

	//{"中", "中", "中", "1万", "2万", "2万", "2万", "4筒", "5筒", "6筒", "4条", "4条", "4条", "6万"},
	//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒"},
	//{"1条", "1条", "1条", "2条", "2条", "2条", "3条", "3条", "3条", "4条", "4条", "8条", "9条"},
	//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "9万", "6筒", "6筒", "7条", "8条"}}

	//{"9万", "9万", "9万", "4万", "4万", "4万", "6筒", "7筒", "8筒", "1条", "2条", "3条", "7条", "7条"},
	//{"1筒", "1筒", "1筒", "2筒", "2筒", "2筒", "3筒", "3筒", "3筒", "4筒", "4筒", "4筒", "5筒"},
	//{"1条", "1条", "5条", "5条", "5条", "3条", "3条", "3条", "6条", "6条", "8条", "9条", "8条"},
	//{"7万", "7万", "7万", "8万", "8万", "8万", "9万", "9万", "6筒", "6筒", "7条", "8条", "8条"}}

}

//测试发牌
func (m *Majong) TestFetchACard(_seatID int) *MCard {

	//有自定义顺序的牌则按顺序取出(同一种牌不可超过4张)
	if len(TestpxArr[_seatID]) > 0 {
		tmp := TestpxArr[_seatID]
		t := make([]string, 0)
		cardName := ""
		for i := 0; i < len(tmp); i++ {
			if i == 0 {
				cardName = tmp[0]
			} else {
				t = append(t, tmp[i])
			}
		}
		TestpxArr[_seatID] = t // 重新赋值数组

		card := m.DelByName(cardName)
		return card
	} else {
		//自定义顺序取完则从剩下的牌中取
		return m.GetFromReamin()
	}

}

//根据名称删除元素,删除一个(不知道data,只知道name)
func (m *Majong) DelByName(cardName string) *MCard {
	c := NewMCard(-1)
	tmp := make([]*MCard, 0)
	hasDel := false
	for _, v := range TestMCards {
		//logs.Info("%v", v)
		if !hasDel {
			if v.String() == cardName {
				hasDel = true
				c = v.Clone()
			} else {
				tmp = append(tmp, v)
			}
		} else {
			tmp = append(tmp, v)
		}
	}
	//重新赋值
	TestMCards = make([]*MCard, 0)
	for _, v := range tmp {
		TestMCards = append(TestMCards, v)
	}
	if c.GetData() < 0 {
		logs.Info("tableId:%v,********** error!!!!! DelByName:%v", m.TableCfg.TableId, cardName)
	} else {
		//logs.Info(" DelByName:%v", cardName)
	}
	return c
}

//从  cs MCards 取牌
func (m *Majong) GetFromReamin() *MCard {
	c := NewMCard(-1)
	tmp := make([]*MCard, 0) //这里写的是按顺序取,可随机取
	for i := 0; i < len(TestMCards); i++ {
		if i == 0 {
			c = TestMCards[i]
		} else {
			tmp = append(tmp, TestMCards[i])
		}
	}

	//重新赋值
	TestMCards = make([]*MCard, 0)
	for _, v := range tmp {
		TestMCards = append(TestMCards, v)
	}
	return c
}
