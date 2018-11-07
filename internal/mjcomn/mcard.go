// 麻将牌对象

package mjcomn

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

//var MCARD_DATA_MJ = make([]int, 0)

//合肥麻将84(2-8万 2-8筒 2-8条)
var MCARD_DATA_HFMJ_84 = [84]int{
	4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, //万
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, //筒
	76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, //条
}

//红中麻将112(1-9万 1-9筒 1-9条 4张红中)
var MCARD_DATA_HZMJ_112 = [112]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条
	124, 125, 126, 127, //中
}

//宜昌血流成河108(1-9万 1-9筒 1-9条)
var MCARD_DATA_XLCH_YC_108 = [108]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
}

//卡五星_襄阳 84(1-9筒 1-9条 中发白)
var MCARD_DATA_KA5XING_XY_84 = [84]int{
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
}

//四川血战到底108(1-9万 1-9筒 1-9条)
var MCARD_DATA_XZDD_SC_108 = [108]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
}

//福州麻将  带花
var TableType_FZMJ_FJ_144 = [144]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
	136, 137, 138, 139, 140, 141, 142, 143, //春夏秋冬梅兰竹菊
}

//福州麻将  不带花
var TableType_FZMJ_FJ_108 = [108]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
}

//南京麻将
var TableType_NJMJ_144 = [144]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
	136, 137, 138, 139, 140, 141, 142, 143, //春夏秋冬梅兰竹菊
}

//阜阳麻将
var MCARD_DATA_FYMJ_136 = [136]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
}

//鞍山麻将
var MCARD_DATA_ASMJ_136 = [136]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
}

//安庆麻将
var TableType_AQMj_136 = [136]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
}

//蚌埠麻将
var TableType_BBMj_136 = [136]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
}

//怀远麻将
var TableType_HYMJ_144 = [144]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
	136, 137, 138, 139, 140, 141, 142, 143, //春夏秋冬梅兰竹菊
}

//国标麻将
var TableType_GBMJ_144 = [144]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
	136, 137, 138, 139, 140, 141, 142, 143, //春夏秋冬梅兰竹菊
}

//怀远麻将
var TableType_XGMJ_144 = [144]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, //万 0 - 35
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, //筒 36 - 71
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, //条 72 - 107
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, //东南西北
	124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, //中发白
	136, 137, 138, 139, 140, 141, 142, 143, //春夏秋冬梅兰竹菊
}

// MCard 牌对象
type MCard struct {
	data  int
	color int
	value int
}

// 创建一张扑克牌
func NewMCard(_data int) *MCard {
	return &MCard{
		data:  _data,
		color: GetColor(_data),
		value: GetVal(_data),
	}
}

//复制一张牌
func (_mcard *MCard) Clone() *MCard {
	mcard := NewMCard(_mcard.data)
	return mcard
}

//是否是花牌
func (_mcard *MCard) IsHuaPai() bool {
	_data := _mcard.data
	if _data >= 136 {
		return true
	}
	return false
}

//是否南京花牌
func (_mcard *MCard) IsNJHuaPai() bool {
	_data := _mcard.data
	if _data >= 124 {
		return true
	}
	return false
}

//是安庆花牌
func (_mcard *MCard) IsAQHuaPai() bool {
	_data := _mcard.data
	if _data >= 124 {
		return true
	} else if _data >= 108 && _data <= 111 {
		return true
	}
	return false
}

//是蚌埠花牌
func (_mcard *MCard) IsBBHuaPai() bool {
	_data := _mcard.data
	if _data >= 124 {
		return true
	} else if _data >= 108 && _data <= 111 {
		return true
	}
	return false
}

//是否是花牌
func (_mcard *MCard) IsXGHuaPai() bool {
	_data := _mcard.data
	if _data >= 136 {
		return true
	}
	return false
}

//是怀远花牌
func (_mcard *MCard) IsHYHuaPai(fengLing int) bool {
	_data := _mcard.data
	if _data >= 124 {
		return true
	} else if _data >= 108 && _data <= 111 && fengLing > 0 {
		return true
	}
	return false
}

//获取值
func GetVal(_data int) int {
	if _data < 108 {
		return _data%36/4 + 1
	} else if _data < 124 {
		return (_data-108)/4 + 1
	} else if _data < 136 {
		return (_data-124)/4 + 1
	}
	return (_data-136)/4 + 1
}

//获取花色
func GetColor(_data int) int {
	if _data < 108 {
		return _data / 36
	} else if _data < 124 {
		return Color_Feng
	} else if _data < 136 {
		return Color_Zfb
	}
	return Color_Hua
}

//返回数据
func (card *MCard) GetData() int {
	return card.data
}

//返回花色
func (card *MCard) GetColor() int {
	return card.color
}

func (card *MCard) SetColor(_color int) {
	card.color = _color
}

//返回牌值
func (card *MCard) GetValue() int {
	return card.value
}

func (card *MCard) SetValue(_value int) {
	card.value = _value
}

//获取牌花色名称
func (card *MCard) getColorName() string {

	switch card.color {
	case Color_Wan:
		return "万"
	case Color_Tong:
		return "筒"
	case Color_Tiao:
		return "条"
	default:
		return ""
	}
}

//获取牌值
func (card *MCard) getValueName() string {
	return strconv.Itoa(card.value)
}

//花色\牌值 相等
func (card *MCard) Equal(c2 *MCard) bool {
	return card.EqualColor(c2) && card.EqualValue(c2)
}

func (card *MCard) EqualByData(_data int) bool {
	return card.color == GetColor(_data) && card.value == GetVal(_data)
}

//牌值相等
func (card *MCard) EqualValue(c2 *MCard) bool {
	return card.value == c2.value
}

//花色相等
func (card *MCard) EqualColor(c2 *MCard) bool {
	return card.color == c2.color
}

//同一张牌
func (card *MCard) Same(c2 *MCard) bool {
	return card.data == c2.data
}

//获取牌名称
func (card *MCard) String() string {
	if card.data < 108 {
		return card.getValueName() + card.getColorName()
	} else if card.data < 112 {
		return "东"
	} else if card.data < 116 {
		return "南"
	} else if card.data < 120 {
		return "西"
	} else if card.data < 124 {
		return "北"
	} else if card.data < 128 {
		return "中"
	} else if card.data < 132 {
		return "发"
	} else if card.data < 136 {
		return "白"
	} else if card.data < 137 {
		return "春"
	} else if card.data < 138 {
		return "夏"
	} else if card.data < 139 {
		return "秋"
	} else if card.data < 140 {
		return "冬"
	} else if card.data < 141 {
		return "梅"
	} else if card.data < 142 {
		return "兰"
	} else if card.data < 143 {
		return "竹"
	} else if card.data < 144 {
		return "菊"
	}
	return "CARD_ERROR!"
}

//获取牌详细信息
func (card *MCard) Detail() string {
	return card.String() + "[" + strconv.Itoa(card.data) + "]"
}

// Cards 牌组对象 ------------------------------------------------------------------------------
type MCards []*MCard

//  创建牌组 _cardCt: 使用的牌组
func NewMCards(_tableType int) MCards {

	MCARD_DATA_MJ := make([]int, 0)

	if _tableType == TableType_HFMJ { // 合肥麻将
		for _, value := range MCARD_DATA_HFMJ_84 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}

	} else if _tableType == TableType_HZMJ { // 红中麻将
		for _, value := range MCARD_DATA_HZMJ_112 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}

	} else if _tableType == TableType_YC_XLCH { //宜昌血流成河
		for _, value := range MCARD_DATA_XLCH_YC_108 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_XY_KA5XING { //襄阳卡五星
		for _, value := range MCARD_DATA_KA5XING_XY_84 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}

	} else if _tableType == TableType_SC_XZDD { //四川血战到底
		for _, value := range MCARD_DATA_XZDD_SC_108 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_FJ_FZMJ { //福建福州麻将
		for _, value := range TableType_FZMJ_FJ_144 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_NJMJ { //南京麻将
		for _, value := range TableType_NJMJ_144 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_FYMJ { //阜阳麻将
		for _, value := range MCARD_DATA_FYMJ_136 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_ASMJ { //鞍山麻将
		for _, value := range MCARD_DATA_ASMJ_136 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_GBMJ { //国标麻将
		for _, value := range TableType_GBMJ_144 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_AQMJ { //安庆麻将
		for _, value := range TableType_AQMj_136 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_HYMJ { //怀远麻将
		for _, value := range TableType_HYMJ_144 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	} else if _tableType == TableType_XGMJ { //香港麻将
		for _, value := range TableType_XGMJ_144 {
			MCARD_DATA_MJ = append(MCARD_DATA_MJ, value)
		}
	}

	cards := make([]*MCard, 0)
	for _, value := range MCARD_DATA_MJ {
		cards = append(cards, NewMCard(value))
	}
	return cards
}

//根据名称获取元素下标,并记录下标
func (cs MCards) GetIndexByName(cardName string) int {
	for k, v := range cs {
		if v.String() == cardName {
			return k
		}
	}
	return -1
}

// Len 长度
func (ds MCards) Len() int {
	return len(ds)
}

// Swap 交换
func (ds MCards) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

// Less 排序接口
func (ds MCards) Less(i, j int) bool {
	return ds[i].data < ds[j].data
}

// Sort 排序
func (ds MCards) Sort() {
	sort.Sort(ds)
}

//牌组信息
func (ds MCards) String() string {
	values := make([]string, 0, len(ds))
	for _, card := range ds {
		values = append(values, card.String())
	}
	return strings.Join(values, ",")
}

//乱序排列
func (ds MCards) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, length := 0, len(ds); i < length; i++ {
		rand := r.Intn(length - i)
		ds[rand], ds[i] = ds[i], ds[rand]
	}
}
