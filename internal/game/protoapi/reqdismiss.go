package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/table"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

//请求解散桌子
func ReqDismiss(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success

	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~ReqDismiss 请求解散桌子~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int, _tableId int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~ReqDismiss~~~~~~~~~~~~_tableId:%v,uid:%v,_result:%v", _tableId, uid, _result)
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

	//查找玩家
	_player := playerSvr.GetPlayer(int(uid))
	if _player == nil {
		result = consts.ErrorNotLoginInGame
		errfunc(result, 0)
		return
	}

	//判断玩家牌桌
	_table := TableMap.GetTable(_player.GetTableID())
	if _table == nil {
		result = consts.ErrorTableNotExist
		errfunc(result, 0)
		return
	}
	_table.Go(func() {
		if !_table.IsPlayerSit(int(uid)) {
			result = consts.ErrorTableNoPlayer
			errfunc(result, _table.ID)
			return
		}

		//判断状态,如果已经处于等待解散中,则不能再发送解散请求
		if _table.ReqDismissSeatID != consts.DefaultIndex {
			result = consts.ErrorDismissTableState
			errfunc(result, _table.ID)
			return
		}

		//发送请求解散牌桌请求
		_table.ReqDismissTable(_player)

		//发送请求成功
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(int(uid))),
				MsgID:  proto.Int32(int32(cmd.GetHead().GetMsgID())),
				Result: proto.Int32(int32(result)),
			},
		}
		writeResponse(pcmd)

	})
}
