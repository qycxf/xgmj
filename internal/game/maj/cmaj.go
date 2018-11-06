// 玩家手牌管理及算法
package maj

import (
	"sort"

	"strconv"

	"qianuuu.com/lib/logs"
	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/consts"
	. "qianuuu.com/xgmj/internal/mjcomn"
)

// CMaj 单个玩家牌管理对象
type CMaj struct {
	TableCfg *config.TableCfg //桌子配置
	SeatID   int              //座位号

	HandPaiArr [][]*MCard //手牌
	PengPaiArr [][]*MCard //碰牌
	GangPaiArr [][]*MCard //杠牌

	SaveCalHuInfo *CalHuInfo //胡牌的牌型信息（）
	HuPxIdArr     []int      //胡牌牌型IdArr
	ExtPxIdArr    []int      //额外牌型IdArr 额外加番会
	PxScoreArr    [][]string //牌型\分数列表
	OptInfo       *OptInfo   //当前可进行操作对象
	PGCArr        [][]int    //记录碰\杠\吃的牌 (按顺序)
	BuHua         *BuHua     //座位补花信息

	SendTipArr    []*SendTip //出牌提示信息
	TingCards     []int      //如果位置听牌,保存听牌数据
	HuIsBaoPai    bool       //位置本次胡牌是否为 "包牌"
	BaoPaiSeatId  int        //负责"包牌"位置
	BaoPaiScoreCt int        //"包牌"分数倍数

	YinZi          bool     //怀远印子
	IsFristSend    bool     //是否有算天胡
	IsTing         bool     //是否有报听
	HowTingCt      int      //报听出牌的次数
	LianZhuangCt   int      //连庄数
	TempGangPaiArr []*MCard //临时记录可杠的牌
	TempBuHuaArr   []*MCard //临时记录可杠的牌
	HuMCard        *MCard   //记录位置所胡的牌
	HuType         int      //胡牌类型
	HuTypeDetail   int      //胡牌详细类型
	HuSeq          int      //如果胡牌,记录胡牌顺序

	OutPaiArr      []*MCard //打出的牌
	HuaPaiArr      []*MCard //打出的花牌
	LastFetchMCard *MCard   //位置最近一次拿到的牌
	LastSendMCard  *MCard   //位置最近打出的牌
	LastGangType   int      //最近一次杠牌类型
	BaoGangSeatId  int      //阜阳麻将,最后14张放杠,记录放直杠的位置

	// 特殊状态使用-----------------
	GangShangPao    bool   //进入杠上炮状态(杠完后出的第一张牌)
	GuoShouHu       bool   //进入过手胡限制
	GuoShouPeng     bool   //进入过手碰限制
	GuoShouHuCard   *MCard //进入过手碰的那张牌
	GuoShouPengCard *MCard //进入过手碰的那张牌
	GuoShouFanCt    int    //检测过手胡番数
	ChaDaJiaoData   int    //如果结算状态位置听牌并查大叫,则记录这张牌

	//如果杠上炮被胡牌,则记录相关数据,用于检查转雨
	GspZyArr []*GspZy //杠上炮转雨,可能会发生多个

	GfxyArr  []int   //临时保存刮风下雨所得积分,按位置存,本位置存总的,其它位置存扣除的
	GfxyRec  [][]int //记录每次杠操作各个位置被扣积分,等到单局结算时一起结算
	HuPaiRec []int   //记录位置胡牌时 各个位置被扣积分,等到单局结算时一起结算

	//游戏数据统计使用 -----------------
	Score      int //一局游戏积分
	TotalScore int //游戏牌局总积分(游戏实时积分)
	FanCt      int //一局游戏番数
	ZhiGCt     int //一局直杠次数
	MianGCt    int //一局面杠次数
	AgCt       int //一局暗杠次数
	DgCt       int //一局点杠次数
	CdjCt      int //一局查大叫次数

	//最终结算统计 ---------------------
	ZimoCt     int //总自摸次数
	JiePaoCt   int //总接炮次数
	DianPaoCt  int //总点炮次数
	AnGangCt   int //总暗杠次数
	MingGangCt int //总明杠次数
	ChaDaJiao  int //总查大叫次数
	ChaHuaZhu  int //总查花猪次数
}

//新建玩家牌管理对象
func NewCMaj(_seatID int, _tableCfg *config.TableCfg) *CMaj {
	ret := &CMaj{
		SeatID: _seatID,

		SaveCalHuInfo:  nil,
		HandPaiArr:     make([][]*MCard, _tableCfg.MaxCardColorIndex),
		PengPaiArr:     make([][]*MCard, _tableCfg.MaxCardColorIndex),
		GangPaiArr:     make([][]*MCard, _tableCfg.MaxCardColorIndex),
		TempGangPaiArr: make([]*MCard, 0),
		GfxyRec:        make([][]int, 0),

		GangShangPao:  false,
		GuoShouHu:     false,
		GuoShouPeng:   false,
		GuoShouFanCt:  0,
		ChaDaJiaoData: consts.DefaultIndex,

		GspZyArr: make([]*GspZy, 0),

		BuHua:         NewBuHua(),
		SendTipArr:    make([]*SendTip, 0),
		TingCards:     make([]int, 0),
		HuIsBaoPai:    false,
		IsTing:        false,
		IsFristSend:   false,
		YinZi:         false,
		HowTingCt:     0,
		BaoPaiSeatId:  consts.DefaultIndex,
		BaoPaiScoreCt: 1,

		LastFetchMCard: nil,
		LastSendMCard:  nil,

		PGCArr:          make([][]int, 0),
		LianZhuangCt:    0,
		HuPxIdArr:       make([]int, 0),
		ExtPxIdArr:      make([]int, 0),
		PxScoreArr:      make([][]string, 0),
		OutPaiArr:       make([]*MCard, 0),
		HuaPaiArr:       make([]*MCard, 0),
		OptInfo:         NewOptInfo(),
		HuType:          consts.DefaultIndex,
		HuTypeDetail:    consts.DefaultIndex,
		HuMCard:         NewMCard(consts.DefaultIndex),
		GuoShouHuCard:   NewMCard(consts.DefaultIndex),
		GuoShouPengCard: NewMCard(consts.DefaultIndex),
		LastGangType:    consts.DefaultIndex,
		BaoGangSeatId:   consts.DefaultIndex,
		Score:           0,
		FanCt:           0,
		HuSeq:           0,
		ZhiGCt:          0,
		MianGCt:         0,
		AgCt:            0,
		DgCt:            0,
		CdjCt:           0,
	}

	ret.TableCfg = _tableCfg
	_maxPlayerCt := ret.TableCfg.PlayerCt
	ret.GfxyArr = make([]int, _maxPlayerCt)
	ret.HuPaiRec = make([]int, _maxPlayerCt)
	return ret
}

//一局结束后重置数据
func (cmaj *CMaj) Reset() {

	cmaj.SaveCalHuInfo = nil
	cmaj.HandPaiArr = make([][]*MCard, cmaj.TableCfg.MaxCardColorIndex)
	cmaj.PengPaiArr = make([][]*MCard, cmaj.TableCfg.MaxCardColorIndex)
	cmaj.GangPaiArr = make([][]*MCard, cmaj.TableCfg.MaxCardColorIndex)
	cmaj.TempGangPaiArr = make([]*MCard, 0)
	cmaj.GfxyArr = make([]int, cmaj.TableCfg.PlayerCt)
	cmaj.GfxyRec = make([][]int, 0)
	cmaj.HuPaiRec = make([]int, cmaj.TableCfg.PlayerCt)

	cmaj.BuHua = NewBuHua()
	cmaj.PGCArr = make([][]int, 0)
	cmaj.HuPxIdArr = make([]int, 0)
	cmaj.ExtPxIdArr = make([]int, 0)
	cmaj.PxScoreArr = make([][]string, 0)
	cmaj.OutPaiArr = make([]*MCard, 0)
	cmaj.HuaPaiArr = make([]*MCard, 0)
	cmaj.OptInfo = NewOptInfo()
	cmaj.HuType = consts.DefaultIndex
	cmaj.HuTypeDetail = consts.DefaultIndex
	cmaj.LastGangType = consts.DefaultIndex
	cmaj.BaoGangSeatId = consts.DefaultIndex
	cmaj.Score = 0
	cmaj.FanCt = 0

	cmaj.GangShangPao = false
	cmaj.IsTing = false
	cmaj.IsFristSend = false
	cmaj.YinZi = false
	cmaj.HowTingCt = 0
	cmaj.GuoShouHu = false
	cmaj.GuoShouPeng = false
	cmaj.GuoShouHuCard = NewMCard(consts.DefaultIndex)
	cmaj.GuoShouPengCard = NewMCard(consts.DefaultIndex)
	cmaj.GuoShouFanCt = 0
	cmaj.ChaDaJiaoData = consts.DefaultIndex

	cmaj.GspZyArr = make([]*GspZy, 0)

	cmaj.SendTipArr = make([]*SendTip, 0)
	cmaj.TingCards = make([]int, 0)
	cmaj.HuIsBaoPai = false
	cmaj.BaoPaiSeatId = consts.DefaultIndex
	cmaj.BaoPaiScoreCt = 1

	cmaj.LastFetchMCard = nil
	cmaj.LastSendMCard = nil

	cmaj.HuSeq = 0
	cmaj.ZhiGCt = 0
	cmaj.MianGCt = 0
	cmaj.AgCt = 0
	cmaj.DgCt = 0
	cmaj.CdjCt = 0
}

//杠上炮转雨数据对象
type GspZy struct {
	GSPData           int   //杠上炮所胡牌值,用这个作为唯一标识
	GSPIsHu           bool  //杠上炮是否被胡牌
	GSPGangType       int   //杠上炮杠类型
	GSPDianGangSeatId int   //如果杠上炮是直杠,记录点杠位置id
	GSPHuPaiSeatIdArr []int //杠上炮可能会被多个人胡牌,一炮多响
	GSPGfxyRecIndex   int   //记录要转雨的记录的杠数据下标值
}

//新建杠上炮转雨数据对象
func NewGspZy() *GspZy {
	ret := GspZy{
		GSPData:           -1,                  //杠上炮所胡牌值,用这个作为唯一标识
		GSPIsHu:           false,               //杠上炮是否被胡牌
		GSPGangType:       consts.DefaultIndex, //杠上炮杠类型
		GSPDianGangSeatId: consts.DefaultIndex, //如果杠上炮是直杠,记录点杠位置id
		GSPHuPaiSeatIdArr: make([]int, 0),      //杠上炮可能会被多个人胡牌,一炮多响
		GSPGfxyRecIndex:   consts.DefaultIndex,
	}
	return &ret
}

//判断听牌是否含有胡一张的情况
func (cmaj *CMaj) IsHuOne() bool {
	if len(cmaj.SendTipArr) > 0 {
		for _, _sendTip := range cmaj.SendTipArr {
			if len(_sendTip.HuCards) == 1 {
				return true
			}
		}
	}
	return false
}

//查找杠上炮记录数据
func (cmaj *CMaj) GetGspZy(_cardData int) *GspZy {
	for _, v := range cmaj.GspZyArr {
		if v.GSPData == _cardData {
			return v
		}
	}
	return nil
}

//根据杠牌数据下标 查找杠上炮记录数据(一个人可能会多次杠上炮)
func (cmaj *CMaj) GetGspZyByRecIndex(_recIndex int) *GspZy {
	logs.Info("=========================================>cmaj.GspZyArr:%v,_recIndex:%v", cmaj.GspZyArr, _recIndex)
	for _, v := range cmaj.GspZyArr {
		if v.GSPGfxyRecIndex == _recIndex {
			return v
		}
	}
	return nil
}

//位置可接炮时检测 过手胡
func (cmaj *CMaj) CheckGSH(curtFanCt int) bool {
	if !cmaj.GuoShouHu {
		return false
	}
	//if curtFanCt > cmaj.GuoShouFanCt { //当前的番数 > 上次检测到的番数 ,则不限制
	//	return false
	//}
	logs.Info("tableId:%v,-------->位置%v 被检测到过手胡,curtFanCt:%v,cmaj.GuoShouFanCt:%v", cmaj.TableCfg.TableId, cmaj.SeatID, curtFanCt, cmaj.GuoShouFanCt)
	return true
}

//位置可接炮时检测 过手碰
func (cmaj *CMaj) CheckGSP() bool {
	if !cmaj.GuoShouPeng {
		return false
	}
	logs.Info("tableId:%v,-------->位置%v 被检测到过手碰", cmaj.TableCfg.TableId, cmaj.SeatID)
	return true
	//return false
}

//取消过手胡标志
func (cmaj *CMaj) CancerGSH() {
	if cmaj.GuoShouHu {
		logs.Info("tableId:%v,-------->位置%v 取消过手胡", cmaj.TableCfg.TableId, cmaj.SeatID)
		cmaj.GuoShouHu = false
		cmaj.GuoShouHuCard = NewMCard(consts.DefaultIndex)
	}
}

//取消过手碰标志
func (cmaj *CMaj) CancerGSP() {
	if cmaj.GuoShouPeng {
		logs.Info("tableId:%v,-------->位置%v 取消过手碰", cmaj.TableCfg.TableId, cmaj.SeatID)
		cmaj.GuoShouPeng = false
		cmaj.GuoShouPengCard = NewMCard(consts.DefaultIndex)
	}
}

//重置操作数据
func (cmaj *CMaj) ResetOptInfo() {
	cmaj.OptInfo = NewOptInfo()
}

//删除位置某个操作
func (cmaj *CMaj) RemoveOpt(_optType int) {

	if _optType == OptTypePeng {
		cmaj.OptInfo.Peng = false
	} else if _optType == OptTypeGang {
		cmaj.OptInfo.Gang = false
	} else if _optType == OptTypeChi {
		cmaj.OptInfo.Chi = false
	} else if _optType == OptTypeHu {
		cmaj.OptInfo.Hu = false
	} else if _optType == OptTypeBu {
		cmaj.OptInfo.Bu = false
	} else if _optType == OptTypeTing {
		cmaj.OptInfo.Ting = false
	}
}

//设置位置可用操作 0 1 2 放操作
func (cmaj *CMaj) SetOpt(_optType int) {

	if _optType == OptTypePeng {
		cmaj.OptInfo.Peng = true
	} else if _optType == OptTypeGang {
		cmaj.OptInfo.Gang = true
	} else if _optType == OptTypeChi {
		cmaj.OptInfo.Chi = true
	} else if _optType == OptTypeHu {
		cmaj.OptInfo.Hu = true
	} else if _optType == OptTypeBu {
		cmaj.OptInfo.Bu = true
	} else if _optType == OptTypeTing {
		cmaj.OptInfo.Ting = true
	}

	cmaj.OptInfo.Cancer = true //任何时候都可以放弃
}

//位置操作类型是否可用
func (cmaj *CMaj) CanOpt(_optType int) bool {

	if _optType == OptTypePeng {
		return cmaj.OptInfo.Peng
	} else if _optType == OptTypeGang {
		return cmaj.OptInfo.Gang
	} else if _optType == OptTypeChi {
		return cmaj.OptInfo.Chi
	} else if _optType == OptTypeHu {
		return cmaj.OptInfo.Hu
	} else if _optType == OptTypeBu {
		return cmaj.OptInfo.Bu
	} else if _optType == OptTypeCancel {
		return cmaj.OptInfo.Cancer
	} else if _optType == OptTypeTing {
		return cmaj.OptInfo.Ting
	}
	return false
}

//获取杠牌,机器人使用
func (cmaj *CMaj) RobotGetGangPai() int {
	for i := 0; i < len(cmaj.OptInfo.GangCard); i++ {
		if cmaj.OptInfo.GangCard[i] != consts.DefaultIndex {
			return cmaj.OptInfo.GangCard[i]
		}
	}
	logs.Info("tableId:%v,******************************** error RobotGetGangPai", cmaj.TableCfg.TableId)
	return -1
}

// 设置手牌数组
func (cmaj *CMaj) SetHandPai(cards MCards) {

	//按牌花色赋值
	for _, v := range cards {
		c := v.GetColor()
		card := NewMCard(v.GetData())
		cmaj.HandPaiArr[c] = append(cmaj.HandPaiArr[c], card)
	}
	//初始化手牌后排序
	cmaj.SortHandPai()
	logs.Info("tableId:%v---------------->SetHandPai,座位:%v:%v", cmaj.TableCfg.TableId, cmaj.SeatID, cmaj.HandPaiArr)
}

// 排序手牌
func (cmaj *CMaj) SortHandPai() {

	for k, _ := range cmaj.HandPaiArr {
		cardSort := NewMCardSorter(cmaj.HandPaiArr[k])
		cardSort.Sort()
		cmaj.HandPaiArr[k] = cardSort
	}
}

//// 排序的手牌数组对象
//func (cmaj *CMaj) SortHandIntArr() {
//
//	for i := 0; i < len(cmaj.HandIntArr); i++ {
//		cmaj.HandIntArr[i] = utils.SortIntArrAsc(cmaj.HandIntArr[i])
//	}
//	//cmaj.HandIntArr = make([][]int, 0)
//	//for _, v := range cmaj.HandPaiArr {
//	//	arr := make([]int, 0)
//	//	for _, n := range v {
//	//		arr = append(arr, n.GetData())
//	//	}
//	//	arr = utils.SortIntArrAsc(arr)
//	//	cmaj.HandIntArr = append(cmaj.HandIntArr, arr)
//	//}
//}

//检测手牌中是否含有这张牌
func (cmaj *CMaj) IsHandHasSamePai(_data int) bool {

	for _, v := range cmaj.HandPaiArr {
		for _, n := range v {
			if n.GetData() == _data {
				return true
			}
		}
	}
	return false
}

//手牌中含有等值的牌数量
func (cmaj *CMaj) HandEqualPaiCt(_card *MCard) int {

	ct := 0
	for _, v := range cmaj.HandPaiArr {
		for _, n := range v {
			if n.Equal(_card) {
				ct++
			}
		}
	}
	return ct
}

//获得手牌中等值的牌
func (cmaj *CMaj) GetHandEqualPai(_card *MCard) *MCard {

	for _, v := range cmaj.HandPaiArr {
		for _, n := range v {
			if n.Equal(_card) {
				return n
			}
		}
	}
	logs.Info("tableId:%v,************** error! GetHandEqualPai _card:%v", cmaj.TableCfg.TableId, _card)
	return nil
}

//删除一张手牌,完全相同
func (cmaj *CMaj) RemoveHandCard(_delMCard *MCard) {
	color := _delMCard.GetColor()
	tmpArr := make([]*MCard, 0)
	len1 := len(cmaj.HandPaiArr[color])
	hasRemove := false
	for _, v := range cmaj.HandPaiArr[color] {
		if !hasRemove && v.Same(_delMCard) { //只删除一张,后期计算查大叫时会有重复的牌
			hasRemove = true
		} else {
			tmpArr = append(tmpArr, v)
		}
	}
	cmaj.HandPaiArr[color] = tmpArr
	len2 := len(cmaj.HandPaiArr[color])
	if len1-len2 != 1 {
		logs.Info("tableId:%v,RemoveHandCard****************************** error!,删除手牌失败,RemoveHandCard:%v,len1-len2:%v,cmaj.HandPaiArr[color]:%v", cmaj.TableCfg.TableId, _delMCard.String(), len1-len2, cmaj.HandPaiArr[color])
	}
}

//删除一张手牌,相等的牌
func (cmaj *CMaj) RemoveHandEqualMCard(_delMCard *MCard) *MCard {
	rtnMCard := NewMCard(-1)
	color := _delMCard.GetColor()
	tmpArr := make([]*MCard, 0)
	hasRemove := false
	len1 := len(cmaj.HandPaiArr[color])
	for _, v := range cmaj.HandPaiArr[color] {
		if !hasRemove && v.Equal(_delMCard) { //一次只删除一个
			rtnMCard = v.Clone()
			hasRemove = true
		} else {
			tmpArr = append(tmpArr, v)
		}
	}
	cmaj.HandPaiArr[color] = tmpArr
	len2 := len(cmaj.HandPaiArr[color])

	if len1-len2 != 1 {
		logs.Info("tableId:%v,RemoveHandEqualMCard*******************************  error!,删除手牌失败,RemoveHandEqualMCard _delMCard:%v,len1-len2:%v,cmaj.HandPaiArr:%v",
			cmaj.TableCfg.TableId, _delMCard.String(), len1-len2, cmaj.HandPaiArr)
	}
	return rtnMCard
}

//删除一张碰牌,相等的牌 (面杠使用)
func (cmaj *CMaj) RemovePengEqualMCard(_delMCard *MCard) *MCard {
	rtnMCard := NewMCard(-1)
	color := _delMCard.GetColor()
	tmpArr := make([]*MCard, 0)
	hasDel := false
	len1 := len(cmaj.PengPaiArr[color])
	for _, v := range cmaj.PengPaiArr[color] {
		if v.Equal(_delMCard) && !hasDel { //一次只删除一个
			rtnMCard = v.Clone()
			hasDel = true
		} else {
			tmpArr = append(tmpArr, v)
		}
	}
	cmaj.PengPaiArr[color] = tmpArr
	len2 := len(cmaj.PengPaiArr[color])

	if len1-len2 != 1 {
		logs.Info("tableId:%v,*********************************************** error!,删除碰牌失败,RemovePengEqualMCard _delMCard:%v,len1-len2:%v,cmaj.PengPaiArr[color]:%",
			cmaj.TableCfg.TableId, _delMCard, len1-len2, cmaj.PengPaiArr[color])
	}
	return rtnMCard
}

//添加一张出牌
func (cmaj *CMaj) AddOutPai(_sendMCard *MCard) {
	cmaj.OutPaiArr = append(cmaj.OutPaiArr, _sendMCard)
}

//添加一张花牌
func (cmaj *CMaj) AddHuaPaiArr(_sendMCard *MCard) {
	cmaj.HuaPaiArr = append(cmaj.HuaPaiArr, _sendMCard)
}

// 添加一张手牌 isFetch:ture 是自己拿牌,false:是添加用于检测
func (cmaj *CMaj) AddHandPai(_addMCard *MCard, isFetch bool) {

	//if isFetch {
	//	logs.Info("tableId:%v ----------------> seatID:%v ,AddHandPai,添加手牌:%v [正常拿牌]", cmaj.TableCfg.TableId, cmaj.SeatID, _addMCard.String())
	//} else {
	//	logs.Info("tableId:%v ----------------> seatID:%v ,AddHandPai,添加手牌:%v [胡牌检测]", cmaj.TableCfg.TableId, cmaj.SeatID, _addMCard.String())
	//}

	color := _addMCard.GetColor()
	cmaj.HandPaiArr[color] = append(cmaj.HandPaiArr[color], _addMCard)
}

//选缺(系统推荐)
func (cmaj *CMaj) SelectQue() []int {
	//推荐数量最少的
	countArr := make([]int, 3)
	selectArr := make([]int, 0)
	for i := 0; i < 3; i++ {
		countArr[i] = len(cmaj.HandPaiArr[i])
	}

	minCount := 13
	minColor := -1
	for i := 0; i < len(countArr); i++ {
		if countArr[i] < minCount {
			minColor = i
			minCount = countArr[i]
		}
	}
	selectArr = append(selectArr, minColor) //记录这个颜色
	//可能存在有2个一样多的
	for i := 0; i < len(countArr); i++ {
		if i != minColor && countArr[i] == minCount { //如果值相同,也计入
			selectArr = append(selectArr, i)
		}
	}

	return selectArr
}

//选取换三张的牌(系统推荐)
func (cmaj *CMaj) RecomThreeMCard() []*MCard {

	//从数量最少的一种中去取
	selectCars := make([]*MCard, 0)
	countArr := make([]int, 3)
	for i := 0; i < 3; i++ {
		//必须有三张才加入
		if len(cmaj.HandPaiArr[i]) >= 3 {
			countArr[i] = len(cmaj.HandPaiArr[i])
		}
	}

	for {
		minCount := 13
		minColor := -1
		for i := 0; i < len(countArr); i++ {
			if countArr[i] > 0 && countArr[i] <= minCount {
				minColor = i
				minCount = countArr[i]
			}
		}

		for i := 0; i < len(cmaj.HandPaiArr[minColor]); i++ {
			selectCars = append(selectCars, cmaj.HandPaiArr[minColor][i])
			if len(selectCars) == 3 {
				return selectCars
			}
		}
		//没有找满三张,说明1种或2种花色的牌总数不够3张,将countArr[i] 数量置为0即可
		countArr[minColor] = 0
	}

}

// 碰牌检测
func (cmaj *CMaj) Check_Peng(_card *MCard) bool {

	if !config.Opts().OpenPeng {
		return false
	}

	ct := 0
	//遍历手牌
	handArr := cmaj.GetHandPai()
	for _, v := range handArr {
		if v.Equal(_card) {
			ct++
			if ct == 2 { //有2张即可碰牌
				return true
			}
		}
	}

	return false
}

// 碰牌
func (cmaj *CMaj) DoPeng(_card *MCard) {

	color := _card.GetColor()
	//删除手牌中两张碰牌
	c1 := cmaj.RemoveHandEqualMCard(_card)
	c2 := cmaj.RemoveHandEqualMCard(_card)

	//添加到碰牌
	cmaj.PengPaiArr[color] = append(cmaj.PengPaiArr[color], _card)
	cmaj.PengPaiArr[color] = append(cmaj.PengPaiArr[color], c1)
	cmaj.PengPaiArr[color] = append(cmaj.PengPaiArr[color], c2)

	tmpArr := make([]int, 0)
	tmpArr = append(tmpArr, PGCTypePeng)
	tmpArr = append(tmpArr, _card.GetData())
	tmpArr = append(tmpArr, c1.GetData())
	tmpArr = append(tmpArr, c2.GetData())

	cmaj.PGCArr = append(cmaj.PGCArr, tmpArr)
}

// 面下杠检测
func (cmaj *CMaj) Check_MianGang() bool {

	if !config.Opts().OpenGang {
		return false
	}

	cmaj.SortHandPai()
	cmaj.TempGangPaiArr = make([]*MCard, 0)

	ct := 0
	//遍历碰牌,如果手中有同花色\值的牌 则可面杠
	pengArr := cmaj.GetPengPai()
	for i := 0; i < len(pengArr); i = i + 3 {
		pengPai := pengArr[i]

		if cmaj.HandEqualPaiCt(pengPai) > 0 {
			mianGangPai := cmaj.GetHandEqualPai(pengPai)
			cmaj.TempGangPaiArr = append(cmaj.TempGangPaiArr, mianGangPai)
			for k := 0; k < 3; k++ {
				cmaj.TempGangPaiArr = append(cmaj.TempGangPaiArr, pengPai)
			}
			ct++
		}
	}

	if ct > 0 {
		return true
	}
	return false
}

// 杠牌 -- 面杠
func (cmaj *CMaj) DoMianGang(_card *MCard) {

	color := _card.GetColor()

	//从手牌中删除
	c0 := cmaj.RemoveHandEqualMCard(_card)
	//删除碰牌中3张牌
	c1 := cmaj.RemovePengEqualMCard(_card)
	c2 := cmaj.RemovePengEqualMCard(_card)
	c3 := cmaj.RemovePengEqualMCard(_card)

	//添加到杠牌
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c0)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c1)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c2)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c3)

	tmpArr := make([]int, 0)
	tmpArr = append(tmpArr, PGCTypeMingGang)
	tmpArr = append(tmpArr, c0.GetData())
	tmpArr = append(tmpArr, c1.GetData())
	tmpArr = append(tmpArr, c2.GetData())
	tmpArr = append(tmpArr, c3.GetData())

	//找到这组碰,并把它改成杠类型
	for k, v := range cmaj.PGCArr {
		if v[0] == PGCTypePeng {
			tmpMCard := NewMCard(v[1])
			if tmpMCard.Equal(_card) {
				//重新赋值
				cmaj.PGCArr[k] = make([]int, 0)
				cmaj.PGCArr[k] = tmpArr
			}
		}
	}

}

// 明杠检测
func (cmaj *CMaj) Check_MingGang(_card *MCard) bool {

	if !config.Opts().OpenGang {
		return false
	}

	cmaj.SortHandPai()
	cmaj.TempGangPaiArr = make([]*MCard, 0)
	color := _card.GetColor()

	ct := 0
	for _, v := range cmaj.HandPaiArr[color] {
		if v.Equal(_card) {
			cmaj.TempGangPaiArr = append(cmaj.TempGangPaiArr, v)
			ct++
		}
	}

	if ct == 3 {
		//加入这张牌
		cmaj.TempGangPaiArr = append(cmaj.TempGangPaiArr, _card)
		return true
	}
	return false
}

// 杠牌 -- 明杠
func (cmaj *CMaj) DoMingGang(_card *MCard) {

	cmaj.MingGangCt++ //记录明杠数
	color := _card.GetColor()

	//删除手牌中3张牌
	c1 := cmaj.RemoveHandEqualMCard(_card)
	c2 := cmaj.RemoveHandEqualMCard(_card)
	c3 := cmaj.RemoveHandEqualMCard(_card)

	//添加到杠牌
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], _card)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c1)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c2)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c3)

	tmpArr := make([]int, 0)
	tmpArr = append(tmpArr, PGCTypeMingGang)
	tmpArr = append(tmpArr, _card.GetData())
	tmpArr = append(tmpArr, c1.GetData())
	tmpArr = append(tmpArr, c2.GetData())
	tmpArr = append(tmpArr, c3.GetData())

	cmaj.PGCArr = append(cmaj.PGCArr, tmpArr)
}

// 杠牌 -- 暗杠
func (cmaj *CMaj) DoAnGang(_card *MCard) {

	color := _card.GetColor()

	//删除手牌中4张牌
	c1 := cmaj.RemoveHandEqualMCard(_card)
	c2 := cmaj.RemoveHandEqualMCard(_card)
	c3 := cmaj.RemoveHandEqualMCard(_card)
	c4 := cmaj.RemoveHandEqualMCard(_card)

	//添加到杠牌
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c1)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c2)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c3)
	cmaj.GangPaiArr[color] = append(cmaj.GangPaiArr[color], c4)

	tmpArr := make([]int, 0)
	tmpArr = append(tmpArr, PGCTypeAnGang)
	tmpArr = append(tmpArr, c1.GetData())
	tmpArr = append(tmpArr, c2.GetData())
	tmpArr = append(tmpArr, c3.GetData())
	tmpArr = append(tmpArr, c4.GetData())

	cmaj.PGCArr = append(cmaj.PGCArr, tmpArr)
}

//额外牌型检测 [支]:某一门有 8 张是胡牌的的基本要求，如有 9 张则加 1 支， 10 张加 2 支，以此类推，多支多分 ( 一支+1嘴) 。
func (cmaj *CMaj) Check_ZHI() int {

	countArr := make([]int, 3) //万筒条
	allCard := cmaj.GetAllPai()
	for _, v := range allCard {
		_color := v.GetColor()
		countArr[_color]++
	}

	for _, v := range countArr {
		if v >= 8 {
			return v - 8
		}
	}
	return 0
}

//额外牌型检测 [卡] 只胡一张牌的牌型所胡的牌称为卡。+1嘴 23胡4 78胡6也算卡,将不算
//func (cmaj *CMaj) Check_KA() bool {
//
//	huCard := cmaj.HuMCard
//
//	//判断胡的是否是将牌
//	jiangMCard := cmaj.GetJiangPai()
//	if jiangMCard != nil {
//		if jiangMCard.Equal(huCard) {
//			return false
//		}
//	}
//
//	logs.Info("tableId:%v-------------Check_KA------>huCard:%v ,cmaj.PxList:%v ", cmaj.TableCfg.TableId, huCard.Detail(), cmaj.PxList)
//	//检测是否胡的 ABC 牌型
//	for _, v := range cmaj.PxList {
//		if len(v) == 0 {
//			continue
//		}
//		if cmaj.CheckABCPai(v[0], v[1], v[2]) {
//			huValue := huCard.GetValue()
//			if utils.HasElement(v, huCard.GetData()) {
//				if (GetVal(v[1]) == 2 && huValue == 4) || (huValue == 6 && GetVal(v[2]) == 8) || huCard.GetData() == v[1] {
//					return true
//				}
//			}
//		}
//	}
//
//	return false
//}

//额外牌型检测 ［缺门］   胡牌时所有牌只存在有2门花色。+2嘴
func (cmaj *CMaj) Check_QUE_MEN() bool {

	//只能缺1门,缺2门就是清一色了
	queMengCt := 0
	if len(cmaj.HandPaiArr[Color_Wan]) == 0 {
		queMengCt++
	}
	if len(cmaj.HandPaiArr[Color_Tong]) == 0 {
		queMengCt++
	}
	if len(cmaj.HandPaiArr[Color_Tiao]) == 0 {
		queMengCt++
	}

	if queMengCt == 1 {
		return true
	}

	return false
}

//额外牌型检测 ［同］：所有牌中数字一样的牌从 4 张起数， 每多一张多+1嘴（如， 3 张四万、 2 张四条、 3 张四饼的牌得 4倍每基础分）
func (cmaj *CMaj) Check_TONG() int {

	maxTongCt := 0
	for i := 2; i <= 8; i++ {
		tongCt := 0
		handPai := cmaj.GetHandPai()
		for _, v := range handPai {
			if v.GetValue() == i {
				tongCt++
			}
		}
		if tongCt > maxTongCt {
			maxTongCt = tongCt
		}
	}

	if maxTongCt > 4 {
		return maxTongCt - 4
	}

	return 0
}

//额外牌型检测 ［10同］：   同超过10张 +10嘴
func (cmaj *CMaj) Check_10TONG() bool {

	maxTongCt := 0
	for i := 2; i <= 8; i++ {
		tongCt := 0
		handPai := cmaj.GetHandPai()
		for _, v := range handPai {
			if v.GetValue() == i {
				tongCt++
			}
		}
		if tongCt > maxTongCt {
			maxTongCt = tongCt
		}
	}

	if maxTongCt >= 10 {
		return true
	}

	return false
}

//额外牌型检测 ［坎］：三张一样的牌在手，且符合基本胡牌牌型中的刻 ，即三张牌未分开叫做一坎，每一坎+1嘴
//func (cmaj *CMaj) Check_KAN() int {
//
//	//logs.Info("tableId:%v----------------Check_KAN---cmaj.PxList:%v", cmaj.TableCfg.TableId, cmaj.PxList)
//	kanCt := 0
//	//检测 AAA 牌型 数
//	for _, v := range cmaj.PxList {
//		if len(v) == 0 {
//			continue
//		}
//		for i := 0; i < len(v); i = i + 3 {
//			if cmaj.CheckAAAPai(v[i], v[i+1], v[i+2]) {
//				kanCt++
//				//logs.Info("tableId:%v----------------kanCt++:%v", cmaj.TableCfg.TableId, kanCt)
//			}
//		}
//
//	}
//	logs.Info("tableId:%v----------------Check_KAN---kanCt:%v ", cmaj.TableCfg.TableId, kanCt)
//	return kanCt
//}

//额外牌型检测 ［四暗刻］：  4坎加一个搭牌，坎牌必须是自己摸的才算，TODO 点炮不算 +10嘴
//func (cmaj *CMaj) Check_4ANKAN() bool {
//
//	kanCt := cmaj.Check_KAN()
//	if kanCt == 4 {
//		return true
//	}
//	return false
//}

//额外牌型检测 ［3连坎］：   同花色连在一起的3个坎 +10嘴
//func (cmaj *CMaj) Check_3LIAN_KAN() bool {
//
//	//取出所有坎牌
//	kanArr := make([]int, 0)
//
//	for _, v := range cmaj.PxList {
//		if len(v) == 0 {
//			continue
//		}
//		for i := 0; i < len(v); i = i + 3 {
//			if cmaj.CheckAAAPai(v[i], v[i+1], v[i+2]) {
//				kanArr = append(kanArr, v[i])
//			}
//		}
//
//	}
//
//	logs.Info("tableId:%v----------------Check_3LIAN_KAN---kanArr:%v ", cmaj.TableCfg.TableId, kanArr)
//	if len(kanArr) >= 3 {
//		kanArr = utils.SortIntArrAsc(kanArr) //升序排列
//		for i := 0; i < len(kanArr); i = i + 2 {
//			if i+2 < len(kanArr) {
//				v1 := GetVal(kanArr[i])
//				v2 := GetVal(kanArr[i+1])
//				v3 := GetVal(kanArr[i+2])
//				if v2 == v1+1 && v3 == v2+1 {
//					return true
//				}
//			}
//		}
//	}
//
//	return false
//}

//额外牌型检测 ［四活］：    某种牌有四张且没有开杠，只要有四张的就算活（豪华七对也算活）。+4嘴
func (cmaj *CMaj) Check_SIHUO() int {

	ct := 0
	for _, v := range cmaj.HandPaiArr {
		iSize := len(v)
		if iSize >= 4 {
			count := 0
			value := -1
			for _, n := range v {
				val := n.GetValue()
				if val != value {
					count = 1
					value = n.GetValue()
				} else {
					count++
				}

				if count == 4 {
					ct++
				}
			}
		}
	}

	return ct
}

//额外牌型检测 ［双铺子］：    两个一样的顺，如345万 345万，若胡牌是胡其中的一张算明双铺，+2嘴 若不是胡其中的牌则为暗双铺+4嘴 自摸+10嘴
//func (cmaj *CMaj) Check_SHUANG_PUZI() []int {
//
//	shuanPuArr := make([]int, 0) //记录双铺子牌值
//	tmpArr := make([]int, 0)
//	//找出 ABC 牌型
//	for _, v := range cmaj.PxList {
//		if len(v) == 0 {
//			continue
//		}
//		for i := 0; i < len(v); i = i + 3 {
//			if cmaj.CheckABCPai(v[i], v[i+1], v[i+2]) {
//				tmpArr = append(tmpArr, v[i])
//				tmpArr = append(tmpArr, v[i+1])
//				tmpArr = append(tmpArr, v[i+2])
//			}
//		}
//	}
//
//	for i := 0; i < len(tmpArr); i = i + 3 {
//		ci := GetColor(tmpArr[i])
//		vi := GetVal(tmpArr[i])
//		for j := i + 3; j < len(tmpArr); j = j + 3 {
//			cj := GetColor(tmpArr[j])
//			vj := GetVal(tmpArr[j])
//			if ci == cj && vi == vj {
//				shuanPuArr = append(shuanPuArr, tmpArr[i])
//				shuanPuArr = append(shuanPuArr, tmpArr[i+1])
//				shuanPuArr = append(shuanPuArr, tmpArr[i+2])
//				shuanPuArr = append(shuanPuArr, tmpArr[j])
//				shuanPuArr = append(shuanPuArr, tmpArr[j+1])
//				shuanPuArr = append(shuanPuArr, tmpArr[j+2])
//				break
//			}
//		}
//	}
//
//	return shuanPuArr
//}

//检测 带幺九（2番）：每副顺子、刻子、将牌都包含一或九
//func (cmaj *CMaj) Check_DAI19() bool {
//
//	//先判断将
//	jiangMCard := cmaj.GetJiangPai()
//	if jiangMCard != nil {
//		jiangVal := jiangMCard.GetValue()
//		if jiangVal != 1 && jiangVal != 9 {
//			return false
//		}
//	}
//
//	//判断碰
//	pengPai := cmaj.GetPengPai()
//	if len(pengPai) > 0 {
//		for i := 0; i < len(pengPai); i += 3 {
//			val1 := pengPai[i].GetValue()
//			val2 := pengPai[i+1].GetValue()
//			val3 := pengPai[i+2].GetValue()
//			if val1 != 1 && val1 != 9 && val2 != 1 && val2 != 9 && val3 != 1 && val3 != 9 {
//				return false
//			}
//		}
//	}
//
//	//判断杠
//	gangPai := cmaj.GetGangPai()
//	if len(gangPai) > 0 {
//		for i := 0; i < len(gangPai); i += 4 {
//			val1 := gangPai[i].GetValue()
//			if val1 != 1 && val1 != 9 {
//				return false
//			}
//		}
//	}
//
//	//判断胡牌牌型中的 AAA  ABC 是否有1,9
//	for _, v := range cmaj.PxList {
//		if len(v) > 0 {
//			for i := 0; i < len(v); i += 3 {
//				val1 := GetVal(v[i])
//				val2 := GetVal(v[i+1])
//				val3 := GetVal(v[i+2])
//				if val1 != 1 && val1 != 9 && val2 != 1 && val2 != 9 && val3 != 1 && val3 != 9 {
//					return false
//				}
//			}
//		}
//	}
//	return true
//}

//获取胡牌的将牌
func (cmaj *CMaj) GetJiangPai(_pxList [][]int) *MCard {

	handPai := cmaj.GetHandPai()
	//logs.Info("GetJiangPai--------------->handTing111:%v", handTing)
	if len(handPai) == 2 {
		return handPai[0]
	}
	//删去胡牌时AAA ABC 剩下的即是将
	for _, v := range _pxList {
		for _, n := range v {
			handPai = cmaj.DeleteArrayCard(handPai, n)
		}
		logs.Info("tableId:%v,~~~~~~~~~~~~~~~ GetJiangPai handPai:%v", cmaj.TableCfg.TableId, handPai)

	}
	logs.Info("tableId:%v, GetJiangPai handPai:%v", cmaj.TableCfg.TableId, handPai)

	if len(handPai) != 2 {
		logs.Info("tableId:%v,*******************************error!!  GetJiangPai handPai:%v _pxList:%v",
			cmaj.TableCfg.TableId, handPai, _pxList)
		return nil
	}
	if len(handPai) == 0 {
		return nil
	}
	return handPai[0]
}

//删除牌数组中的 值为_data的所有牌,返回新数组
func (cmaj *CMaj) DeleteArrayCard(arr []*MCard, _data int) []*MCard {

	tmp := make([]*MCard, 0)
	hasDel := false
	for _, v := range arr {
		if !hasDel && v.GetData() == _data { //一次删除一个

			hasDel = true
		} else {
			tmp = append(tmp, v)
		}
	}

	return tmp
}

//全带幺检测
func (cmaj *CMaj) Check_QDY(calHuInfo *CalHuInfo) bool {
	jiang := cmaj.GetJiangPai(calHuInfo.PxList)
	if jiang == nil {
		return false
	}

	if jiang.GetColor() > Color_Tiao {
		return false
	}
	if jiang.GetValue() != 1 && jiang.GetValue() != 9 {
		return false
	}

	for _, v := range calHuInfo.PxList {
		if len(v) == 3 {
			if NewMCard(v[0]).GetColor() > Color_Tiao {
				return false
			}
			if NewMCard(v[0]).GetValue() == 1 || NewMCard(v[2]).GetValue() == 9 {
				continue
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

//三杠
func (cmaj *CMaj) Check_SG() bool {
	if cmaj.GetGangCt()/4 == 3 {
		return true
	}
	return false
}

//清龙检测
func (cmaj *CMaj) Check_QL() bool {
	MostCards := cmaj.HandPaiArr[cmaj.MostColor()]
	for i := 1; i < 10; i++ {
		isHas := false
		for _, v := range MostCards {
			if v.GetValue() == i {
				isHas = true
			}
		}

		if isHas {

		} else {
			return false
		}
	}
	return true
}

//混幺九检测
func (cmaj *CMaj) Check_HYJ(calHuInfo *CalHuInfo) bool {

	//由字牌和序数牌一、九的刻子、将牌组成的糊牌。不计对对糊。
	if !cmaj.Check_DDHu(calHuInfo) {
		return false
	}

	allCards := cmaj.GetAllPai()
	for _, v := range allCards {
		if v.GetColor() < Color_Feng {
			if v.GetValue() != 1 && v.GetValue() != 9 {
				return false
			}
		}
	}
	return true
}

//检测 清一色（2番）：全是一种花色的平胡
func (cmaj *CMaj) Check_QYS() bool {
	//牌全是一种花色，全字牌不计清一色

	allMCardArr := cmaj.GetAllPai()
	color := allMCardArr[0].GetColor()
	if color > Color_Tiao {
		return false
	}

	for _, v := range allMCardArr {
		if v.GetColor() != color {
			return false
		}
	}
	return true
}

//四暗刻检测
func (cmaj *CMaj) Check_SAK() bool {
	//四个暗刻（杠），必须在手上，
	handCards := cmaj.GetHandPai()
	anKeCt := 0
	for _, v := range handCards {
		cardCount := 0
		for _, k := range handCards {
			if v.Equal(k) {
				cardCount++
			}
		}
		if cardCount >= 3 {
			anKeCt++
		}
	}
	if anKeCt == 4 {
		return true
	}
	return false
}

//字一色检测检测
func (cmaj *CMaj) Check_ZYS() bool {

	//牌中全是字牌
	allcards := cmaj.GetAllPai()

	for _, v := range allcards {
		if v.GetColor() < 3 {
			return false
		}
	}
	return true
}

//小四喜检测
func (cmaj *CMaj) Check_XSX(_calHuInfo *CalHuInfo) bool {

	// 有三幅风刻子，（包括杠、碰）并且风牌做将，
	dongCt := 0
	nanCt := 0
	xiCt := 0
	beiCt := 0
	card_Dong := NewMCard(108)
	card_Nan := NewMCard(112)
	card_Xi := NewMCard(116)
	card_Bei := NewMCard(120)

	jiang := cmaj.GetJiangPai(_calHuInfo.PxList)

	allCards := cmaj.GetAllPai()

	for _, v := range allCards {
		if v.Equal(card_Dong) {
			dongCt++
		}
		if v.Equal(card_Nan) {
			nanCt++
		}
		if v.Equal(card_Xi) {
			xiCt++
		}
		if v.Equal(card_Bei) {
			beiCt++
		}
	}
	fengKeCt := 0
	if dongCt >= 3 {
		fengKeCt++
	}
	if nanCt >= 3 {
		fengKeCt++
	}
	if xiCt >= 3 {
		fengKeCt++
	}
	if beiCt >= 3 {
		fengKeCt++
	}
	if fengKeCt == 3 && jiang != nil && jiang.GetColor() == Color_Feng {
		return true
	}
	return false
}

//混一色检测
func (cmaj *CMaj) Check_HYS() bool {
	xuCt := 0
	ziCt := 0
	WanCt := 0
	TongCt := 0
	TiaoCt := 0
	ZfbCt := 0
	FengCt := 0
	AllPai := cmaj.GetAllPai()
	for _, v := range AllPai {
		if v.GetColor() == Color_Wan {
			WanCt++
		}
		if v.GetColor() == Color_Tong {
			TongCt++
		}
		if v.GetColor() == Color_Tiao {
			TiaoCt++
		}
		if v.GetColor() == Color_Zfb {
			ZfbCt++
		}
		if v.GetColor() == Color_Feng {
			FengCt++
		}
	}

	if WanCt > 0 {
		xuCt++
	}
	if TongCt > 0 {
		xuCt++
	}
	if TiaoCt > 0 {
		xuCt++
	}
	if ZfbCt > 0 {
		ziCt++
	}
	if FengCt > 0 {
		ziCt++
	}
	if xuCt == 1 && ziCt != 0 {
		return true
	}
	return false
}

//大三元检测
func (cmaj *CMaj) Check_DSY() bool {

	// 中发白都为3张 ，包括碰牌
	zhongCt := 0
	faCt := 0
	baiCt := 0
	card_zhong := NewMCard(124)
	card_fa := NewMCard(128)
	card_bai := NewMCard(132)

	allCards := cmaj.GetAllPai()

	for _, v := range allCards {
		if v.Equal(card_zhong) {
			zhongCt++
		}
		if v.Equal(card_fa) {
			faCt++
		}
		if v.Equal(card_bai) {
			baiCt++
		}

	}

	if zhongCt == 3 && faCt == 3 && baiCt == 3 {
		return true
	}
	return false
}

//返回手牌中数量最多的花色
func (cmaj *CMaj) MostColor() int {
	wan := len(cmaj.HandPaiArr[0])
	tong := len(cmaj.HandPaiArr[1])
	tiao := len(cmaj.HandPaiArr[2])
	if wan > tong && wan > tiao {
		return 0
	}
	if tong > wan && tong > tiao {
		return 1
	}
	if tiao > wan && tiao > tong {
		return 2
	}
	return 0
}

//清幺九检测

func (cmaj *CMaj) Check_QYJ() bool {
	//全是1和9 不计对对胡,碰牌算
	allCards := cmaj.GetAllPai()

	for _, v := range allCards {
		if v.GetValue() != 9 && v.GetValue() != 1 {
			return false
		}
	}
	return true
}

//四杠
func (cmaj *CMaj) Check_SIGANG() bool {

	//玩家杠了四次
	if cmaj.GetGangCt()/4 > 3 {
		return true
	}
	return false
}

//九连宝灯
func (cmaj *CMaj) Check_JLBD() bool {
	//由一种花色序数牌按１１１２３４５６７８９９９组成的特定排型，见同花色任何一张序数牌极为糊牌。不计清一色。
	//可以听该门花色的任何一张牌
	MostCards := cmaj.HandPaiArr[cmaj.MostColor()]
	valueOneCt := 0
	valueNineCt := 0

	for i := 1; i < 10; i++ {
		isHas := false
		for _, v := range MostCards {
			if v.GetValue() == i {
				isHas = true
			}
		}

		if isHas {

		} else {
			return false
		}
	}
	for _, v := range MostCards {
		if v.GetValue() == 1 {
			valueOneCt++
		}
		if v.GetValue() == 10 {
			valueNineCt++
		}
	}
	if valueOneCt >= 3 && valueNineCt >= 3 {
		return true
	}
	return false
}

//绿一色检测
func (cmaj *CMaj) Check_LYS() bool {

	//手上的牌全是绿的 不计混一色
	allCards := cmaj.GetAllPai()
	card_fa := NewMCard(128)
	for _, v := range allCards {
		if v.GetColor() != Color_Tiao && !v.Equal(card_fa) {
			return false
		}
	}
	return true
}

//大四喜检测
func (cmaj *CMaj) Check_DSX() bool {
	// 东南西北 都是大于等于3张的，包括碰杠
	dongCt := 0
	nanCt := 0
	xiCt := 0
	beiCt := 0
	card_Dong := NewMCard(108)
	card_Nan := NewMCard(112)
	card_Xi := NewMCard(116)
	card_Bei := NewMCard(120)

	allCards := cmaj.GetAllPai()

	for _, v := range allCards {
		if v.Equal(card_Dong) {
			dongCt++
		}
		if v.Equal(card_Nan) {
			nanCt++
		}
		if v.Equal(card_Xi) {
			xiCt++
		}
		if v.Equal(card_Bei) {
			beiCt++
		}
	}

	if dongCt >= 3 && nanCt >= 3 && xiCt >= 3 && beiCt >= 3 {
		return true
	}
	return false
}

//检测 门清：没有碰过牌,没有杠
func (cmaj *CMaj) Check_MQ() bool {

	if cmaj.GetPengCt() > 0 {
		//logs.Info("Check_MQ------------------------------cmaj.GetPengCt():%v---------------->>>>", cmaj.GetPengCt())
		return false
	}

	if cmaj.ZhiGCt > 0 || cmaj.MianGCt > 0 {
		//logs.Info("Check_MQ------------------------------cmaj.GetPengCt():%v---------------->>>>", cmaj.GetPengCt())
		return false
	}

	return true
}

//检测 卡2条（1番）：胡牌时，牌型为一条和三条，胡二条。
//func (cmaj *CMaj) Check_Ka2Tiao() bool {
//
//	if cmaj.TableCfg != nil {
//		//if cmaj.TableCfg.TableType != consts.TableTypeFour { //只有4人2房有
//		//	return false
//		//}
//		//if cmaj.TableCfg.Ka2Tiao == consts.No { //开启了选项
//		//	return false
//		//}
//	}
//
//	//找出胡牌牌型中含有胡牌
//	huMCard := cmaj.HuMCard
//	if huMCard.GetColor() == Color_Tiao && huMCard.GetValue() == 2 {
//		tiaoArr := cmaj.PxList[Color_Tiao]
//		if len(tiaoArr) > 0 {
//			for i := 0; i < len(tiaoArr); i += 3 {
//				if tiaoArr[i+1] == cmaj.HuMCard.GetData() && GetVal(tiaoArr[i]) == 1 && GetVal(tiaoArr[i+2]) == 3 {
//					return true
//				}
//			}
//		}
//	}
//
//	return false
//}

//检测 夹心5（1番）：胡牌时，牌值为46，胡5。
//func (cmaj *CMaj) Check_JiaXin5() bool {
//
//	//找出胡牌牌型中含有胡牌
//	huMCard := cmaj.HuMCard
//	if huMCard.GetValue() == 5 {
//		if len(cmaj.PxList[huMCard.GetColor()]) > 0 {
//			for i := 0; i < len(cmaj.PxList[huMCard.GetColor()]); i += 3 {
//
//				data1 := cmaj.PxList[huMCard.GetColor()][i+1]
//				data0 := cmaj.PxList[huMCard.GetColor()][i]
//				data2 := cmaj.PxList[huMCard.GetColor()][i+2]
//
//				//这个用 value 判断,只要有等值的牌 就说明可以胡这张牌
//				if GetVal(data1) == huMCard.GetValue() && GetVal(data0) == 4 && GetVal(data2) == 6 {
//					return true
//				}
//			}
//		}
//	}
//
//	return false
//}

//检测 中张（1番）：胡牌时所有牌没有1和9
func (cmaj *CMaj) Check_ZZ() bool {

	if cmaj.TableCfg != nil {
		//if cmaj.TableCfg.Mqzz == consts.No {
		//	return false
		//}
	}

	allPai := cmaj.GetAllPai()
	for _, v := range allPai {
		if v.GetValue() == 1 || v.GetValue() == 9 {
			return false
		}
	}

	return true
}

// 补花检测
func (cmaj *CMaj) Check_BB_BuHua() bool {
	//东风、红中、发财、白板作为花牌，放置到手牌最右侧
	cmaj.TempBuHuaArr = make([]*MCard, 0)
	//cmaj.SortHandPai()
	//cmaj.TempBuHuaArr = make([]*MCard, 0)
	//
	ct := 0
	//遍历手牌,如果手中有同花色\值的牌 则可面杠
	handpai := cmaj.GetHandPai()
	//for i := 0; i < len(handpai); i ++
	for _, v := range handpai {
		if v.IsBBHuaPai() {
			ct++
			cmaj.BuHua.HandHuaCards = append(cmaj.BuHua.HandHuaCards, v.GetData())
			cmaj.TempBuHuaArr = append(cmaj.TempBuHuaArr, v)
		}
	}

	if ct > 0 {
		return true
	}
	return false
}

// 补花检测
func (cmaj *CMaj) Check_HY_BuHua(fengLing int) bool {
	//东风、红中、发财、白板作为花牌，放置到手牌最右侧
	cmaj.TempBuHuaArr = make([]*MCard, 0)
	cmaj.SortHandPai()
	//cmaj.TempBuHuaArr = make([]*MCard, 0)
	//
	ct := 0
	//遍历手牌,如果手中有同花色\值的牌 则可面杠
	cmaj.BuHua.HandHuaCards = make([]int, 0)
	handpai := cmaj.GetHandPai()
	//for i := 0; i < len(handpai); i ++
	for _, v := range handpai {
		if v.IsHYHuaPai(fengLing) {
			ct++
			cmaj.BuHua.HandHuaCards = append(cmaj.BuHua.HandHuaCards, v.GetData())
			cmaj.TempBuHuaArr = append(cmaj.TempBuHuaArr, v)
		}
	}

	if ct > 0 {
		return true
	}
	return false
}

//检测 对对胡（1番）：四副刻(或杠)加一对将
func (cmaj *CMaj) Check_DDHu(_calHuInfo *CalHuInfo) bool {
	_pxList := _calHuInfo.PxList
	shun := cmaj.GetShunZiArr(_pxList)
	for _, v := range shun {
		if len(v) > 0 {
			return false
		}
	}
	return true
	//ddCount := 0 //刻数量
	//for _, v := range _pxList {
	//	for i := 0; i < len(v); i += 3 {
	//		if GetVal(v[i]) == GetVal(v[i+1]) && GetVal(v[i+1]) == GetVal(v[i+2]) {
	//			ddCount++
	//		}
	//	}
	//
	//}
	//logs.Info("tableId:%v----------Check_DDHu----------->_pxList:%v,ddCount:%v", cmaj.TableCfg.TableId, _pxList, ddCount)
	////加上碰,杠数
	//ddCount += cmaj.GetPengCt() / 3
	//ddCount += cmaj.GetGangCt() / 4
	//
	//if ddCount == 4 {
	//	return true
	//}
	//return false
}

//怀远牌型检测-------------------------------------------
//单吊手牌为一张
func (cmaj *CMaj) Check_DanDiao() bool {
	handpai := cmaj.GetHandPai()
	if len(handpai) > 2 {
		return false
	}
	return true
}

//一般高
func (cmaj *CMaj) Check_YiBanGao(_calHuInfo *CalHuInfo) int {
	ybgct := 0
	shunArr := cmaj.GetShunZiArr(_calHuInfo.PxList)
	logs.Info("shunArr---%v", shunArr)
	for _, v := range shunArr {
		if len(v) < 6 {
			continue
		} else {
			for i := 0; i < (len(v) - 5); {
				if _calHuInfo.CheckAABBCCPai_hy(v[i], v[i+1], v[i+2], v[i+3], v[i+4], v[i+5]) &&
					GetVal(v[i]) == GetVal(v[i+3]) && GetVal(v[i+1]) == GetVal(v[i+4]) && GetVal(v[i+2]) == GetVal(v[i+5]) {
					logs.Info("------%v,%v,%v,%v,%v,%v,", v[i], v[i+1], v[i+2], v[i+3], v[i+4], v[i+5])
					ybgct++
					return ybgct

				} else {
					i += 3

				}
			}
		}

	}
	keZiArr := cmaj.GetKeZiArr(_calHuInfo.PxList)
	for _, v := range keZiArr {
		if len(v) < 9 {
			continue
		} else {
			for i := 0; i < (len(v) - 6); {
				if _calHuInfo.CheckABCPai_hy(v[i], v[i+3], v[i+6]) {
					ybgct++
					return ybgct
				} else {
					i += 3

				}
			}
		}

	}
	return ybgct
}
func (cmaj *CMaj) Check_YiBanGao_7DUI(_calHuInfo *CalHuInfo) int {
	ybgct := 0
	shunArr := make([][]int, 0)
	shunArr = append(shunArr, DataToVal(_calHuInfo.WanIntArr))
	shunArr = append(shunArr, DataToVal(_calHuInfo.TongIntArr))
	shunArr = append(shunArr, DataToVal(_calHuInfo.TiaoIntArr))
	//shunArr[0] = DataToVal(_calHuInfo.WanIntArr)
	//shunArr[1] = DataToVal(_calHuInfo.TongIntArr)
	//shunArr[2] = DataToVal(_calHuInfo.TiaoIntArr)
	for _, v := range shunArr {
		if len(v) < 6 {
			continue
		} else {
			// 四张一样的去除两个
			_tempval := make([]int, 0)
			for i := 0; i < len(v); {
				if HasElementCt(v, v[i]) == 4 {
					_tempval = append(_tempval, v[i])
					_tempval = append(_tempval, v[i])
					i += 4
				} else {
					_tempval = append(_tempval, v[i])
					i++
				}
			}
			if len(_tempval) < 6 {
				continue
			}
			_ybgArr := make([]int, 0)
			for i := 0; i < len(_tempval); {
				if _tempval[i] == (_tempval[i+2]-1) && _tempval[i+2] == (_tempval[i+4]-1) {
					for j := i; j < (i + 6); j++ {
						_ybgArr = append(_ybgArr, _tempval[j])
					}
					ybgct++
					//i += 6
					break
				} else {
					i++
					if (len(_tempval) - i) < 6 {
						break
					}
				}
			}
			if len(_ybgArr) == 6 && len(v) >= 12 {
				// 去重
				_partArr := make([]int, 0)
				_partArr = DeletePartArr(v, _ybgArr)
				for i := 0; i < len(_partArr); {
					if _partArr[i] == (_partArr[i+2]-1) && _partArr[i+2] == (_partArr[i+4]-1) {
						ybgct++
						break
					} else {
						i++
						if (len(_partArr) - i) < 6 {
							break
						}
					}
				}
			}
		}

	}
	return ybgct
}

//四归一(杠不算四归一)
func (cmaj *CMaj) Check_SiGuiYi() int {
	PengPai := cmaj.GetPengPai()
	AllPai := cmaj.GetHandPai()
	AllPai = append(AllPai, PengPai...)
	sgy := 0
	for i := 0; i < len(AllPai); i++ {
		count := 0
		for k := 0; k < len(AllPai); k++ {
			if AllPai[i].Equal(AllPai[k]) {
				count++
			}
		}
		if count == 4 {
			sgy++
		}
	}
	return sgy
}

//一条龙
func (cmaj *CMaj) Check_YiTiaoLong() bool {
	hand := cmaj.GetHandPai()
	wan := make([]int, 0)
	tong := make([]int, 0)
	tiao := make([]int, 0)
	for _, v := range hand {
		if v.GetColor() == Color_Wan {
			wan = append(wan, v.GetValue())
		} else if v.GetColor() == Color_Tong {
			tong = append(tong, v.GetValue())
		} else if v.GetColor() == Color_Tiao {
			tiao = append(tiao, v.GetValue())
		}
	}
	_wan, _tong, _tiao := true, true, true
	for i := 1; i < 10; i++ {
		if HasElementCt(wan, i) <= 0 {
			_wan = false
		}
		if HasElementCt(tong, i) <= 0 {
			_tong = false
		}
		if HasElementCt(tiao, i) <= 0 {
			_tiao = false
		}
	}
	return _wan || _tiao || _tong
}

//合肥麻将8支检测
func (cmaj *CMaj) Check_8Zhi() bool {

	countArr := make([]int, 3) //万筒条
	allCard := cmaj.GetAllPai()
	for _, v := range allCard {
		_color := v.GetColor()
		countArr[_color]++
		if countArr[_color] >= 8 {
			return true
		}
	}
	return false
}

//判断2个牌组是否完全相同（下标相同）
func (cmaj *CMaj) SameArray(t1 []int, t2 []int) bool {
	t1 = SortIntArrAsc(t1)
	t2 = SortIntArrAsc(t2)
	if len(t1) != len(t2) {
		return false
	}
	for i := 0; i < len(t1); i++ {
		if t1[i] != t2[i] {
			return false
		}
	}
	return true
}

//添加胡牌牌型
func (cmaj *CMaj) AddPxId(_pxId int) {
	cmaj.HuPxIdArr = append(cmaj.HuPxIdArr, _pxId)
}

func (cmaj *CMaj) ClearPxId() {
	//logs.Info("tableId:%v---------------------------------seatId:%v------->ClearPxId()", cmaj.TableCfg.TableId, cmaj.SeatID)
	cmaj.HuPxIdArr = make([]int, 0)
}

func (cmaj *CMaj) ClearExtPxId() {
	cmaj.ExtPxIdArr = make([]int, 0)
}

func (cmaj *CMaj) HasPxId(_pxId int) bool {
	for _, v := range cmaj.HuPxIdArr {
		if v == _pxId {
			return true
		}
	}
	return false
}

//添加额外牌型id
func (cmaj *CMaj) AddExtPxId(_extPxId int) {
	logs.Info("tableId:%v,------->AddExtPxId _extPxId:%v,cmaj.ExtPxIdArr:%v", cmaj.TableCfg.TableId, _extPxId, cmaj.ExtPxIdArr)
	cmaj.ExtPxIdArr = append(cmaj.ExtPxIdArr, _extPxId)
}

//判断是否包含额外牌型id
func (cmaj *CMaj) HasExtPxId(_extPxId int) bool {
	for _, v := range cmaj.ExtPxIdArr {
		if v == _extPxId {
			return true
		}
	}
	return false
}

//返回 分数描述数组
func (cmaj *CMaj) GetPxScoreInfo() []string {
	arr := make([]string, 0)
	for _, v := range cmaj.PxScoreArr {
		str := v[0] + ":" + v[1]
		if v[1] == "" {
			str = v[0]
		}
		arr = append(arr, str)
	}
	return arr
}

//添加牌型(分数获得名称),对应的分数
func (cmaj *CMaj) AddPxScore(_pxStr string, _score int) {

	if _score == 0 { //没有分数,只添加名称
		str := []string{_pxStr, ""}
		cmaj.PxScoreArr = append(cmaj.PxScoreArr, str)
		return
	}

	//检查是否已经包含  _pxStr
	for i := 0; i < len(cmaj.PxScoreArr); i++ {
		pxsc := cmaj.PxScoreArr[i]
		name := pxsc[0]
		scstr := pxsc[1]
		if name == _pxStr { //已经包含,更新分数
			oldsc, _ := strconv.Atoi(scstr)
			updatesc := oldsc + _score
			cmaj.PxScoreArr[i][1] = strconv.Itoa(updatesc)
			return
		}
	}
	//添加新的
	str := []string{_pxStr, strconv.Itoa(_score)}
	cmaj.PxScoreArr = append(cmaj.PxScoreArr, str)

}

//获取玩家所有手牌
func (cmaj *CMaj) GetHandPai() []*MCard {
	cards := make([]*MCard, 0)

	//手牌
	for _, v := range cmaj.HandPaiArr {
		for _, n := range v {
			cards = append(cards, n)
		}
	}
	return cards
}

////获取玩家所有手牌数组
//func (cmaj *CMaj) GetHandIntArr() []int {
//
//	tmpArr := make([]int, 0)
//
//	//手牌
//	for _, v := range cmaj.HandPaiArr {
//		for _, n := range v {
//			tmpArr = append(tmpArr, n.GetData())
//		}
//	}
//	return tmpArr
//}

////按花色返回手牌 刻子 数组
func (cmaj *CMaj) GetKeZiArr(_pxList [][]int) [][]int {
	keArr := make([][]int, cmaj.TableCfg.MaxCardColorIndex)

	//手牌
	for _color, v := range _pxList {
		if len(v) > 0 {
			for i := 0; i < len(v); i += 3 {
				if CheckAAAPai(v[i], v[i+1], v[i+2]) {
					keArr[_color] = append(keArr[_color], v[i])
					keArr[_color] = append(keArr[_color], v[i+1])
					keArr[_color] = append(keArr[_color], v[i+2])
				}
			}
		}
	}

	//排序,便于计算
	for i := 0; i < len(keArr); i++ {
		keArr[i] = SortIntArrAsc(keArr[i])
	}

	return keArr
}

////按花色返回手牌 刻子 数组
func (cmaj *CMaj) GetShunZiArr(_pxList [][]int) [][]int {
	shun := make([][]int, cmaj.TableCfg.MaxCardColorIndex)

	//手牌
	for _color, v := range _pxList {
		if len(v) > 0 {
			for i := 0; i < len(v); i += 3 {
				if CheckABCPai(v[i], v[i+1], v[i+2]) {
					shun[_color] = append(shun[_color], v[i])
					shun[_color] = append(shun[_color], v[i+1])
					shun[_color] = append(shun[_color], v[i+2])
				}
			}
		}
	}

	//排序,便于计算
	//for i := 0; i < len(shun); i++ {
	//	shun[i] = SortIntArrAsc(shun[i])
	//}

	return shun
}

func (cmaj *CMaj) GetHandColorArr() [][]int {
	arr := make([][]int, cmaj.TableCfg.MaxCardColorIndex)
	for i := 0; i < len(cmaj.HandPaiArr); i++ {
		for _, n := range cmaj.HandPaiArr[i] {
			_color := n.GetColor()
			arr[_color] = append(arr[_color], n.GetData())
		}
		arr[i] = SortIntArrAsc(arr[i])
	}
	return arr
}

//如果手上有选缺的牌,则返回
func (cmaj *CMaj) GetQueMCard(_que int) *MCard {
	length := len(cmaj.HandPaiArr[_que])
	if length == 0 {
		return nil
	}
	//返回最后一张
	return cmaj.HandPaiArr[_que][length-1]
}

//获取玩家所有碰牌
func (cmaj *CMaj) GetPengPai() []*MCard {
	cards := make([]*MCard, 0)

	//碰牌
	for _, v := range cmaj.PengPaiArr {
		for _, n := range v {
			cards = append(cards, n)
		}
	}
	return cards
}

//获取玩家所有杠牌
func (cmaj *CMaj) GetGangPai() []*MCard {
	cards := make([]*MCard, 0)

	//杠牌
	for _, v := range cmaj.GangPaiArr {
		for _, n := range v {
			cards = append(cards, n)
		}
	}
	return cards
}

//获取玩家所有牌
func (cmaj *CMaj) GetAllPai() []*MCard {
	cards := make([]*MCard, 0)

	//手牌
	for _, v := range cmaj.GetHandPai() {
		cards = append(cards, v)
	}

	//碰牌
	for _, v := range cmaj.GetPengPai() {
		cards = append(cards, v)
	}

	//杠牌
	for _, v := range cmaj.GetGangPai() {
		cards = append(cards, v)
	}
	return cards
}

//获取胡牌的将牌
//func (cmaj *CMaj) GetJiangPai() *MCard {
//	handPai := cmaj.GetHandPai()
//	if len(handPai) == 2 {
//		return handPai[0]
//	}
//
//	logs.Info("tableId:%v --------------------->cmaj.PxList:%v", cmaj.TableCfg.TableId, cmaj.PxList)
//
//	//删去胡牌时AAA ABC 剩下的即是将
//	for _, v := range cmaj.PxList {
//		for _, n := range v {
//			handPai = cmaj.DeleteArrayMCard(handPai, n)
//		}
//	}
//	if len(handPai) != 2 {
//		logs.Info("tableId:%v,*******************************error!!  GetJiangPai len(handPai):%v", cmaj.TableCfg.TableId, len(handPai))
//		return nil
//	}
//	return handPai[0]
//}

type BuHua struct {
	HandHuaCards []int
	BuCards      []int
}

//新建补花数据信息对象
func NewBuHua() *BuHua {
	ret := BuHua{
		HandHuaCards: make([]int, 0),
		BuCards:      make([]int, 0),
	}
	return &ret
}

//删除牌数组中的 值为_data的所有牌,返回新数组 (在查大叫的时候会有重复的牌,一次只删去一个)
func (cmaj *CMaj) DeleteArrayMCard(arr []*MCard, _data int) []*MCard {

	tmp := make([]*MCard, 0)
	hasDel := false
	for _, v := range arr {
		if !hasDel && v.GetData() == _data { //一次删除一个
			hasDel = true
		} else {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

//获取碰牌数量
func (cmaj *CMaj) GetPengCt() int {
	count := 0
	for _, v := range cmaj.PengPaiArr {
		count += len(v)
	}
	return count
}

//获取杠牌数量
func (cmaj *CMaj) GetGangCt() int {

	count := 0
	for _, v := range cmaj.GangPaiArr {
		count += len(v)
	}
	return count
}

//位置操作数据对象
type OptInfo struct {
	Peng   bool //是否可碰
	Gang   bool //是否可杠
	Hu     bool //是否可胡
	Chi    bool //是否可吃
	Bu     bool //是否可补张
	Cancer bool //是否可取消
	Ting   bool //是否可听

	PengCard []int //可碰牌值 最大2个值
	GangCard []int //可杠牌值 最大3个值
	HuCard   int   //可胡牌值 最大1个值
	ChiCard  []int //可吃牌值 3个为一组 最大9个值
	BuCard   []int //可补张牌值 最大2个值
}

func (optInfo *OptInfo) String() string {
	str := "[碰:" + BoolToStr(optInfo.Peng) + "]" +
		"[杠:" + BoolToStr(optInfo.Gang) + "]" +
		"[胡:" + BoolToStr(optInfo.Hu) + "]" +
		"[吃:" + BoolToStr(optInfo.Chi) + "]" +
		"[补:" + BoolToStr(optInfo.Bu) + "]" +
		"[取消:" + BoolToStr(optInfo.Cancer) + "]" +
		"[报听:" + BoolToStr(optInfo.Ting) + "]"

	return str
}

//添加碰牌
func (optInfo *OptInfo) AddPengCard(_data int) {
	optInfo.PengCard = append(optInfo.PengCard, _data)
}

//添加杠牌
func (optInfo *OptInfo) AddGangCard(_data int) {
	optInfo.GangCard = append(optInfo.GangCard, _data)
}

func BoolToStr(_value bool) string {
	if _value {
		return "true"
	}
	return "false"
}

func NewOptInfo() *OptInfo {
	ret := OptInfo{

		Peng:   false, //是否可碰
		Gang:   false, //是否可杠
		Hu:     false, //是否可胡
		Chi:    false, //是否可吃
		Bu:     false, //是否可补张
		Cancer: false, //是否可取消
		Ting:   false, //是否可听牌

		PengCard: make([]int, 0),      //可碰牌值 最大2个值
		GangCard: make([]int, 0),      //可杠牌值 最大3个值
		HuCard:   consts.DefaultIndex, //可胡牌值 最大1个值
		ChiCard:  make([]int, 0),      //可吃牌值 3个为一组 最大9个值
		BuCard:   make([]int, 0),      //可补张牌值 最大3个值
	}
	return &ret
}

//按牌值大小排序牌数组--------------------------------

type MCardSorter []*MCard

func NewMCardSorter(cardArr []*MCard) MCardSorter {
	cs := make([]*MCard, 0)
	for _, v := range cardArr {
		cs = append(cs, v)
	}
	return cs
}

func (cs MCardSorter) Sort() {
	sort.Sort(cs)
}
func (cs MCardSorter) Len() int {
	return len(cs)
}

func (cs MCardSorter) Less(i, j int) bool {
	return cs[j].GetData() > cs[i].GetData() // 按牌值升序
}

func (cs MCardSorter) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}
