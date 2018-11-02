package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/table"
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

//进入房间桌子
func EnterTable(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success
	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~EnterTable 玩家进入桌子请求~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_tableId int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~EnterTable~~~~~~~~~~~~_tableId:%v,uid:%v,_result:%v", _tableId, uid, result)
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
		errfunc(0)
		return
	}

	//检测是否已经在桌子中
	if _player.GetTableID() > 0 {
		result = consts.ErrorHasInTable
		errfunc(_player.GetTableID())
		return
	}

	//验证桌子编号
	_tableID := cmd.GetSimple().GetIntValue()
	_table := TableMap.GetTable(int(_tableID))
	if _table == nil {
		//弹窗提示桌子不存在
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(int(uid))),
				MsgID:  proto.Int32(int32(0)),
				Result: proto.Int32(int32(result)),
			},
			Simple: &protobuf.RespSimple{
				Tag:      protobuf.RespSimple_WIN_TIP_M.Enum(),
				StrValue: proto.String("房间不存在,请重新输入!"),
			},
		}
		writeResponse(pcmd)

		//发送 EnterTable失败
		pcmd = &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(int(uid))),
				MsgID:  proto.Int32(int32(cmd.GetHead().GetMsgID())),
				Result: proto.Int32(int32(result)),
			},
			Simple: &protobuf.RespSimple{
				Tag: protobuf.RespSimple_ENTER_TABLE.Enum(),
			},
		}
		writeResponse(pcmd)

		return
	}

	_table.Go(func() {

		_idleSeatID := _table.SearchIdleSeat()
		if _idleSeatID == -1 {
			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(int(uid))),
					MsgID:  proto.Int32(int32(0)),
					Result: proto.Int32(int32(result)),
				},
				Simple: &protobuf.RespSimple{
					Tag:      protobuf.RespSimple_WIN_TIP_M.Enum(),
					StrValue: proto.String("房间人已满，游戏已经开始！"),
				},
			}
			writeResponse(pcmd)

			//发送 EnterTable失败
			pcmd = &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(int(uid))),
					MsgID:  proto.Int32(int32(cmd.GetHead().GetMsgID())),
					Result: proto.Int32(int32(result)),
				},
				Simple: &protobuf.RespSimple{
					Tag: protobuf.RespSimple_ENTER_TABLE.Enum(),
				},
			}
			writeResponse(pcmd)
			return
		}

		//收费\检测玩家房卡
		actNeedCt:=0
		if _table.TableCfg.PayWay==0 {
			actNeedCt=_table.TableCfg.AvgPerCard//均摊房卡
		}else {
			actNeedCt=_table.TableCfg.AvgPerCard*_table.TableCfg.PlayerCt//房主付费或者大赢家付费
		}

		if actNeedCt > 0 {
			if config.Opts().OpenCharge {
				if playerSvr.GetUsableFangka(int(uid)) < actNeedCt {
					//提示房卡不足
					pcmd := &protobuf.ResponseCmd{
						Head: &protobuf.RespHead{
							Uid:    proto.Int32(int32(_player.ID())),
							MsgID:  proto.Int32(0),
							Result: proto.Int32(0),
						},
						Simple: &protobuf.RespSimple{
							Tag:      protobuf.RespSimple_WIN_TIP_M.Enum(),
							StrValue: proto.String("房卡不足,请充值!"),
						},
					}
					writeResponse(pcmd)

					//发送创建桌子失败
					pcmd = &protobuf.ResponseCmd{
						Head: &protobuf.RespHead{
							Uid:    proto.Int32(int32(uid)),
							MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
							Result: proto.Int32(int32(result)),
						},
					}
					writeResponse(pcmd)
					return
				}
			}
		}

		//发送 EnterTable
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(int(uid))),
				MsgID:  proto.Int32(int32(cmd.GetHead().GetMsgID())),
				Result: proto.Int32(int32(result)),
			},
			Simple: &protobuf.RespSimple{
				Tag: protobuf.RespSimple_ENTER_TABLE.Enum(),
			},
		}
		writeResponse(pcmd)

		//进入桌子坐下
		_table.JoinTable(_player)

	})

}
