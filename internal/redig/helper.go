//
// Author: leafsoar
// Date: 2016-05-11 16:12:36
//

package redig

import (
	"fmt"

	"qianuuu.com/ahmj/internal/config"
	"qianuuu.com/ahmj/internal/protobuf"

	"golang.org/x/protobuf/proto"
)

// func init() {
// 	// 获取响应消息的属性 id
// 	cmd := protobuf.ResponseCmd{}
// 	cmdtype := reflect.TypeOf(cmd)
// 	length := cmdtype.NumField()
// 	for i := 0; i < length; i++ {
// 		field := cmdtype.Field(i)
// 		fp := field.Tag.Get("protobuf")
// 		if fp != "" {
// 			id := strings.Split(fp, ",")[1]
// 			fmt.Println(field, fp, id)
// 		}
// 	}
// }

// SetRespIDs 注意与 proto 定义保持一致
func SetRespIDs(cmd *protobuf.ResponseCmd) {
	// 根据消息内容设置对应的 ids
	var ids int32
	if cmd.Head != nil {
		ids = ids | 1<<1
	}
	if cmd.Simple != nil {
		ids = ids | 1<<2
	}
	if cmd.PlayerInfo != nil {
		ids = ids | 1<<3
	}
	if cmd.TableInfo != nil {
		ids = ids | 1<<4
	}
	if cmd.ChatInfo != nil {
		ids = ids | 1<<5
	}
	if cmd.WordChat != nil {
		ids = ids | 1<<6
	}

	// fmt.Println("ids:", ids)
	cfg := config.Opts()
	cmd.Head.App = proto.String(cfg.AppName)
	cmd.Head.RespIDs = proto.Int32(ids)
}

// ReqCmdInfo 请求消息说明
func ReqCmdInfo(cmd *protobuf.RequestCmd) interface{} {
	return cmd.GetHead()
}

// RespCmdInfo 响应消息说明
func RespCmdInfo(cmd *protobuf.ResponseCmd) interface{} {
	var msgs []string
	if cmd.Head != nil {
		msgs = append(msgs, "head")
	}
	if cmd.Simple != nil {
		msgs = append(msgs, "simple")
	}
	if cmd.PlayerInfo != nil {
		msgs = append(msgs, "playerinfo")
	}
	if cmd.TableInfo != nil {
		msgs = append(msgs, "tableinfo")
	}
	if cmd.ChatInfo != nil {
		msgs = append(msgs, "chatinfo")
	}
	if cmd.WordChat != nil {
		msgs = append(msgs, "wrodchat")
	}

	head := cmd.GetHead()
	return fmt.Sprintf("uid: %v msgid: %v msg: %v", head.GetUid(), head.GetMsgID(), msgs)
}

func simpleReq(cmd *protobuf.RequestCmd) string {
	if cmd.GetSimple() != nil {
		return fmt.Sprintf("%v %v", cmd.GetHead(), cmd.GetSimple())
	} else if cmd.GetCrateTable() != nil {
		// return fmt.Sprintf("%v %v", cmd.GetHead(), cmd.GetCrateTable().GetPlayType())
	}
	return fmt.Sprintf("%v", cmd.GetHead())
}
