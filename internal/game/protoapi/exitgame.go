package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/table"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

// ExitGame 退出游戏
func ExitGame(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {
	uid := cmd.Head.GetUid()
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

	//正常登陆游戏, 根据uid查询游戏数据库,没有则新建
	logs.Info("玩家退出游戏 _player : %v", _player)

	//检测玩家牌桌
	exitSucc := TableMap.ExitTable(_player)
	if exitSucc {
		playerSvr.Logout(_player.ID())
	}

	pcmd := &protobuf.ResponseCmd{
		Head: &protobuf.RespHead{
			Uid:    proto.Int32(int32(uid)),
			MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
			Result: proto.Int32(0),
		},
	}
	writeResponse(pcmd)

}
