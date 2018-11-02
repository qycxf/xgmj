package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/table"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

//响应解散桌子
func RespDismiss(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success

	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~RespDismiss 玩家响应解散桌子请求~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int, _tableId int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~RespDismiss~~~~~~~~~~~~_tableId:%v,uid:%v,_result:%v", _tableId, uid, _result)
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

		_seat := _table.GetSeats().GetSeatByUID(int(uid))
		_agreeOrNot := cmd.GetSimple().GetIntValue()

		//第一局未开始情况下 不能发送解散房间请求
		if _table.Majhong.GameCt == 0 && _table.GetState() <= consts.TableStateWaiteReady {
			result = consts.ErrorDismissGameNotStart
			errfunc(result, _table.ID)
			return
		}

		if _agreeOrNot == int32(consts.Yes) {
			_table.AgreeDismissArr[_seat.GetId()] = consts.DismissTableStateAgree
			_table.DisMissTable(_seat.GetId())
		} else {

			//有一个人不同意则发送不同意消息
			_table.ClearDismissInfo()

			for _, v := range _table.GetSeats() {
				if !v.GetPlayer().IsRobot() {
					logs.Info("---->_seat.Player.User.NickName:%v v.Player.User.ID:%v", _seat.GetPlayer().NickName(), v.GetPlayer().ID())
					pcmd := &protobuf.ResponseCmd{
						Head: &protobuf.RespHead{
							Uid:    proto.Int32(int32(v.GetPlayer().ID())),
							MsgID:  proto.Int32(0),
							Result: proto.Int32(0),
						},
						Simple: &protobuf.RespSimple{
							Tag:      protobuf.RespSimple_DISMISS_RESULT.Enum(),
							IntValue: proto.Int32(int32(consts.No)),
							StrValue: proto.String("由于玩家[" + _seat.GetPlayer().NickName() + "]拒绝,房间解散失败,游戏继续"),
						},
					}
					writeResponse(pcmd)
				}
			}
		}

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
