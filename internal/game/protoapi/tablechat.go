package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/table"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/player"
)

//牌桌聊天
func TableChat(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success
	errfunc := func(_result int) {
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
		errfunc(result)
		return
	}

	//判断玩家牌桌
	_table := TableMap.GetTable(_player.GetTableID())
	if _table == nil {
		result = consts.ErrorTableNotExist
		errfunc(result)
		return
	}
	_table.Go(func() {
		if !_table.IsPlayerSit(int(uid)) {
			result = consts.ErrorTableNoPlayer
			errfunc(result)
			return
		}

		_strVale := cmd.GetSimple().GetStrValue()
		if _strVale == "" {
			result = consts.ErrorChatContentIsNull
			errfunc(result)
			return
		}

		//限制发送聊天时间间隔  TODO 敏感词屏蔽

		//发送聊天内容
		_fromSeatID := _table.GetSeats().GetSeatByUID(int(uid)).GetId()
		for _, v := range _table.GetSeats() {
			_sendUid := v.GetPlayer().ID()
			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(_sendUid)),
					MsgID:  proto.Int32(0),
					Result: proto.Int32(int32(result)),
				},
				ChatInfo: &protobuf.ChatInfo{
					SeatID:  proto.Int32(int32(_fromSeatID)),
					Content: proto.String(_strVale),
				},
			}
			writeResponse(pcmd)
		}

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
