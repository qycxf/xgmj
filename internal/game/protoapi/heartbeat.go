package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/player"
)

// 游戏心跳
func HeartBeat(cmd *protobuf.RequestCmd, playerSvr player.Server) {

	uid := int(cmd.Head.GetUid())
	//logs.Info("HeartBeat:游戏心跳, uid:%v", uid)

	result := consts.Success

	errfunc := func(_result int) {
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(uid)),
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
		errfunc(result)
		return
	}

	_player.UpdateMsgRecevTime()
	pcmd := &protobuf.ResponseCmd{
		Head: &protobuf.RespHead{
			Uid:    proto.Int32(int32(uid)),
			MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
			Result: proto.Int32(0),
		},
	}
	writeResponse(pcmd)

}
