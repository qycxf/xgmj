package protoapi

import (
	"qianuuu.com/xgmj/internal/game/table"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

// 断开游戏
func OffLine(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := int(cmd.Head.GetUid())
	logs.Info("OffLine:断开游戏, uid:%v", uid)

}
