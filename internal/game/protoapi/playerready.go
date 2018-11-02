package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/table"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

//玩家准备
func PlayerReady(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success
	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~PlayerReady 玩家准备请求~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int, _tableId int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~PlayerReady~~~~~~~~~~~~_tableId:%v,uid:%v,_result:%v", _tableId, uid, _result)
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

		// 判断牌桌状态
		_tableState := _table.GetState()
		if _tableState != consts.TableStateIdle &&
			_tableState != consts.TableStateWaiteReady {
			result = consts.ErrorTableStateReady
			errfunc(result, _table.ID)
			return
		}

		//判断玩家准备状态
		_seat := _table.GetSeats().GetSeatByUID(int(uid))
		if _seat.GetState() >= consts.SeatStateGameReady {
			result = consts.ErrorPlayerHasReady
			errfunc(result, _table.ID)
			return
		}

		//牌桌准备
		_table.SeatReady(int(uid))

		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(uid)),
				MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
				Result: proto.Int32(int32(result)),
			},
		}
		writeResponse(pcmd)

	})

}
