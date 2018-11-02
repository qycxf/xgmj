// 麻将算法工具类
package mjcomn

import (
	"sync"

	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/qo"
)

//胡牌运算数据对象 ---------------------------------------------------
type CalHuObject struct {
	MjCfg      MjCfg
	HandPaiArr [][]*MCard //手牌

}

func NewCalHuObject(_mjcfg MjCfg) *CalHuObject {
	ret := CalHuObject{
		MjCfg:      _mjcfg,
		HandPaiArr: make([][]*MCard, _mjcfg.MaxColorCt),
	}
	return &ret
}

//创建胡牌对象
func GetCalHuInfo(_mjcfg MjCfg, _handPaiArr [][]*MCard) *CalHuInfo {

	_calHuInfo := &CalHuInfo{
		Check7Dui:  _mjcfg.Check7Dui,
		WanIntArr:  make([]int, 0),
		TongIntArr: make([]int, 0),
		TiaoIntArr: make([]int, 0),
		FengIntArr: make([]int, 0),
		ZfbIntArr:  make([]int, 0),
		PxType:     PXTYPE_UNKNOW,
		PxList:     make([][]int, _mjcfg.MaxColorCt), //记录胡牌后 手牌中的 AAA和ABC 的牌型 按花色存储
	}

	//手牌
	_laiziData := _mjcfg.LaiziData
	for _, v := range _handPaiArr {
		for _, n := range v {

			//IMT! IMT!
			// 1)红中麻将 红中癞子不能当红中本身用,否则导致胡牌检测错误
			// 2)鞍山麻将 癞子牌可当本身牌值使用
			//if _mjcfg.TableType == TableType_HZMJ {
			if _laiziData != DefaultIndex && n.EqualByData(_laiziData) {
				continue
			}
			//}

			_color := n.GetColor()
			if _color == Color_Wan {
				_calHuInfo.WanIntArr = append(_calHuInfo.WanIntArr, n.GetData())
			} else if _color == Color_Tong {
				_calHuInfo.TongIntArr = append(_calHuInfo.TongIntArr, n.GetData())
			} else if _color == Color_Tiao {
				_calHuInfo.TiaoIntArr = append(_calHuInfo.TiaoIntArr, n.GetData())
			} else if _color == Color_Feng {
				_calHuInfo.FengIntArr = append(_calHuInfo.FengIntArr, n.GetData())
			} else if _color == Color_Zfb {
				_calHuInfo.ZfbIntArr = append(_calHuInfo.ZfbIntArr, n.GetData())
			}
		}
	}

	return _calHuInfo
}

//返回手上红中癞子牌个数
func GetLaiZiCt(_handPaiArr [][]*MCard, _laiziData int) int {

	if _laiziData == DefaultIndex {
		return 0
	}

	count := 0
	for _, v := range _handPaiArr {
		for _, n := range v {
			if n.EqualByData(_laiziData) {
				count++
			}
		}
	}

	return count
}

//胡牌检测 -----------------------------------
func CalHuPai(_mjcfg MjCfg, _handPaiArr [][]*MCard) *CalHuInfo {

	_tableType := _mjcfg.TableType
	_laiziData := _mjcfg.LaiziData

	if _laiziData != DefaultIndex {
		//手中是否有癞子牌
		laiZiCt := GetLaiZiCt(_handPaiArr, _laiziData)
		//if laiZiCt>0 {
		//	logs.Info("------------------------->laiZiCt:%v ,_laiziData:%v,_handPaiArr:%v",
		//		laiZiCt,_laiziData,_handPaiArr)
		//}
		if laiZiCt > 0 { //带癞子检测胡牌

			//得到所有配牌数组
			peiPaiArr := getLzPeiPaiArr(laiZiCt, _tableType)

			_tmpCalHuInfo := GetCalHuInfo(_mjcfg, _handPaiArr)

			for _, _cvArr := range peiPaiArr {

				//获得手牌
				_calHuInfo := _tmpCalHuInfo.Clone()

				//将 1~3 张配牌添加到手牌中
				for _, _data := range _cvArr {

					if len(_mjcfg.LZExceptArr) > 0 { //如果有牌包含在LZExceptArr,则不检测
						if ContainsEqualCard(_mjcfg.LZExceptArr, _data) {
							continue
						}
					}

					_color := GetColor(_data)
					if _color == Color_Wan {
						_calHuInfo.WanIntArr = append(_calHuInfo.WanIntArr, _data)
					} else if _color == Color_Tong {
						_calHuInfo.TongIntArr = append(_calHuInfo.TongIntArr, _data)
					} else if _color == Color_Tiao {
						_calHuInfo.TiaoIntArr = append(_calHuInfo.TiaoIntArr, _data)
					} else if _color == Color_Feng {
						_calHuInfo.FengIntArr = append(_calHuInfo.FengIntArr, _data)
					} else if _color == Color_Zfb {
						_calHuInfo.ZfbIntArr = append(_calHuInfo.ZfbIntArr, _data)
					}
				}

				//排序
				_calHuInfo.SortHand()

				//检测胡牌
				isHu := _calHuInfo.ChkHu()

				if isHu {

					for _, _data := range _cvArr { //记录癞子配牌
						_calHuInfo.PeiArr = append(_calHuInfo.PeiArr, _data)
					}

					logs.Info("tableId-------------------->癞子胡牌,癞子配牌为:[%v],_calHuInfo.PxList:%v", _cvArr, _calHuInfo.PxList)
					return _calHuInfo
				}
			}

			//所有配牌都不能胡牌,返回不再运算
			return _tmpCalHuInfo

		}
	}

	//不带癞子检测
	_calHuInfo := GetCalHuInfo(_mjcfg, _handPaiArr)
	_calHuInfo.SortHand()
	isHu := _calHuInfo.ChkHu()
	if isHu {
		return _calHuInfo
	}
	return _calHuInfo
}

//检测 _array是否含有等值的 _data牌
func ContainsEqualCard(_array []int, _data int) bool {
	for _, v := range _array {
		if NewMCard(v).Equal(NewMCard(_data)) {
			return true
		}
	}
	return false
}

//听牌检测 -----------------------------------
func (chuobj *CalHuObject) CalTingPai() []*SendTip {

	//获取要检测打出的牌(去相同牌值的牌)
	allHandPai := make([]*MCard, 0)
	for _, v := range chuobj.HandPaiArr {
		for _, n := range v {
			allHandPai = append(allHandPai, n)
		}
	}

	notEqualCards := make([]int, 0)
	//for _, _card := range allHandPai {
	//	hasEqual := false
	//	for _, v := range notEqualCards {
	//		if _card.Equal(NewMCard(v)) {
	//			hasEqual = true
	//			break
	//		}
	//	}
	//	if !hasEqual {
	//		notEqualCards = append(notEqualCards, _card.data)
	//	}
	//}
	for _, _card := range allHandPai {
		notEqualCards = append(notEqualCards, _card.data)
	}
	//notEqualCards = append(notEqualCards, 28)

	length := len(notEqualCards) //最多 len(14) 个听牌提示
	logs.Info("SendCheckTing-----------> len:%v,notEqualCards:%v", length, notEqualCards)
	sendMap := make(map[int][]int)
	// chkArr := make([][]int, len)
	for i := 0; i < length; i++ {
		sendData := notEqualCards[i]
		arr := make([]int, 0)
		for _, v := range allHandPai {
			if v.GetData() != sendData {
				arr = append(arr, v.GetData())
			}
		}
		sendMap[sendData] = arr //保存
	}

	_sendTipArr := make([]*SendTip, 0)

	keys := make([]int, 0, len(sendMap))
	for k, _ := range sendMap {
		keys = append(keys, k)
	}

	ql := qo.New()
	var wg sync.WaitGroup
	// 并发不能试用 range 遍历 map，会改变遍历的过程
	// for _sendData, _chkArr := range sendMap {
	for i, lens := 0, len(keys); i < lens; i++ {
		_sendData := keys[i]
		_chkArr := sendMap[_sendData]
		wg.Add(1)
		go func() {
			sendTip := getSendTip(_sendData, _chkArr, chuobj)
			//logs.Info("sendTip==============%v=%v=%v=%v",sendTip,_sendData,_chkArr,chuobj)
			if sendTip != nil {
				//添加一个听牌提示
				ql.Go(func() {
					_sendTipArr = append(_sendTipArr, sendTip)
					wg.Done()
				})
			} else {
				wg.Done()
			}
		}()
	}
	wg.Wait()
	return _sendTipArr
}

func getSendTip(_sendData int, _chkArr []int, _chuobj *CalHuObject) *SendTip {

	_tableType := _chuobj.MjCfg.TableType
	_check7Dui := _chuobj.MjCfg.Check7Dui
	_maxCardColorIndex := _chuobj.MjCfg.MaxColorCt

	logs.Info("SendCheckTing-----------> sendCard:%v,_chkArr:%v  _check7Dui:%v",
		NewMCard(_sendData).String(), _chkArr, _check7Dui)

	//所有可能胡的牌
	huCards := make([]int, 0)
	HuInfos := make([]CalHuInfo, 0)

	_peiPaiArr := GetTableTypePai(_tableType)
	for _, data := range _peiPaiArr {
		//配牌中不能有癞子牌
		if _chuobj.MjCfg.LaiziData != DefaultIndex && NewMCard(data).EqualByData(_chuobj.MjCfg.LaiziData) {
			continue
		}
		//logs.Info("SendCheckTing----------->添加%v检测 ", NewMCard(data).String())
		_copyChkArr := IntSliceCopy(_chkArr)
		_copyChkArr = append(_copyChkArr, data)

		//设置检测手牌
		_handPaiArr := make([][]*MCard, _maxCardColorIndex)
		for _, v := range _copyChkArr {
			_card := NewMCard(v)
			_color := _card.color
			_handPaiArr[_color] = append(_handPaiArr[_color], _card)
		}

		//检测胡牌
		_calHuInfo := CalHuPai(_chuobj.MjCfg, _handPaiArr)

		if _calHuInfo.PxType > PXTYPE_UNKNOW {
			huCards = append(huCards, data) //记录这张牌
			HuInfos = append(HuInfos, *_calHuInfo)
		}
	}

	if len(huCards) > 0 {
		sendTip := NewSendTip()
		//打出去能听得牌
		sendTip.SendCard = _sendData
		for i := 0; i < len(huCards); i++ {
			sendTip.HuCards = append(sendTip.HuCards, huCards[i])
			sendTip.HuInfos = append(sendTip.HuInfos, HuInfos[i])
		}
		//添加一个听牌提示
		// _sendTipArr = append(_sendTipArr, sendTip)
		//logs.Info("sendTip========%v", sendTip)
		return sendTip

	}
	return nil
}
