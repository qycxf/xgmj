package protoapi

import (
	"strconv"

	"golang.org/x/protobuf/proto"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/table"
	"qianuuu.com/xgmj/internal/mjcomn"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

//玩家操作
func PlayerOpt(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success
	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~PlayerOpt 玩家操作请求~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int, _tableId int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~PlayerOpt~~~~~~~~~~~~_tableId:%v,uid:%v,_result:%v", _tableId, uid, _result)
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

		//验证值的合法性
		_intVale := int(cmd.GetSimple().GetIntValue()) //操作值
		//if _intVale < mjcomn.OptTypePeng || _intVale > mjcomn.OptTypeCancel {
		if _intVale < mjcomn.OptTypePeng || _intVale > mjcomn.OptTypeTing {
			result = consts.ErrorTableOptVal
			errfunc(result, _table.ID)
			return
		}
		_strVale := cmd.GetSimple().GetStrValue() //出牌\杠牌 值
		_cardData, _ := strconv.Atoi(_strVale)

		_valiNum := cmd.GetSimple().GetValiNum()
		logs.Info("tableId%v~~~~~~~~~~~~PlayerOpt 验证随机数~~~~~~~~~~~~uid:%v-----_valiNum:%v", _table.ID, uid, _valiNum)

		//验证操作状态
		_seat := _table.GetSeats().GetSeatByUID(int(uid))
		result = _table.CheckOpt(_seat, _intVale, _cardData)
		if result != consts.Success {
			errfunc(result, _table.ID)
			return
		}

		//执行位置操作
		_table.SeatOpt(_seat, _intVale, _cardData, false)

		//发送位置操作信息
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(uid)),
				MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
				Result: proto.Int32(int32(0)),
			},
		}
		writeResponse(pcmd)
	})

}
