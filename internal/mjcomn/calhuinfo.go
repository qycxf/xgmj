package mjcomn

import (
	"qianuuu.com/lib/logs"
)

//需要进行胡牌运算的数据对象
type CalHuInfo struct {
	Check7Dui bool //是否检测七对

	WanIntArr  []int //手牌 -万
	TongIntArr []int //手牌 -筒
	TiaoIntArr []int //手牌 -条
	FengIntArr []int //手牌 -风
	ZfbIntArr  []int //手牌 -中发白

	PxType int     //胡牌牌型 (七对\平胡)
	PxList [][]int //记录胡牌后 手牌中的 AAA 和 ABC 的牌型 按花色存储
	PeiArr []int   //如果胡牌,记录癞子所配的牌值
}

func (ch *CalHuInfo) Clone() *CalHuInfo {
	_calHuInfo := &CalHuInfo{
		Check7Dui:  ch.Check7Dui,
		WanIntArr:  IntSliceCopy(ch.WanIntArr),
		TongIntArr: IntSliceCopy(ch.TongIntArr),
		TiaoIntArr: IntSliceCopy(ch.TiaoIntArr),
		FengIntArr: IntSliceCopy(ch.FengIntArr),
		ZfbIntArr:  IntSliceCopy(ch.ZfbIntArr),
		PxList:     make([][]int, len(ch.PxList)), //记录胡牌后 手牌中的 AAA 和 ABC 的牌型 按花色存储
		PeiArr:     make([]int, 0),
	}
	return _calHuInfo
}

//排序手牌
func (ch *CalHuInfo) SortHand() {
	ch.WanIntArr = SortIntArrAsc(ch.WanIntArr)
	ch.TongIntArr = SortIntArrAsc(ch.TongIntArr)
	ch.TiaoIntArr = SortIntArrAsc(ch.TiaoIntArr)
	ch.FengIntArr = SortIntArrAsc(ch.FengIntArr)
	ch.ZfbIntArr = SortIntArrAsc(ch.ZfbIntArr)
}

//清空 pxList,在检测一种情况不符合后,应清空记录的 牌型列表,防止重复,如果胡牌 cmaj.PxList 长度 = 手牌长度-2
//_color:只删除 与 _color 相同的牌组
func (ch *CalHuInfo) ClearPxColor(_color int) {
	ch.PxList[_color] = make([]int, 0)
}

//记录到pxList 中去
func (ch *CalHuInfo) RecToPxList(i1 int, i2 int, i3 int) {

	//相同的一组牌不重复加入
	arr := make([]int, 3)
	arr[0] = i1
	arr[1] = i2
	arr[2] = i3

	_color := NewMCard(i1).GetColor()

	//isFind := false
	//for i := 0; i < len(ch.PxList[_color]); i++ {
	//	if arr[0] == ch.PxList[_color][i] {
	//		isFind = true
	//		break
	//	}
	//}
	//
	//if isFind {
	//	return
	//}

	ch.PxList[_color] = append(ch.PxList[_color], arr[0])
	ch.PxList[_color] = append(ch.PxList[_color], arr[1])
	ch.PxList[_color] = append(ch.PxList[_color], arr[2])

}

//从 PxList中 找到含有 _data 的那一组,一般用于胡牌牌型检测
func GetPxArr(_pxList [][]int, _data int) []int {

	for _, v := range _pxList {
		for i := 0; i < len(v); i = i + 3 {
			if v[0] == _data || v[1] == _data || v[2] == _data {
				return []int{v[0], v[1], v[2]}
			}
		}
	}
	return []int{}
}

var TestCt = 0

// 单次胡牌检测 -----------------------------------
func (ch *CalHuInfo) ChkHu() bool {

	//logs.Info("----------------------------------------------------->ch.Check7Dui:%v", ch.Check7Dui)
	TestCt++
	if ch.Check7Dui {
		handCt := len(ch.WanIntArr) + len(ch.TongIntArr) + len(ch.TiaoIntArr) + len(ch.FengIntArr) + len(ch.ZfbIntArr)
		if handCt == 14 {
			isQiDui := ch.Check_7DUI()
			if isQiDui {
				ch.PxType = PXTYPE_7DUI
				return true
			}
		}
	}

	if ch.Check_SSY() {
		ch.PxType = PXTYPE_SSY
		return true
	}

	if ch.Check_ShiSanLan() {
		ch.PxType = PXTYPE_SSL
		return true
	}
	isPingHu := ch.Check_PingHU()
	if isPingHu {
		ch.PxType = PXTYPE_PINGHU
		return true
	}

	return false
}

// 平胡检测 -----------------------------------
func (ch *CalHuInfo) Check_PingHU() bool {

	//logs.Info("Check_PingHU---------------[%v] [%v] [%v] [%v] [%v]",
	//	ch.WanIntArr, ch.TongIntArr, ch.TiaoIntArr, ch.FengIntArr, ch.ZfbIntArr)

	iJiangNum := 0 // 将数量

	_handIntArr := [][]int{ch.WanIntArr, ch.TongIntArr, ch.TiaoIntArr, ch.FengIntArr, ch.ZfbIntArr}
	for i := 0; i < len(_handIntArr); i++ {
		tmpPaiArr := _handIntArr[i]

		iSize := len(tmpPaiArr)

		if iSize > 0 {
			v1 := -1
			v2 := -1
			v3 := -1
			v4 := -1
			v5 := -1
			v6 := -1
			v7 := -1
			v8 := -1
			v9 := -1
			v10 := -1
			v11 := -1
			v12 := -1
			v13 := -1
			v14 := -1

			v1 = tmpPaiArr[0]
			if iSize > 1 {
				v2 = tmpPaiArr[1]
			}
			if iSize > 2 {
				v3 = tmpPaiArr[2]
			}
			if iSize > 3 {
				v4 = tmpPaiArr[3]
			}
			if iSize > 4 {
				v5 = tmpPaiArr[4]
			}
			if iSize > 5 {
				v6 = tmpPaiArr[5]
			}
			if iSize > 6 {
				v7 = tmpPaiArr[6]
			}
			if iSize > 7 {
				v8 = tmpPaiArr[7]
			}
			if iSize > 8 {
				v9 = tmpPaiArr[8]
			}
			if iSize > 9 {
				v10 = tmpPaiArr[9]
			}
			if iSize > 10 {
				v11 = tmpPaiArr[10]
			}
			if iSize > 11 {
				v12 = tmpPaiArr[11]
			}
			if iSize > 12 {
				v13 = tmpPaiArr[12]
			}
			if iSize > 13 {
				v14 = tmpPaiArr[13]
			}

			if iSize == 2 {
				if !CheckAAPai(v1, v2) {
					return false
				} else {
					iJiangNum++
				}
			} else if iSize == 3 {
				if !ch.Check3Pai(v1, v2, v3) {
					return false
				}

			} else if iSize == 5 {
				if !ch.Check5Pai(v1, v2, v3, v4, v5) {
					return false
				} else {
					iJiangNum++
				}
			} else if iSize == 6 {
				if !ch.Check6Pai(v1, v2, v3, v4, v5, v6) {
					return false
				}
			} else if iSize == 8 {
				if !ch.Check8Pai(v1, v2, v3, v4, v5, v6, v7, v8) {
					return false
				} else {
					iJiangNum++
				}
			} else if iSize == 9 {
				if !ch.Check9Pai(v1, v2, v3, v4, v5, v6, v7, v8, v9) {
					return false
				}
			} else if iSize == 11 {
				if !ch.Check11Pai(v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11) {
					return false
				} else {
					iJiangNum++
				}
			} else if iSize == 12 {
				if !ch.Check12Pai(v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12) {
					return false
				}
			} else if iSize == 14 {
				if !ch.Check14Pai(v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14) {
					return false
				} else {
					iJiangNum++
				}
			} else {
				return false
			}
		}
	}

	if iJiangNum == 1 {
		//logs.Info("Check_PingHU---------------可以平胡:>%v", ch.HandIntArr)
		return true
	}

	return false

}

//检测七对牌型
func (ch *CalHuInfo) Check_7DUI() bool {

	_handIntArr := [][]int{ch.WanIntArr, ch.TongIntArr, ch.TiaoIntArr, ch.FengIntArr, ch.ZfbIntArr}
	for _, v := range _handIntArr {
		length := len(v)
		if length%2 != 0 {
			return false
		}

		for i := 0; i < length; i = i + 2 {
			if GetColor(v[i]) != GetColor(v[i+1]) || GetVal(v[i]) != GetVal(v[i+1]) {
				return false
			}
		}
	}

	// logs.Info("七对胡牌---------------_handIntArr:>%v", _handIntArr)
	return true
}

//检测十三幺胡牌
func (ch *CalHuInfo) Check_SSY() bool {

	handCt := len(ch.WanIntArr) + len(ch.TongIntArr) + len(ch.TiaoIntArr) + len(ch.FengIntArr) + len(ch.ZfbIntArr)
	if handCt != 14 {
		return false
	}
	oneCount_wan := 0
	nineCount_wan := 0
	for _, v := range ch.WanIntArr {
		if GetVal(v) == 1 {
			oneCount_wan++
		}
		if GetVal(v) == 9 {
			nineCount_wan++
		}
	}

	oneCount_tong := 0
	nineCount_tong := 0
	for _, v := range ch.TongIntArr {
		if GetVal(v) == 1 {
			oneCount_tong++
		}
		if GetVal(v) == 9 {
			nineCount_tong++
		}
	}
	oneCount_tiao := 0
	nineCount_tiao := 0
	for _, v := range ch.TiaoIntArr {
		if GetVal(v) == 1 {
			oneCount_tiao++
		}
		if GetVal(v) == 9 {
			nineCount_tiao++
		}
	}

	dongCt := 0
	nanCt := 0
	xiCt := 0
	beiCt := 0
	for _, v := range ch.FengIntArr {
		if GetVal(v) == 1 {
			dongCt++
		}
		if GetVal(v) == 2 {
			nanCt++
		}
		if GetVal(v) == 3 {
			xiCt++
		}
		if GetVal(v) == 4 {
			beiCt++
		}
	}
	zhongCt := 0
	faCt := 0
	baiCt := 0
	for _, v := range ch.ZfbIntArr {
		if GetVal(v) == 1 {
			zhongCt++
		}
		if GetVal(v) == 2 {
			faCt++
		}
		if GetVal(v) == 3 {
			baiCt++
		}
	}

	if oneCount_wan >= 1 && nineCount_wan >= 1 &&
		oneCount_tong >= 1 && nineCount_tong >= 1 &&
		oneCount_tiao >= 1 && nineCount_tiao >= 1 &&
		dongCt >= 1 && nanCt >= 1 && xiCt >= 1 && beiCt >= 1 &&
		zhongCt >= 1 && faCt >= 1 && baiCt >= 1 {
		return true
	}
	return false
}

//十三不搭（七星不靠）
func (ch *CalHuInfo) Check_ShiSanLan() bool {
	//ch.SortHand()

	//手牌总长度
	size := len(ch.WanIntArr) + len(ch.TongIntArr) + len(ch.TiaoIntArr) + len(ch.FengIntArr) + len(ch.ZfbIntArr)
	if size != 14 {
		return false
	}

	wan := SortIntArrAsc(ch.WanIntArr)
	tong := SortIntArrAsc(ch.TongIntArr)
	tiao := SortIntArrAsc(ch.TiaoIntArr)

	wanArr, tongArr, tiaoArr := CheckQIBuKao(wan, tong, tiao)
	if len(ch.FengIntArr)+len(ch.ZfbIntArr) == 7 {

	} else {
		return false
	}

	if len(ch.FengIntArr) > 1 {
		for i := 0; i < len(ch.FengIntArr)-1; i++ {
			if NewMCard(ch.FengIntArr[i]).Equal(NewMCard(ch.FengIntArr[i+1])) {
				return false
			}
		}
	}

	if len(ch.ZfbIntArr) > 1 {
		for i := 0; i < len(ch.ZfbIntArr)-1; i++ {
			if NewMCard(ch.ZfbIntArr[i]).Equal(NewMCard(ch.ZfbIntArr[i+1])) {
				return false
			}
		}
	}

	if wanArr > 0 && tongArr > 0 && tiaoArr > 0 {

		if wanArr == 1 && tongArr == 2 && tiaoArr == 3 {
			return true

		} else if wanArr == 1 && tongArr == 3 && tiaoArr == 2 {
			return true
		} else if wanArr == 2 && tongArr == 1 && tiaoArr == 3 {
			return true
		} else if wanArr == 2 && tongArr == 3 && tiaoArr == 1 {
			return true
		} else if wanArr == 3 && tongArr == 1 && tiaoArr == 2 {
			return true
		} else if wanArr == 3 && tongArr == 2 && tiaoArr == 1 {
			return true
		} else {
			return false
		}
	}

	return false
}

func CheckQIBuKao(wan []int, tong []int, tiao []int) (int, int, int) {

	wanArr := 0
	tongArr := 0
	tiaoArr := 0
	if len(wan) == 3 {
		if GetVal(wan[0]) == 1 {
			if GetVal(wan[1]) == 4 {
				if GetVal(wan[2]) == 7 {
					wanArr = 1
				} else {
					return 0, 0, 0
				}

			} else {
				return 0, 0, 0
			}

		} else if GetVal(wan[0]) == 2 {
			if GetVal(wan[1]) == 5 {
				if GetVal(wan[2]) == 8 {
					wanArr = 2
				} else {
					return 0, 0, 0
				}
			} else {
				return 0, 0, 0
			}
		} else if GetVal(wan[0]) == 3 {
			if GetVal(wan[1]) == 6 {
				if GetVal(wan[2]) == 9 {
					wanArr = 3
				} else {
					return 0, 0, 0
				}
			} else {
				return 0, 0, 0
			}
		} else {
			return 0, 0, 0
		}

	} else if len(wan) == 2 {
		if GetVal(wan[0]) == 1 {
			if GetVal(wan[1]) == 4 || GetVal(wan[1]) == 7 {
				wanArr = 1
			} else {
				return 0, 0, 0
			}

		} else if GetVal(wan[0]) == 2 {
			if GetVal(wan[1]) == 5 || GetVal(wan[1]) == 8 {
				wanArr = 2
			} else {
				return 0, 0, 0
			}
		} else if GetVal(wan[0]) == 3 {
			if GetVal(wan[1]) == 6 || GetVal(wan[1]) == 9 {
				wanArr = 3
			} else {
				return 0, 0, 0
			}
		} else {
			return 0, 0, 0
		}
	} else {
		return 0, 0, 0
	}

	if len(tong) == 3 {
		if GetVal(tong[0]) == 1 {
			if GetVal(tong[1]) == 4 {
				if GetVal(tong[2]) == 7 {
					tongArr = 1
				} else {
					return 0, 0, 0
				}

			} else {
				return 0, 0, 0
			}

		} else if GetVal(tong[0]) == 2 {
			if GetVal(tong[1]) == 5 {
				if GetVal(tong[2]) == 8 {
					tongArr = 2
				} else {
					return 0, 0, 0
				}
			} else {
				return 0, 0, 0
			}
		} else if GetVal(tong[0]) == 3 {
			if GetVal(tong[1]) == 6 {
				if GetVal(tong[2]) == 9 {
					tongArr = 3
				} else {
					return 0, 0, 0
				}
			} else {
				return 0, 0, 0
			}
		} else {
			return 0, 0, 0
		}

	} else if len(tong) == 2 {
		if GetVal(tong[0]) == 1 {
			if GetVal(tong[1]) == 4 || GetVal(tong[1]) == 7 {
				tongArr = 1
			} else {
				return 0, 0, 0
			}

		} else if GetVal(tong[0]) == 2 {
			if GetVal(tong[1]) == 5 || GetVal(tong[1]) == 8 {
				tongArr = 2
			} else {
				return 0, 0, 0
			}
		} else if GetVal(tong[0]) == 3 {
			if GetVal(tong[1]) == 6 || GetVal(tong[1]) == 9 {
				tongArr = 3
			} else {
				return 0, 0, 0
			}
		} else {
			return 0, 0, 0
		}
	} else {
		return 0, 0, 0
	}

	if len(tiao) == 3 {
		if GetVal(tiao[0]) == 1 {
			if GetVal(tiao[1]) == 4 {
				if GetVal(tiao[2]) == 7 {
					tiaoArr = 1
				} else {
					return 0, 0, 0
				}

			} else {
				return 0, 0, 0
			}

		} else if GetVal(tiao[0]) == 2 {
			if GetVal(tiao[1]) == 5 {
				if GetVal(tiao[2]) == 8 {
					tiaoArr = 2
				} else {
					return 0, 0, 0
				}
			} else {
				return 0, 0, 0
			}
		} else if GetVal(tiao[0]) == 3 {
			if GetVal(tiao[1]) == 6 {
				if GetVal(tiao[2]) == 9 {
					tiaoArr = 3
				} else {
					return 0, 0, 0
				}
			} else {
				return 0, 0, 0
			}
		} else {
			return 0, 0, 0
		}

	} else if len(tiao) == 2 {
		if GetVal(tiao[0]) == 1 {
			if GetVal(tiao[1]) == 4 || GetVal(tiao[1]) == 7 {
				tiaoArr = 1
			} else {
				return 0, 0, 0
			}

		} else if GetVal(tiao[0]) == 2 {
			if GetVal(tiao[1]) == 5 || GetVal(tiao[1]) == 8 {
				tiaoArr = 2
			} else {
				return 0, 0, 0
			}
		} else if GetVal(tiao[0]) == 3 {
			if GetVal(tiao[1]) == 6 || GetVal(tiao[1]) == 9 {
				tiaoArr = 3
			} else {
				return 0, 0, 0
			}
		} else {
			return 0, 0, 0
		}
	} else {
		return 0, 0, 0
	}
	return wanArr, tongArr, tiaoArr
}

//检测牌型 [豪华七对]：牌型为七对，但玩家手中有四张一样的牌（没有杠出）
func (ch *CalHuInfo) Check_HAOHUA7DUI() bool {

	_handIntArr := [][]int{ch.WanIntArr, ch.TongIntArr, ch.TiaoIntArr}
	logs.Info("豪华七对---------------_handIntArr:>%v", _handIntArr)

	for _, v := range _handIntArr {
		iSize := len(v)
		if iSize%2 != 0 {
			return false
		}
		if iSize >= 4 {
			count := 0
			value := -1
			for _, n := range v {
				val := GetVal(n)
				if val != value {
					count = 1
					value = val
				} else {
					count++
				}
				//logs.Info("豪华七对---------------count:>%v", count)
				if count == 4 {
					return true
				}
			}
		}
	}

	return false
}

//检测牌型 [超豪华七对]：牌型为七对，但玩家手中2组 四张一样的牌（没有杠出）
func (ch *CalHuInfo) Check_SUPERHAOHUA7DUI() bool {

	_handIntArr := [][]int{ch.WanIntArr, ch.TongIntArr, ch.TiaoIntArr}
	ct := 0
	for _, v := range _handIntArr {
		iSize := len(v)
		if iSize%2 != 0 {
			return false
		}
		if iSize >= 4 {
			count := 0
			value := -1
			for _, n := range v {
				val := GetVal(n)
				if val != value {
					count = 1
					value = val
				} else {
					count++
				}

				if count == 4 {
					ct++
				}
			}
		}
	}

	if ct >= 2 {
		return true
	}
	return false
}

//检测牌型 [超超豪华七对]：牌型为七对，但玩家手中3组 四张一样的牌（没有杠出）
func (ch *CalHuInfo) Check_SUPER_SUPERHAOHUA7DUI() bool {

	_handIntArr := [][]int{ch.WanIntArr, ch.TongIntArr, ch.TiaoIntArr}
	ct := 0
	for _, v := range _handIntArr {
		iSize := len(v)
		if iSize%2 != 0 {
			return false
		}
		if iSize >= 4 {
			count := 0
			value := -1
			for _, n := range v {
				val := GetVal(n)
				if val != value {
					count = 1
					value = val
				} else {
					count++
				}

				if count == 4 {
					ct++
				}
			}
		}
	}

	if ct >= 3 {
		return true
	}
	return false
}

// 检测将牌（2张）
func CheckAAPai(i1 int, i2 int) bool {
	return GetVal(i1) == GetVal(i2)
}

// 检测是否胡牌（3张）
func (ch *CalHuInfo) Check3Pai(iValue1 int, iValue2 int, iValue3 int) bool {
	if ch.CheckABCPai(iValue1, iValue2, iValue3) {
		return true
	}

	if ch.CheckAAAPai(iValue1, iValue2, iValue3) {
		return true
	}
	return false
}

// 检测是否胡牌（14张）
func (ch *CalHuInfo) Check14Pai(iValue1 int, iValue2 int, iValue3 int, iValue4 int, iValue5 int,
	iValue6 int, iValue7 int, iValue8 int, iValue9 int, iValue10 int, iValue11 int, iValue12 int, iValue13 int, iValue14 int) bool {

	// 如果是左边两个为将，右边为三重张或三连张
	if CheckAAPai(iValue1, iValue2) {
		// 无AAA，全ABC
		if ch.Check12Pai(iValue3, iValue4, iValue5, iValue6, iValue7,
			iValue8, iValue9, iValue10, iValue11, iValue12, iValue13,
			iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是中间两个为将，左右边为三重张或三连张
	if CheckAAPai(iValue4, iValue5) {
		// 无AAA，全ABC
		if ch.Check3Pai(iValue1, iValue2, iValue3) && ch.Check9Pai(iValue6, iValue7, iValue8, iValue9, iValue10,
			iValue11, iValue12, iValue13, iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是中间两个为将，左右边为三重张或三连张
	if CheckAAPai(iValue7, iValue8) {
		// 无AAA，全ABC
		if ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue5, iValue6) && ch.Check6Pai(iValue9, iValue10, iValue11, iValue12,
			iValue13, iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是中间两个为将，左右边为三重张或三连张
	if CheckAAPai(iValue10, iValue11) {
		// 无AAA，全ABC
		if ch.Check3Pai(iValue12, iValue13, iValue14) && ch.Check9Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
			iValue6, iValue7, iValue8, iValue9) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是右边两个为将，左右边为三重张或三连张
	if CheckAAPai(iValue12, iValue13) {
		// 无AAA，全ABC
		if ch.Check12Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
			iValue6, iValue7, iValue8, iValue9, iValue10, iValue11, iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 13,14为将
	if CheckAAPai(iValue13, iValue14) {
		// 无AAA，全ABC
		if ch.Check12Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
			iValue6, iValue7, iValue8, iValue9, iValue10, iValue11, iValue12) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况 23
	if CheckAAPai(iValue2, iValue3) {
		// 无AAA，全ABC
		if ch.Check12Pai(iValue1, iValue4, iValue5, iValue6, iValue7,
			iValue8, iValue9, iValue10, iValue11, iValue12, iValue13, iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况 56
	if CheckAAPai(iValue5, iValue6) {
		// 无AAA，全ABC
		if ch.Check12Pai(iValue1, iValue2, iValue3, iValue4, iValue7,
			iValue8, iValue9, iValue10, iValue11, iValue12, iValue13, iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况 89
	if CheckAAPai(iValue8, iValue9) {
		// 无AAA，全ABC
		if ch.Check12Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
			iValue6, iValue7, iValue10, iValue11, iValue12, iValue13, iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况 11 12
	if CheckAAPai(iValue11, iValue12) {
		// 无AAA，全ABC
		if ch.Check12Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
			iValue6, iValue7, iValue8, iValue9, iValue10, iValue13, iValue14) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	return false
}

// 检测是否胡牌（12张）
func (ch *CalHuInfo) Check12Pai(iValue1 int, iValue2 int, iValue3 int, iValue4 int, iValue5 int,
	iValue6 int, iValue7 int, iValue8 int, iValue9 int, iValue10 int, iValue11 int, iValue12 int) bool {

	if ch.Check3Pai(iValue1, iValue2, iValue3) && ch.Check9Pai(iValue4, iValue5, iValue6, iValue7, iValue8,
		iValue9, iValue10, iValue11, iValue12) {
		return true
	}
	if ch.Check3Pai(iValue10, iValue11, iValue12) && ch.Check9Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
		iValue6, iValue7, iValue8, iValue9) {
		return true
	}
	if ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue5, iValue6) && ch.Check6Pai(iValue7, iValue8, iValue9, iValue10, iValue11,
		iValue12) {
		return true
	}
	//                123456789
	//一种特殊情况的牌  ABBCCCDDDEEF
	//if Check3Pai(iValue1, iValue2, iValue4) && Check3Pai(iValue3, iValue5, iValue7) &&
	//	Check3Pai(iValue6, iValue8, iValue10) && Check3Pai(iValue9, iValue11, iValue12) {
	//	return true
	//}
	if ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue5, iValue7) &&
		ch.Check6Pai(iValue6, iValue8, iValue9, iValue10, iValue11, iValue12) {
		return true
	}

	////一种特殊情况的牌  AABBBCCCCDDE
	if ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue6, iValue7) &&
		ch.Check6Pai(iValue5, iValue8, iValue9, iValue10, iValue11, iValue12) {
		return true
	}

	return false
}

// 检测是否胡牌（11张）
func (ch *CalHuInfo) Check11Pai(iValue1 int, iValue2 int, iValue3 int, iValue4 int, iValue5 int,
	iValue6 int, iValue7 int, iValue8 int, iValue9 int, iValue10 int, iValue11 int) bool {

	// 如果是左边两个为将
	if CheckAAPai(iValue1, iValue2) {
		if ch.Check9Pai(iValue3, iValue4, iValue5, iValue6, iValue7, iValue8,
			iValue9, iValue10, iValue11) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况 左边
	if CheckAAPai(iValue2, iValue3) {
		if ch.Check9Pai(iValue1, iValue4, iValue5, iValue6, iValue7, iValue8,
			iValue9, iValue10, iValue11) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况 中间
	if CheckAAPai(iValue5, iValue6) {
		if ch.Check9Pai(iValue1, iValue2, iValue3, iValue4, iValue7, iValue8,
			iValue9, iValue10, iValue11) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况 右边
	if CheckAAPai(iValue8, iValue9) {
		if ch.Check9Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
			iValue6, iValue7, iValue10, iValue11) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是中间两个为将
	if CheckAAPai(iValue4, iValue5) {
		// 无AAA，全ABC
		if ch.Check3Pai(iValue1, iValue2, iValue3) &&
			ch.Check6Pai(iValue6, iValue7, iValue8, iValue9, iValue10,
				iValue11) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是右边两个为将
	if CheckAAPai(iValue7, iValue8) {
		// 无AAA，全ABC
		if ch.Check3Pai(iValue9, iValue10, iValue11) &&
			ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
				iValue6) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是右边两个为将
	if CheckAAPai(iValue10, iValue11) {
		if ch.Check9Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
			iValue6, iValue7, iValue8, iValue9) {
			return true
		}

		ch.ClearPxColor(GetColor(iValue1))
	}

	return false
}

// 检测是否胡牌（9张）
func (ch *CalHuInfo) Check9Pai(iValue1 int, iValue2 int, iValue3 int, iValue4 int, iValue5 int, iValue6 int, iValue7 int, iValue8 int, iValue9 int) bool {

	if ch.Check3Pai(iValue1, iValue2, iValue3) && ch.Check6Pai(iValue4, iValue5, iValue6, iValue7, iValue8,
		iValue9) {
		return true
	}
	if ch.Check3Pai(iValue7, iValue8, iValue9) && ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue5,
		iValue6) {
		return true
	}

	//一种特殊情况 AABBBCCCD
	if ch.Check3Pai(iValue5, iValue8, iValue9) && ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue6,
		iValue7) {
		return true
	}

	//logs.Info("------------------->")
	//一种特殊情况 ABBCCCDDE
	if ch.Check3Pai(iValue1, iValue2, iValue4) && ch.Check3Pai(iValue3, iValue5, iValue7) && ch.Check3Pai(iValue6, iValue8, iValue9) {
		return true
	}

	//一种特殊情况 ABBBCCCDD
	if ch.Check3Pai(iValue1, iValue2, iValue5) && ch.Check6Pai(iValue3, iValue4, iValue6, iValue7, iValue8,
		iValue9) {
		return true
	}

	return false
}

// 检测是否胡牌（8张）
func (ch *CalHuInfo) Check8Pai(iValue1 int, iValue2 int, iValue3 int, iValue4 int, iValue5 int,
	iValue6 int, iValue7 int, iValue8 int) bool {
	// 如果是左边两个为将，右边为三重张或三连张
	if CheckAAPai(iValue1, iValue2) {
		if ch.Check6Pai(iValue3, iValue4, iValue5, iValue6, iValue7, iValue8) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//左边 ABBBC的情况
	if CheckAAPai(iValue2, iValue3) {
		if ch.Check6Pai(iValue1, iValue4, iValue5, iValue6, iValue7, iValue8) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是中间两个为将，左右边为三重张或三连张
	if CheckAAPai(iValue4, iValue5) {
		if ch.Check3Pai(iValue1, iValue2, iValue3) && ch.Check3Pai(iValue6, iValue7, iValue8) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是右边两个为将，左边为三重张或三连张
	if CheckAAPai(iValue7, iValue8) {
		if ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue5, iValue6) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//右边边 ABBBC的情况
	if CheckAAPai(iValue5, iValue6) {
		if ch.Check6Pai(iValue1, iValue2, iValue3, iValue4, iValue7, iValue8) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	//有一种暗杠的牌形式 如 5 666 7777
	if CheckAAPai(iValue2, iValue3) {
		if ch.Check6Pai(iValue1, iValue4, iValue5, iValue6, iValue7, iValue8) {
			return true
		}
		ch.ClearPxColor(GetColor(iValue1))
	}

	return false
}

// 检测是否胡牌（6张）
func (ch *CalHuInfo) Check6Pai(iValue1 int, iValue2 int, iValue3 int, iValue4 int, iValue5 int, iValue6 int) bool {

	//AAABBB \AAAABC
	if ch.Check3Pai(iValue1, iValue2, iValue3) && ch.Check3Pai(iValue4, iValue5, iValue6) {
		return true
	}
	//ABBCCD
	if ch.CheckABBCCD(iValue1, iValue2, iValue3, iValue4, iValue5, iValue6) {
		return true
	}
	//AABBCC
	if ch.CheckAABBCCPai(iValue1, iValue2, iValue3, iValue4, iValue5, iValue6) {
		return true
	}
	//ABBBBC 的形式
	if ch.Check3Pai(iValue1, iValue2, iValue6) && ch.Check3Pai(iValue3, iValue4, iValue5) {
		return true
	}

	return false
}

// 检测是否胡牌（5张）
func (ch *CalHuInfo) Check5Pai(iValue1 int, iValue2 int, iValue3 int, iValue4 int, iValue5 int) bool {

	// 如果是左边两个为将，右边为三重张或三连张
	if CheckAAPai(iValue1, iValue2) {
		if ch.Check3Pai(iValue3, iValue4, iValue5) {
			return true
		}

		ch.ClearPxColor(GetColor(iValue1))
	}

	// 如果是右边两个为将，左边为三重张或三连张
	if CheckAAPai(iValue4, iValue5) {

		if ch.Check3Pai(iValue1, iValue2, iValue3) {
			return true
		}

		ch.ClearPxColor(GetColor(iValue1))
	}

	//ABBBC的情况
	if CheckAAPai(iValue2, iValue3) {
		if ch.Check3Pai(iValue1, iValue4, iValue5) {
			return true
		}

		ch.ClearPxColor(GetColor(iValue1))
	}

	return false
}

// 检测是否六连对
func (ch *CalHuInfo) CheckAABBCCDDEEFFPai(i1 int, i2 int, i3 int, i4 int, i5 int, i6 int, i7 int, i8 int, i9 int,
	i10 int, i11 int, i12 int) bool {
	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	iValue4 := GetVal(i4)
	iValue5 := GetVal(i5)
	iValue6 := GetVal(i6)
	iValue7 := GetVal(i7)
	iValue8 := GetVal(i8)
	iValue9 := GetVal(i9)
	iValue10 := GetVal(i10)
	iValue11 := GetVal(i11)
	iValue12 := GetVal(i12)

	if iValue1 == iValue2 && iValue3 == iValue4 && iValue5 == iValue6 &&
		iValue7 == iValue8 && iValue9 == iValue10 &&
		iValue11 == iValue12 {
		if (iValue1 == iValue3-1) && (iValue3 == iValue5-1) &&
			(iValue5 == iValue7-1) && (iValue7 == iValue9-1) &&
			(iValue9 == iValue11-1) {
			ch.RecToPxList(i1, i3, i5)
			ch.RecToPxList(i2, i4, i6)
			ch.RecToPxList(i7, i9, i11)
			ch.RecToPxList(i8, i10, i12)
			return true
		}
	}
	return false

}

// 检测是否三连刻
func (ch *CalHuInfo) CheckAAAABBBBCCCCPai(i1 int, i2 int, i3 int, i4 int, i5 int, i6 int, i7 int, i8 int, i9 int,
	i10 int, i11 int, i12 int) bool {
	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	iValue4 := GetVal(i4)
	iValue5 := GetVal(i5)
	iValue6 := GetVal(i6)
	iValue7 := GetVal(i7)
	iValue8 := GetVal(i8)
	iValue9 := GetVal(i9)
	iValue10 := GetVal(i10)
	iValue11 := GetVal(i11)
	iValue12 := GetVal(i12)

	if (iValue1 == iValue2 && iValue2 == iValue3 && iValue3 == iValue4) &&
		(iValue5 == iValue6 && iValue6 == iValue7 && iValue7 == iValue8) &&
		(iValue9 == iValue10 && iValue10 == iValue11 && iValue11 == iValue12) {
		if (iValue1 == iValue5-1) && (iValue5 == iValue9-1) {
			ch.RecToPxList(i1, i5, i9)
			ch.RecToPxList(i2, i6, i10)
			ch.RecToPxList(i3, i7, i11)
			ch.RecToPxList(i4, i8, i12)
			return true
		}
	}
	return false

}

// 检测是否三连高压
func (ch *CalHuInfo) CheckAAABBBCCCPai(i1 int, i2 int, i3 int, i4 int, i5 int, i6 int, i7 int, i8 int, i9 int) bool {
	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	iValue4 := GetVal(i4)
	iValue5 := GetVal(i5)
	iValue6 := GetVal(i6)
	iValue7 := GetVal(i7)
	iValue8 := GetVal(i8)
	iValue9 := GetVal(i9)

	if (iValue1 == iValue2 && iValue2 == iValue3) &&
		(iValue4 == iValue5 && iValue5 == iValue6) && (iValue7 == iValue8 && iValue8 == iValue9) {
		if (iValue1 == iValue4-1) && (iValue4 == iValue7-1) {
			ch.RecToPxList(i1, i4, i7)
			ch.RecToPxList(i2, i5, i8)
			ch.RecToPxList(i3, i6, i9)
			return true
		}
	}
	return false

}

// 检测是否三连对 AABBCC
func (ch *CalHuInfo) CheckAABBCCPai(i1 int, i2 int, i3 int, i4 int, i5 int, i6 int) bool {

	if ch.CheckABCPai(i1, i3, i5) && ch.CheckABCPai(i2, i4, i6) {
		return true
	}

	//iValue1 := GetVal(i1)
	//iValue2 := GetVal(i2)
	//iValue3 := GetVal(i3)
	//iValue4 := GetVal(i4)
	//iValue5 := GetVal(i5)
	//iValue6 := GetVal(i6)
	//
	//if iValue1 == iValue2 && iValue3 == iValue4 && iValue5 == iValue6 {
	//	if (iValue1 == iValue3-1) && (iValue3 == iValue5-1) {
	//		ch.RecToPxList(i1, i3, i5)
	//		ch.RecToPxList(i2, i4, i6)
	//		return true
	//	}
	//}

	return false

}

// 检测是否三连对 AABBCC
func (ch *CalHuInfo) CheckAABBCCPai_hy(i1 int, i2 int, i3 int, i4 int, i5 int, i6 int) bool {

	if ch.CheckABCPai_hy(i1, i2, i3) && ch.CheckABCPai_hy(i4, i5, i6) {
		return true
	}

	//iValue1 := GetVal(i1)
	//iValue2 := GetVal(i2)
	//iValue3 := GetVal(i3)
	//iValue4 := GetVal(i4)
	//iValue5 := GetVal(i5)
	//iValue6 := GetVal(i6)
	//
	//if iValue1 == iValue2 && iValue3 == iValue4 && iValue5 == iValue6 {
	//	if (iValue1 == iValue3-1) && (iValue3 == iValue5-1) {
	//		ch.RecToPxList(i1, i3, i5)
	//		ch.RecToPxList(i2, i4, i6)
	//		return true
	//	}
	//}

	return false

}

//ABBCCD
func (ch *CalHuInfo) CheckABBCCD(i1 int, i2 int, i3 int, i4 int, i5 int, i6 int) bool {

	if ch.CheckABCPai(i1, i2, i4) && ch.CheckABCPai(i3, i5, i6) {
		return true
	}

	//iValue1 := GetVal(i1)
	//iValue2 := GetVal(i2)
	//iValue3 := GetVal(i3)
	//iValue4 := GetVal(i4)
	//iValue5 := GetVal(i5)
	//iValue6 := GetVal(i6)
	//if iValue2 == iValue3 && iValue4 == iValue5 {
	//	if (iValue1 == iValue2-1) && (iValue3 == iValue4-1) && (iValue5 == iValue6-1) {
	//		ch.RecToPxList(i1, i2, i4)
	//		ch.RecToPxList(i3, i5, i6)
	//		return true
	//	}
	//}
	return false

}

// 检测是否三重张
func (ch *CalHuInfo) CheckAAAPai(i1 int, i2 int, i3 int) bool {

	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	if iValue1 == iValue2 && iValue2 == iValue3 {
		ch.RecToPxList(i1, i2, i3)
		return true
	}
	return false
}

// 检测是否三连张
func (ch *CalHuInfo) CheckABCPai(i1 int, i2 int, i3 int) bool {

	//如果是中发白\风  则不检测ABC牌型
	color := GetColor(i1)
	if color == Color_Feng || color == Color_Zfb || color == Color_Hua {
		return false
	}

	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	if iValue1 == (iValue2-1) && iValue2 == (iValue3-1) {
		ch.RecToPxList(i1, i2, i3)
		return true
	}
	return false
} // 检测是否三连张
func (ch *CalHuInfo) CheckABCPai_hy(i1 int, i2 int, i3 int) bool {

	//如果是中发白\风  则不检测ABC牌型
	color := GetColor(i1)
	if color == Color_Feng || color == Color_Zfb || color == Color_Hua {
		return false
	}

	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	if iValue1 == (iValue2-1) && iValue2 == (iValue3-1) {
		//ch.RecToPxList(i1, i2, i3)
		return true
	}
	return false
}

// 检测是否三重张
func CheckAAAPai(i1 int, i2 int, i3 int) bool {

	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	if iValue1 == iValue2 && iValue2 == iValue3 {
		return true
	}
	return false
}

// 检测是否三连张
func CheckABCPai(i1 int, i2 int, i3 int) bool {

	//如果是中发白\风  则不检测ABC牌型
	color := GetColor(i1)
	if color == Color_Feng || color == Color_Zfb {
		return false
	}

	iValue1 := GetVal(i1)
	iValue2 := GetVal(i2)
	iValue3 := GetVal(i3)
	if iValue1 == (iValue2-1) && iValue2 == (iValue3-1) {
		return true
	}
	return false
}
