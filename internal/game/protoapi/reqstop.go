package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/table"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

// 程序后台运行
func ReqStop(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success

	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~ReqStop 程序后台运行请求~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int, _tableId int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~ReqStop~~~~~~~~~~~~_tableId:%v,uid:%v,_result:%v", _tableId, uid, _result)
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(uid),
				MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
				Result: proto.Int32(int32(result)),
				Tip:    proto.String(consts.GetErrTip(result)),
			},
		}
		writeResponse(pcmd)
	}

	_player := playerSvr.GetPlayer(int(uid))
	if _player == nil {
		result = consts.ErrorNotLoginInGame
		errfunc(result, 0)
		return
	}

	if !_player.IsOffline() {
		_player.SetOffline(true)
		if _player.GetTableID() > 0 {
			_table := TableMap.GetTable(_player.GetTableID())
			_table.Go(func() {
				_table.SendTableInfo()
			})
		}
	}
}
