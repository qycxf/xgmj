package config

import (
	json "encoding/json"
	io "io/ioutil"
	"strconv"

	"strings"

	"qianuuu.com/lib/logs"
	"qianuuu.com/xgmj/internal/mjcomn"
)

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {

	return &JsonStruct{}

}

func (self *JsonStruct) Load(filename string, v interface{}) {

	data, err := io.ReadFile(filename)

	if err != nil {

		return

	}

	datajson := []byte(data)

	err = json.Unmarshal(datajson, v)

	if err != nil {

		return

	}

}

type ValueTestAtmp struct {
	StringValue string

	NumericalValue int

	BoolValue bool
}

type Id struct {
	id   int
	name string
}

type testdata struct {
	UseIndex int
	UsePxArr []string

	TestGameCt string
	TestIdArr  []map[string]interface{}
	TestPxArr  []map[string]interface{} //interface{}

}

var testd = &testdata{
	UseIndex:   0,
	UsePxArr:   make([]string, 0),
	TestGameCt: "",
}

func TestData() *testdata {
	return testd
}

func ReadJson() {

	JsonParse := NewJsonStruct()
	JsonParse.Load("test.json", testd)
	testd.RandPxArr()
	logs.Info("-------------------------------------------------------------->ReadJson test_json success!!")
}

func (td *testdata) PrintIds() {
	for _, v := range td.TestIdArr {
		logs.Info("------------>v:%v", v["id"])
	}
}

func (td *testdata) HasGameCt(_gameCt int) bool {

	gameCtArr := strings.Split(td.TestGameCt, ",")
	for _, v := range gameCtArr {
		if strconv.Itoa(_gameCt) == v {
			return true
		}
	}
	return false
}

func (td *testdata) HasId(_id int) bool {
	logs.Info("------HasId------>td.TestIdArr:%v", td.TestIdArr)
	for _, v := range td.TestIdArr {
		logs.Info("------HasId------>:%v;%v", v["id"], _id)
		if v["id"].(string) == strconv.Itoa(_id) {
			return true
		}
	}
	return false
}

//返回 count 不同牌型
//func (td testdata) GetPxArrs(gameCt int) []string {
//	totalCt := len(td.TestPxArr)
//	arr := make([]int, 0)
//	for {
//		ranIndex := utils.GetRanExceptX(totalCt)
//		if !utils.HasElement(arr, ranIndex) {
//			arr = append(arr, ranIndex)
//			if len(arr) == gameCt {
//				break
//			}
//		}
//	}
//	logs.Info("==================================================>GetPxArrs,arr:%v", arr)
//	strArr := make([]string, 0)
//	for _, v := range arr {
//		str := td.TestPxArr[v]["px"].(string)
//		strArr = append(strArr, str)
//	}
//
//	logs.Info("==================================================>GetPxArrs,strArr:%v len(strArr):%v", strArr, len(strArr))
//	return strArr
//}

func (td *testdata) GetRandomPxArr() string {
	count := len(td.TestPxArr)
	ranIndex := mjcomn.GetRanExceptX(count)
	str := td.TestPxArr[ranIndex]["px"].(string)
	logs.Info("==================================================>ranIndex:%v,str:%v", ranIndex, str)
	return str
}

//将TestPxArr乱序存到 UsePxArr 中
func (td *testdata) RandPxArr() {
	count := len(td.TestPxArr)
	indexArr := make([]int, count)
	for i := 0; i < count; i++ {
		indexArr[i] = i
	}
	indexArr = mjcomn.RandArr(indexArr)
	//logs.Info("-----------RandPxArr------indexArr>%v", indexArr)
	td.UsePxArr = make([]string, 0)
	for i := 0; i < count; i++ {
		_index := indexArr[i]
		_str := td.TestPxArr[_index]["px"].(string)
		td.UsePxArr = append(td.UsePxArr, _str)
	}
	//logs.Info("-----------RandPxArr------td.UsePxArr>%v", td.UsePxArr)
}

func (td *testdata) GetUseStr() string {

	if td.UseIndex < 0 || td.UseIndex >= len(td.UsePxArr) {
		td.UseIndex = 0
	}

	str := td.UsePxArr[td.UseIndex]
	logs.Info("---------->td.GetUseStr:  td.UseIndex:%v,len(td.UsePxArr):%v,str:%v", td.UseIndex, len(td.UsePxArr), str)

	td.UseIndex++
	if td.UseIndex >= len(td.UsePxArr) {
		td.UseIndex = 0
	}
	return str
}
