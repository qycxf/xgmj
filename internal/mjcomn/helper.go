package mjcomn

// 同癞子个数的各种牌型值
type countmap map[int][][]int

// 获取不同癞子数的牌型值组
var mplz map[int]countmap

var tableTypePai map[int][]int

func init() {

	tableTypePai = make(map[int][]int)
	// 牌桌类型牌定义
	tableTypePai[TableType_HFMJ] = []int{4, 8, 12, 16, 20, 24, 28, 40, 44, 48, 52, 56, 60, 64, 76, 80, 84, 88, 92, 96, 100}
	tableTypePai[TableType_HZMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	tableTypePai[TableType_YC_XLCH] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	tableTypePai[TableType_XY_KA5XING] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	tableTypePai[TableType_SC_XZDD] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	tableTypePai[TableType_FJ_FZMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	tableTypePai[TableType_NJMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104, 108, 112, 116, 120}
	tableTypePai[TableType_FYMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	tableTypePai[TableType_ASMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104, 108, 112, 116, 120, 124, 128, 132}
	tableTypePai[TableType_AQMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104, 108, 112, 116, 120, 124, 128, 132}
	tableTypePai[TableType_BBMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104, 108, 112, 116, 120, 124, 128, 132}
	tableTypePai[TableType_HYMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104, 108, 112, 116, 120, 124, 128, 132}
	tableTypePai[TableType_GBMJ] = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}

	// 初始化 赖子可能 ，三个癞子组
	mplz = make(map[int]countmap)
	for i := 1; i <= 3; i++ {
		cm := make(countmap)
		for k, v := range tableTypePai {
			cm[k] = preGetLzPeiPaiArr(i, v)
		}
		mplz[i] = cm
	}
}

func GetTableTypePai(tabletype int) []int {
	return tableTypePai[tabletype]
}

// 癞子配牌 (只读，所以不用 copy)
func getLzPeiPaiArr(_laiZiCt int, _typeType int) [][]int {
	if _laiZiCt > 3 || _laiZiCt < 1 {
		return nil
	}
	return mplz[_laiZiCt][_typeType]
}

func preGetLzPeiPaiArr(_laiZiCt int, _dataArr []int) [][]int {

	rtnArr := make([][]int, 0)
	for i := 0; i < len(_dataArr); i++ {
		_tmpArr := make([]int, 0)
		_tmpArr = append(_tmpArr, _dataArr[i])
		if _laiZiCt > 1 {
			for j := i; j < len(_dataArr); j++ {
				_tmpArr2 := make([]int, 0)
				for _, vv := range _tmpArr {
					_tmpArr2 = append(_tmpArr2, vv)
				}

				_tmpArr2 = append(_tmpArr2, _dataArr[j])
				if _laiZiCt > 2 {

					for k := j; k < len(_dataArr); k++ {
						_tmpArr3 := make([]int, 0)
						for _, vv := range _tmpArr2 {
							_tmpArr3 = append(_tmpArr3, vv)
						}

						_tmpArr3 = append(_tmpArr3, _dataArr[k])
						rtnArr = append(rtnArr, _tmpArr3)
					}

				} else {
					rtnArr = append(rtnArr, _tmpArr2)
				}
			}
		} else {
			rtnArr = append(rtnArr, _tmpArr)
		}
	}

	//logs.Info("------------------------->GetPeiPaiArr _laiZiCt:%v,len(rtnArr):%v", _laiZiCt, len(rtnArr))
	return rtnArr
}
