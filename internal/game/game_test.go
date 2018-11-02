package game

import (
	"testing"

	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/maj"
	"qianuuu.com/xgmj/internal/game/table"
	. "qianuuu.com/mahjong/mjcomn"
)

func TestMajong(t *testing.T) {

	_tableCfg := config.NewTableCfg()
	_tableCfg.TableType = TableType_FYMJ
	_tableCfg.KehuQidui = consts.Yes
	_tableCfg.MaxCardColorIndex = 5

	//_testTable := table.NewHZTable(1, table.NewRobots(), _tableCfg)
	_testTable := table.NewFYTable(1, table.NewRobots(), _tableCfg)

	_testTable.Majhong.CMajArr[0] = maj.NewCMaj(0, _tableCfg)
	_testTable.Majhong.CurtSenderIndex = 0

	//var CARDDATA = [108]int{
	//            一                二               三                四                  五                   六                  七                八                  九
	// 万	0, 1, 2, 3,       4, 5, 6, 7,       8, 9, 10, 11,      12, 13, 14, 15,    16,17, 18, 19,     20,21, 22, 23,     24, 25, 26, 27,    28, 29, 30, 31,      32, 33, 34, 35, //万 0 - 35
	// 筒	36, 37, 38, 39,   40, 41, 42,43,    44, 45, 46, 47,    48, 49, 50, 51,    52, 53, 54, 55,    56, 57, 58, 59,    60, 61, 62, 63,    64, 65, 66, 67,      68, 69, 70, 71, //筒 36 - 71
	// 条	72, 73, 74, 75,   76, 77, 78, 79,   80, 81, 82, 83,    84, 85, 86, 87,    88, 89, 90, 91,    92, 93, 94, 95,    96, 97, 98, 99,    100, 101, 102, 103,  104, 105, 106, 107, //条 72 - 107
	//}

	////{"中", "中", "3万", "4万", "5万", "4筒", "5筒", "6筒", "8筒", "8筒", "2条", "3条", "7条", "8条"},
	//arr := []int{124, 125, 8, 12, 16, 48, 52, 56, 64, 65, 76, 80, 96, 100}

	//{"中", "中", "中", "4万", "5万", "4筒", "5筒", "6筒", "8筒", "8筒", "2条", "3条", "7条", "8条"},
	//arr := []int{124, 125, 126, 12, 16, 48, 52, 56, 64, 65, 76, 80, 96, 100}

	//{"中", "中", "中", "1万", "3万", "5万", "9万", "1筒", "3筒", "5筒", "9筒", "1条", "5条", "9条"},
	//arr := []int{124, 125, 126, 0, 8, 16, 32, 36, 44, 52, 68, 72, 88, 104}

	//中中 7万7万1条2条2条3条4条6条8条8条9条  3条
	//arr := []int{124, 125, 24, 25, 72, 76, 77, 80, 84, 92, 100, 101, 104, 40}

	arr := []int{124, 126, 75, 80, 84, 88, 89}
	//{"中", "中", "中", "1万", "2万", "2万", "2万", "4筒", "5筒", "6筒", "4条", "4条", "4条", "6万"},

	for _, v := range arr {
		card := NewMCard(v)
		_testTable.Majhong.CMajArr[0].AddHandPai(card, true)
	}

	_testTable.CheckHuPai(0)
	//_testTable.SendCheckTing()
	//logs.Info("------------->TestCt:%v", TestCt)

}
