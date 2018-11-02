package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

// 请求玩家信息
func GetPlayerInfo(cmd *protobuf.RequestCmd, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success
	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~GetPlayerInfo 请求玩家信息~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~GetPlayerInfo~~~~~~~~~~~,uid:%v,_result:%v", uid, _result)
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

	//刷新数据
	_player := playerSvr.GetPlayer(int(uid))
	if _player == nil {
		result = consts.ErrorNotLoginInGame
		errfunc(result)
		return
	}

	pcmd := &protobuf.ResponseCmd{
		Head: &protobuf.RespHead{
			Uid:    proto.Int32(uid),
			MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
			Result: proto.Int32(int32(result)),
		},
		PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
	}
	writeResponse(pcmd)
}
