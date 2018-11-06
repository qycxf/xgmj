package mjcomn

import (
	"fmt"
	"testing"
)

func TestMajong(t *testing.T) {

	//var CARDDATA = [108]int{
	//            一                二               三                四                  五                   六                  七                八                  九
	// 万	0, 1, 2, 3,       4, 5, 6, 7,       8, 9, 10, 11,      12, 13, 14, 15,    16,17, 18, 19,     20,21, 22, 23,     24, 25, 26, 27,    28, 29, 30, 31,      32, 33, 34, 35, //万 0 - 35
	// 筒	36, 37, 38, 39,   40, 41, 42,43,    44, 45, 46, 47,    48, 49, 50, 51,    52, 53, 54, 55,    56, 57, 58, 59,    60, 61, 62, 63,    64, 65, 66, 67,      68, 69, 70, 71, //筒 36 - 71
	// 条	72, 73, 74, 75,   76, 77, 78, 79,   80, 81, 82, 83,    84, 85, 86, 87,    88, 89, 90, 91,    92, 93, 94, 95,    96, 97, 98, 99,    100, 101, 102, 103,  104, 105, 106, 107, //条 72 - 107
	//}

	//{"中", "中", "中", "4万", "5万", "4筒", "5筒", "6筒", "8筒", "8筒", "2条", "3条", "7条", "8条"},
	//arr := []int{124, 125, 126, 12, 16, 48, 52, 56, 64, 65, 76, 80, 96, 100}

	//中中7万7万1条2条2条3条4条6条8条8条9条
	//arr := []int{124, 125, 24, 25, 72, 76, 77, 80, 84, 92, 100, 101, 104, 40}

	////{"中", "3万", "3万", "4万", "5万", "4筒", "5筒", "6筒", "8筒", "8筒", "2条", "3条", "7条", "8条"},
	//arr := []int{124, 8, 9, 12, 16, 48, 52, 56, 64, 65, 76, 80, 96, 100}

	//{"6万", "6万", "6万", "2筒", "2筒", "2筒", "3万", "3万", "3万", "4筒", "4筒", "4筒", "8万", "9万"},
	//arr := []int{20, 21, 22, 40, 41, 42, 8, 9, 10, 48, 49, 50, 28, 32}

	//_calHuObject := NewCalHuObject(5)
	//_calHuObject.Check7Dui = true
	////_calHuObject.PeiPaiArr = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104}
	//
	////手牌
	//for _, v := range arr {
	//	_card := NewMCard(v)
	//	_color := _card.GetColor()
	//	//_data := _card.GetData()
	//	_calHuObject.HandPaiArr[_color] = append(_calHuObject.HandPaiArr[_color], _card)
	//}

	//_calHuObject.CalHuPai()

	//_sendTipArr := _calHuObject.CalTingPai()
	//for _, v := range _sendTipArr {
	//	sendCard := NewMCard(v.SendCard)
	//	huCard := make([]*MCard, 0)
	//	for _, vv := range v.HuCards {
	//		huCard = append(huCard, NewMCard(vv))
	//	}
	//	logs.Info("------------------SendCheckTing()-------------->打出%v可胡:%v", sendCard, huCard)
	//}

	ch := &CalHuInfo{
		WanIntArr:  []int{1, 12},
		TongIntArr: []int{45, 58, 71},
		TiaoIntArr: []int{76, 91, 103},
		ZfbIntArr:  []int{108, 112, 118},
		FengIntArr: []int{125, 135, 128},
	}

	handCards := make([]int, 0)
	handCards = append(handCards, ch.WanIntArr...)
	handCards = append(handCards, ch.TongIntArr...)
	handCards = append(handCards, ch.TiaoIntArr...)
	handCards = append(handCards, ch.ZfbIntArr...)
	handCards = append(handCards, ch.FengIntArr...)

	for _, v := range handCards {
		fmt.Print(NewMCard(v), ",")
	}
	fmt.Println()
	fmt.Println(ch.Check_ShiSanLan())
}
