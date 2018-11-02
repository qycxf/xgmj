package maj

import (
	"testing"

	"qianuuu.com/xgmj/internal/config"
)

func TestMajong(t *testing.T) {

	//测试胡牌牌型 ------------------------------------------------------------

	//var CARDDATA = [108]int{
	//            一                二               三                四                  五                   六                  七                八                  九
	// 万	0, 1, 2, 3,       4, 5, 6, 7,       8, 9, 10, 11,      12, 13, 14, 15,    16,17, 18, 19,     20,21, 22, 23,     24, 25, 26, 27,    28, 29, 30, 31,      32, 33, 34, 35, //万 0 - 35
	// 筒	36, 37, 38, 39,   40, 41, 42,43,    44, 45, 46, 47,    48, 49, 50, 51,    52, 53, 54, 55,    56, 57, 58, 59,    60, 61, 62, 63,    64, 65, 66, 67,      68, 69, 70, 71, //筒 36 - 71
	// 条	72, 73, 74, 75,   76, 77, 78, 79,   80, 81, 82, 83,    84, 85, 86, 87,    88, 89, 90, 91,    92, 93, 94, 95,    96, 97, 98, 99,    100, 101, 102, 103,  104, 105, 106, 107, //条 72 - 107
	//}

	_tableCfg := config.NewTableCfg()
	_tableCfg.TableType = TableType_HZMJ_AH
	cmaj := NewCMaj(0, _tableCfg)
	//cmaj.GetPeiPaiArr(3)

	//{"中", "中", "中", "3万", "4万", "4筒", "6筒", "7筒", "7筒", "1条", "2条", "4条", "5条", "8条", "9条"},
	arr := []int{124, 125, 126, 8, 12, 48, 56, 60, 61, 72, 76, 84, 88, 100}
	for _, v := range arr {
		card := NewMCard(v)
		cmaj.AddHandPai(card, true)
	}

	//cmaj.DoPeng(NewCard(84))
	//cmaj.DoPeng(NewCard(76))
	//cmaj.DoMingGang(NewCard(104))
	//cmaj.DoMingGang(NewCard(80))

	//isHu := cmaj.Check_PingHU()
	//logs.Info("isHu:%v", isHu)

}
