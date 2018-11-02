package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/table"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

//离开桌子
func ExitTable(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success
	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~ExitTable 玩家离开桌子请求~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int, _tableId int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~ExitTable~~~~~~~~~~~~_tableId:%v,uid:%v,_result:%v", _tableId, uid, _result)
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

		//检测是否可以离桌
		if !_table.ChkExitTable(_player) {
			result = consts.ErrorExitTable
			errfunc(result, _table.ID)
			return

		}

		if _table.GetState() > consts.TableStateWaiteReady || _table.Majhong.GameCt >= 1 {
			result = consts.ErrorExitTableGameStart
			errfunc(result, _table.ID)
			return
		}

		_table.Exit(_player) //离开桌子

		_table.SendTableInfo() //刷新其它玩家TableInfo

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
