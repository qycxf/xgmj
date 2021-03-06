//
// Author: leafsoar
// Date: 2016-02-04 22:32:41
//

package values

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Value 值类型
// type Value interface{}

// Value 任何类型
type Value interface{}

// // NewValueFromJSON 从 JSON 解析值
// func NewValueFromJSON(data []byte) (Value, error) {
// 	var value Value
// 	err := json.Unmarshal(data, &value)
// 	return value, err
// }

// NewValueMapArray 返回 ValueＭap 数组
func NewValueMapArray(data []byte) ([]ValueMap, error) {
	var values []interface{}
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}
	var vms []ValueMap
	for _, item := range values {
		vms = append(vms, item.(map[string]interface{}))
	}
	return vms, nil
}

// ValueMap 字典类型
type ValueMap map[string]interface{}

// NewValuesFromJSON 创建新数据
func NewValuesFromJSON(data []byte) (ValueMap, error) {
	var vm ValueMap
	err := json.Unmarshal(data, &vm)
	return vm, err
}

// GetValueMap 返回 ValueMap
func (vm ValueMap) GetValueMap(name string) ValueMap {
	value := vm[name]
	// fmt.Println("vpt", reflect.TypeOf(value))
	switch value.(type) {
	case map[string]interface{}:
		return value.(map[string]interface{})
	case ValueMap:
		return value.(ValueMap)
	}
	return nil
}

// GetString 获取字符串数据
func (vm ValueMap) GetString(name string) string {
	value := vm[name]
	if value == nil {
		return ""
	}
	if name == "answers" {
		fmt.Printf("=======  %T", value)
	}
	switch value.(type) {
	case int:
		return strconv.Itoa(value.(int))
	case string:
		return strings.TrimSpace(value.(string))
	case []interface{}:
		data, _ := json.Marshal(value)
		return string(data)
	case interface{}:
		data, _ := json.Marshal(value)
		return string(data)
	}
	return ""
}

// GetInt 获取 int 类型数据
func (vm ValueMap) GetInt(name string) int {
	value := vm[name]
	if value == nil {
		return 0
	}
	// fmt.Println(reflect.TypeOf(value))
	switch value.(type) {
	case int:
		return value.(int)
	case float64:
		return int(value.(float64))
	case string:
		ret, err := strconv.Atoi(value.(string))
		if err != nil {
			return 0
		}
		return ret
	}
	return 0
}

// GetInt64 获取 int 类型数据
func (vm ValueMap) GetInt64(name string) int64 {
	value := vm[name]
	if value == nil {
		return 0
	}
	// fmt.Println(reflect.TypeOf(value))
	switch value.(type) {
	case int:
		return int64(value.(int))
	case int64:
		return value.(int64)
	case float64:
		return int64(value.(float64))
	case string:
		ret, err := strconv.Atoi(value.(string))
		if err != nil {
			return 0
		}
		return int64(ret)
	}
	return 0
}

// GetFloat64 获取 float64 数据
func (vm ValueMap) GetFloat64(name string) float64 {
	value := vm[name]
	if value == nil {
		return 0
	}
	switch value.(type) {
	case float64:
		return value.(float64)
	case string:
		ret, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return 0
		}
		return ret
	}
	return 0
}

// GetBool 获取 bool 值
func (vm ValueMap) GetBool(name string) bool {
	value := vm[name]
	if value == nil {
		return false
	}
	switch value.(type) {
	case bool:
		return value.(bool)
	case int:
		return value.(int) == 1
	case string:
		return value.(string) == "true"
	}
	return false
}

// ToJSON json 字符串
func (vm ValueMap) ToJSON() []byte {
	ret, _ := json.MarshalIndent(vm, "", "  ")
	return ret
}

// // SetString 设置数据
// func (f Values) SetString(name, value string) {
// 	f[name] = value
// }
