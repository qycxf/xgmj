package mjcomn

import (
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"
)

//升序排序数组
func SortIntArrAsc(array []int) []int {
	for i := 0; i < len(array); i++ {
		for j := i + 1; j < len(array); j++ {
			if array[j] < array[i] {
				// tmp := array[i]
				// array[i] = array[j]
				// array[j] = tmp
				array[i], array[j] = array[j], array[i]
			}
		}
	}
	return array
}

//判断两个数组是否完全相同
func EqualArr(_arr1 []int, _arr2 []int) bool {

	arr1 := SortIntArrAsc(_arr1)
	arr2 := SortIntArrAsc(_arr2)

	if len(arr1) != len(arr2) {
		return false
	}
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

//数组删除指定元素,相同的全部删除
func RemoveElement(_array []int, _value int) []int {
	tmpArr := make([]int, 0)
	for _, v := range _array {
		if v != _value {
			tmpArr = append(tmpArr, v)
		}
	}
	return tmpArr
}

func RemoveStrElement(_array []string, _value string) []string {
	tmpArr := make([]string, 0)
	isRemove := false
	for _, v := range _array {
		if v != _value {
			tmpArr = append(tmpArr, v)
		} else {
			if !isRemove { //一次只删除一个
				isRemove = true
			} else {
				tmpArr = append(tmpArr, v)
			}
		}
	}
	return tmpArr
}

//乱序一个数组
func RandArr(_array []int) []int {
	tmpArr := make([]int, 0)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(_array)
	for i := 0; i < length; i++ {
		ranIndex := random.Intn(len(_array))
		val := _array[ranIndex]
		tmpArr = append(tmpArr, val)
		_array = RemoveElement(_array, val)
	}
	return tmpArr
}

func RandStrArr(_array []string) []string {
	tmpArr := make([]string, 0)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(_array)
	for i := 0; i < length; i++ {
		ranIndex := random.Intn(len(_array))
		val := _array[ranIndex]
		tmpArr = append(tmpArr, val)
		_array = RemoveStrElement(_array, val)
	}
	return tmpArr
}

//获取一个 count位 的随机数
func GetRandomNum(count int) int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	number := 0
	for i := 0; i < count; i++ {
		ran := random.Intn(9) //0~9
		if i == count-1 {
			ran = random.Intn(8) + 1 //第一位不为0
		}

		value := ran * Power(10, i)
		//logs.Info("value:%v", value)
		number += value
	}

	return number
}

//获取一个 [0,count) 的随机数
func GetRanExceptX(count int) int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	ran := random.Intn(count) //0~9
	return ran
}

//求value的 ct 次冥
func Power(value int, ct int) int {
	if ct == 0 {
		return 1
	}
	_tmp := value //保存底数
	for i := 0; i < ct-1; i++ {
		value *= _tmp
	}
	//log.Info("ct:%v------>value:%v", ct, value)
	return value
}

//获取数组中的最大元素
func GetMaxElement(_array []int) int {
	_data := _array[0]
	for _, v := range _array {
		if v > _data {
			_data = v
		}
	}
	return _data
}

//判断数组中是否含有指定元素
func HasElement(_array []int, _element int) bool {
	for _, v := range _array {
		if v == _element {
			return true
		}
	}
	return false
}

//判断数组中是否含有指定元素
func HasElementCt(_array []int, _element int) int {
	ct := 0
	for _, v := range _array {
		if v == _element {
			ct++
		}
	}
	return ct
}

//判断数组中data==>>val
func DataToVal(_array []int) []int {
	_valArr := make([]int, 0)
	for i := 0; i < len(_array); i++ {
		_valArr = append(_valArr, GetVal(_array[i]))
	}
	return _valArr
}

//数组中去重
func DeletePartArr(_array []int, _partArr []int) []int {
	for i := 0; i < len(_partArr); i++ {
		tmpArr := make([]int, 0)
		flag := false
		for _, v := range _array {
			if v == _partArr[i] && !flag {
				flag = true
			} else {
				tmpArr = append(tmpArr, v)
			}
		}
		_array = tmpArr
	}
	return _array
}

func main() {
	//array := []int{3, 4, 5, 1, 2, 0}
	//array = SortIntArrAsc(array)
	//fmt.Print(array)
	number := GetRandomNum(3)
	log.Info("number:%v", number)
}
