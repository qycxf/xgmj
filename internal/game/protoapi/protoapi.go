//
// Author: leafsoar
// Date: 2016-05-21 14:51:05
//

// proto 简单接口实现接口

package protoapi

import (
	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/lib/values"
	"qianuuu.com/player"
)

// 响应消息函数
var writeResponse func(*protobuf.ResponseCmd)

// Init 初始化
func Init(writeFunc func(*protobuf.ResponseCmd)) {
	writeResponse = writeFunc
}

func requestTableID(playerSvr player.Server, createuid, itype, iname, charge int) (int, error) {
	vm := values.ValueMap{
		"appname":     config.Opts().AppName,
		"app_name":    config.Opts().AppName,
		"inning_type": itype,
		"inning_name": iname,
		"creator_id":  createuid,
		"charge":      charge,
	}
	return playerSvr.GetTableID(vm)
}
