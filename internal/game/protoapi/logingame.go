package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/config"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/game/table"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

// LoginGame 登陆游戏
func LoginGame(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := int(cmd.Head.GetUid())
	logs.Debug("[LoginGame] 登陆游戏, uid:%v", uid)
	defer func() {
		logs.Debug("[LoginGame] 登陆游戏, uid:%v --end", uid)
	}()

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

	//判断是否需要断线重连,仅一种情况重连:在牌桌中且在游戏中
	//isReConnect := false //是否断线重连
	_player := playerSvr.GetPlayer(uid)
	if _player != nil {
		if _player.GetTableID() > 0 {
			_table := TableMap.GetTable(_player.GetTableID())
			if _table != nil {
				//处理断线重连
				if config.Opts().OpenReConnect {
					logs.Info("玩家断线重连,TableId:%v", _player.GetTableID())

					//发送进入桌子成功
					pcmd := &protobuf.ResponseCmd{
						Head: &protobuf.RespHead{
							Uid:    proto.Int32(int32(uid)),
							MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
							Result: proto.Int32(0),
						},
						Simple: &protobuf.RespSimple{
							Tag: protobuf.RespSimple_ENTER_TABLE.Enum(),
						},
					}
					writeResponse(pcmd)
					logs.Info("--------------------------------------->RespSimple_ENTER_TABLE")

					//发送TableInfo
					pcmd = &protobuf.ResponseCmd{
						Head: &protobuf.RespHead{
							Uid:    proto.Int32(int32(uid)),
							MsgID:  proto.Int32(0),
							Result: proto.Int32(0),
						},
						PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
						TableInfo:  protobuf.Helper.GetTableInfo(_table, uid),
					}
					writeResponse(pcmd)
					logs.Info("--------------------------------------->TableInfo")

					//如果当前牌桌有请求解散状态,则发送给该玩家
					if _table.ReqDismissSeatID != consts.DefaultIndex {
						//申请人位置,及各个位置当前是否同意
						intArr := make([]int32, 5)
						intArr[4] = int32(_table.ReqDismissSeatID)
						for _, v := range _table.GetSeats() {
							intArr[v.GetId()] = int32(_table.AgreeDismissArr[v.GetId()])
						}
						logs.Info("-----_tableId:%v----牌桌处于等待解散状态中------>dismissRemainTime:%v,_table.GetSeats():%v", _table.ID, _table.GetDismissRemainTime(), _table.GetSeats())
						pcmd := &protobuf.ResponseCmd{
							Head: &protobuf.RespHead{
								Uid:    proto.Int32(int32(_player.ID())),
								MsgID:  proto.Int32(0),
								Result: proto.Int32(0),
							},
							Simple: &protobuf.RespSimple{
								Tag:      protobuf.RespSimple_DISMISS_TABLE.Enum(),
								IntValue: proto.Int32(int32(_table.GetDismissRemainTime())),
								StrArr:   _table.GetDismissInfo(),
								IntArr:   intArr,
							},
						}
						writeResponse(pcmd)
						logs.Info("--------------------------------------->RespSimple_DISMISS_TABLE")
					}

					return
				} else {
					//未开启断线重连,直接解散房间,发送解散房间信息
					for _, v := range _table.GetSeats() {
						_table.Exit(v.GetPlayer())
						_sendUid := v.GetPlayer().ID()
						if !v.GetPlayer().IsRobot() {
							logs.Info("----------------->>>>>>> _sendUid::::%v ", _sendUid)
							pcmd := &protobuf.ResponseCmd{
								Head: &protobuf.RespHead{
									Uid:    proto.Int32(int32(_sendUid)),
									MsgID:  proto.Int32(0),
									Result: proto.Int32(0),
								},
								Simple: &protobuf.RespSimple{
									Tag:      protobuf.RespSimple_WIN_TIP_S.Enum(),
									StrValue: proto.String("由于" + _player.String() + "断线! 系统强制解散房间~~"),
								},
							}
							writeResponse(pcmd)
						}
					}
				}
			} else {
				logs.Info("********************* LoginGame error _player.TableId:%v,_player:%v", _player.GetTableID(), _player.String())
			}
		}

		//内存中有,保存一次数据
		//playerSvr.SavePlayer(_player)
	}

	_player, _err := playerSvr.Login(uid) //玩家正常登陆
	if _err != nil {
		result = consts.ErrorPlayerNotExist
		errfunc(result)
		return
	}
	logs.Info("玩家登陆游戏 err:%v", _err)

	pcmd := &protobuf.ResponseCmd{
		Head: &protobuf.RespHead{
			Uid:    proto.Int32(int32(uid)),
			MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
			Result: proto.Int32(0),
		},
		PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
	}
	writeResponse(pcmd)
	logs.Info("uid: %d 玩家登陆游戏 --------------------------> over ", uid)
}
