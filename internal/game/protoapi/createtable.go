package protoapi

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
	"qianuuu.com/xgmj/internal/config"
	"qianuuu.com/xgmj/internal/consts"
	"qianuuu.com/xgmj/internal/game/table"
	"qianuuu.com/xgmj/internal/mjcomn"
	"qianuuu.com/xgmj/internal/protobuf"
)

//玩家创建桌子
func CreateTable(cmd *protobuf.RequestCmd, TableMap *table.Tables, playerSvr player.Server) {

	uid := cmd.Head.GetUid()
	result := consts.Success
	logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~CreateTable 玩家创建桌子请求~~~~~~~~~~~~uid:%v", uid)

	errfunc := func(_result int) {
		logs.Custom(logs.NetworkTag, "~~~~~~~~~~~~CreateTable~~~~~~~~~~~,uid:%v,_result:%v", uid, _result)
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

	//验证创建房间参数
	_tableType := cmd.GetCrateTable().GetTableType()                // 桌子类型
	_playerCt := int32(cmd.GetCrateTable().GetPlayerCt())           // 玩家数
	_gameCt := int32(cmd.GetCrateTable().GetGameCt())               // 游戏局数
	_baseScore := int32(cmd.GetCrateTable().GetBaseScore())         // 牌局底分
	_robotCt := int32(cmd.GetCrateTable().GetRobotCt())             // 测试机器人数
	_dianpaoHu := int32(cmd.GetCrateTable().GetDianpaoHu())         // 点炮胡
	_zimoHu := int32(cmd.GetCrateTable().GetZimoHu())               // 自摸胡
	_tiandiHu := int32(cmd.GetCrateTable().GetTiandiHu())           // 天地胡
	_kehuQidui := int32(cmd.GetCrateTable().GetKehuQidui())         // 可胡七对
	_qiangGang := int32(cmd.GetCrateTable().GetQiangGang())         // 可抢杠
	_zhuaNiaoCt := int32(cmd.GetCrateTable().GetZhuaNiaoCt())       // 抓鸟数
	_yiMaQuanZh := int32(cmd.GetCrateTable().GetYiMaQuanZh())       // 一码全中
	_menQing := int32(cmd.GetCrateTable().GetMenQing())             // 门清
	_present := int32(cmd.GetCrateTable().GetPresent())             // 赠送
	_tdqZuiCt := int32(cmd.GetCrateTable().GetTdqZuiCt())           // 合肥麻将,天地胡\清一色嘴数
	_kePengGang := int32(cmd.GetCrateTable().GetKePengGang())       // 红中麻将,可碰杠
	_kaiHuSuanGang := int32(cmd.GetCrateTable().GetKaiHuSuanGang()) // 阜阳麻将,开胡算杠
	_youGangYouFen := int32(cmd.GetCrateTable().GetYouGangYouFen()) // 阜阳麻将,有杠有分(没荒庄)
	_daiHua := int32(cmd.GetCrateTable().GetDaiHua())               // 蚌埠麻将,是否带花
	_fenging := int32(cmd.GetCrateTable().GetFengLing())            // 怀远麻将,是否带风令
	_baoTing := int32(cmd.GetCrateTable().GetBaoTing())             // 怀远麻将,是否报听
	_wuHuaGuo := int32(cmd.GetCrateTable().GetWuHuaGuo())           // 怀远麻将,是否有无花果选项
	_PayWay := int32(cmd.GetCrateTable().GetPayWay())               //付费方式

	_tableConfig := config.NewTableCfg()
	_tableConfig.TableType = int(_tableType)
	_tableConfig.PlayerCt = int(_playerCt)
	_tableConfig.GameCt = int(_gameCt)
	_tableConfig.BaseScore = int(_baseScore)
	_tableConfig.RobotCt = int(_robotCt)
	_tableConfig.DianpaoHu = int(_dianpaoHu)
	_tableConfig.ZimoHu = int(_zimoHu)
	_tableConfig.TiandiHu = int(_tiandiHu)
	_tableConfig.KehuQidui = int(_kehuQidui)
	_tableConfig.QiangGang = int(_qiangGang)
	_tableConfig.ZhuaNiaoCt = int(_zhuaNiaoCt)
	_tableConfig.YiMaQuanZh = int(_yiMaQuanZh)
	_tableConfig.MenQing = int(_menQing)
	_tableConfig.Present = int(_present)
	_tableConfig.TdqZuiCt = int(_tdqZuiCt)
	_tableConfig.KePengGang = int(_kePengGang)
	_tableConfig.KaiHuSuanGang = int(_kaiHuSuanGang)
	_tableConfig.YouGangYouFen = int(_youGangYouFen)
	_tableConfig.CreaterId = _player.ID() //保存创建者uid
	_tableConfig.PayWay = int(_PayWay)
	_tableConfig.DaiHua = int(_daiHua)
	_tableConfig.FengLing = int(_fenging)
	_tableConfig.BaoTing = int(_baoTing)
	_tableConfig.WuHuaGuo = int(_wuHuaGuo)

	//创建牌桌参数验证
	//if _tableType != mjcomn.TableType_BBMJ &&
	//	_tableType != mjcomn.TableType_HYMJ {
	//	result = consts.ErrorCreateTableType
	//	errfunc(result)
	//	return
	//}

	//点炮胡\自摸胡 有且仅选一个
	//if _tableType == mjcomn.TableType_HYMJ {
	//	if _dianpaoHu == int32(consts.Yes) && _zimoHu == int32(consts.Yes) ||
	//		_dianpaoHu == int32(consts.No) && _zimoHu == int32(consts.No) {
	//		result = consts.ErrorCreateTableSelectParam
	//		errfunc(result)
	//		return
	//	}
	//}

	//验证合肥麻将
	//if _tableType == mjcomn.TableType_BBMJ {
	//_tableConfig.DianpaoHu = consts.Yes
	//底分选择
	//selectScoreCtArr := config.Opts().SelectScoreCt
	//isFind := false
	//for _, v := range selectScoreCtArr {
	//	if _baseScore == v {
	//		isFind = true
	//		break
	//	}
	//}
	//if !isFind {
	//	result = consts.ErrorCreateTableFanCt
	//	errfunc(result)
	//	return
	//}

	//天地胡\清一色嘴数
	//selectTdqZuiCt := config.Opts().SelectTdqZuiCt
	//isFind = false
	//for _, v := range selectTdqZuiCt {
	//	if _tdqZuiCt == v {
	//		isFind = true
	//		break
	//	}
	//}
	//if !isFind {
	//	result = consts.ErrorCreateTableTDQZuiCt
	//	errfunc(result)
	//	return
	//}

	//}

	//检测创建房间费用(房卡)
	//roomCards := config.Opts().RoomCards
	//valiGameCt := false
	//for _, v := range roomCards {
	//	if _gameCt == v.GameCt {
	//		_tableConfig.FangKa = int(v.NeedCard)
	//		_tableConfig.AvgPerCard = int(v.AvgPerCard)
	//		valiGameCt = true
	//		break
	//	}
	//}

	//if !valiGameCt {
	//	result = consts.ErrorCreateTableGameCt
	//	errfunc(result)
	//	return
	//}
	//检测创建房间费用(房卡)
	if _tableType == mjcomn.TableType_HYMJ {

		//roomCards := config.Opts().RoomCards
		//for _, v := range roomCards {
		//	if _gameCt == v.GameCt {
		//		_tableConfig.FangKa = int(v.NeedCard)
		//		_tableConfig.AvgPerCard = int(v.AvgPerCard)
		//		//valiGameCt = true
		//		break
		//	}
		//}
		if _gameCt == 6 {
			_tableConfig.FangKa = 1 * int(_playerCt)
			_tableConfig.AvgPerCard = int(1)
		} else if _gameCt == 12 {
			_tableConfig.FangKa = 2 * int(_playerCt)
			_tableConfig.AvgPerCard = int(2)
		}

	} else if _tableType == mjcomn.TableType_XGMJ {
		if _gameCt == 6 {
			_tableConfig.FangKa = 1 * int(_playerCt)
			_tableConfig.AvgPerCard = int(1)
		} else if _gameCt == 12 {
			_tableConfig.FangKa = 2 * int(_playerCt)
			_tableConfig.AvgPerCard = int(2)
		}
	}

	//tid, err := playerSvr.GetTableID()
	tid, err := requestTableID(playerSvr, int(uid),
		_tableConfig.GetItype(), 0, 0)
	if err != nil || tid <= 0 {
		logs.Error("create table error: %v", err.Error())
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(uid),
				MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
				Result: proto.Int32(int32(consts.ErrorTableNotExist)),
				Tip:    proto.String(err.Error()),
			},
		}
		writeResponse(pcmd)
		return
	}
	_tableConfig.TableId = tid

	//创建桌子
	table := CreateByType(TableMap, _tableConfig)

	table.Go(func() {

		needKaCt := _tableConfig.FangKa

		if _tableConfig.Present == consts.Yes { //赠送房间
			needKaCt = _tableConfig.AvgPerCard * _tableConfig.PlayerCt
			table.PresentUID = _player.ID() //保存赠送者uid
		} else {
			//创建房间,平摊房卡
			if _tableConfig.PayWay == 0 {
				needKaCt = _tableConfig.AvgPerCard
			} else {
				needKaCt = _tableConfig.AvgPerCard * _tableConfig.PlayerCt
			}

		}

		if config.Opts().OpenCharge {

			//判断房卡 -----------------------------
			if playerSvr.GetUsableFangka(_player.ID()) < needKaCt {
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
						Result: proto.Int32(int32(consts.Yes)),
					},
				}
				writeResponse(pcmd)
				return
			}
		}

		//赠送房间 ------------------------------------------------------------------------
		if _tableConfig.Present == consts.Yes {

			if config.Opts().OpenCharge {
				//赠送房间,先扣除房卡
				logs.Info("-------------赠送房间--------------->tableId:%v,扣除房卡,PresentUID:%v,needKaCt:%v",
					table.ID, table.PresentUID, needKaCt)
				consumeid, err := table.GetHandler().ConsumeFangka(table.PresentUID, table.ID, needKaCt)
				if err != nil {
					//提示扣卡失败信息
					pcmd := &protobuf.ResponseCmd{
						Head: &protobuf.RespHead{
							Uid:    proto.Int32(int32(_player.ID())),
							MsgID:  proto.Int32(0),
							Result: proto.Int32(0),
						},
						Simple: &protobuf.RespSimple{
							Tag:      protobuf.RespSimple_WIN_TIP_M.Enum(),
							StrValue: proto.String(err.Error()),
						},
					}
					writeResponse(pcmd)

					//发送创建桌子失败
					pcmd = &protobuf.ResponseCmd{
						Head: &protobuf.RespHead{
							Uid:    proto.Int32(int32(uid)),
							MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
							Result: proto.Int32(int32(consts.Yes)),
						},
					}
					writeResponse(pcmd)
					return
				}

				//赠送桌子成功
				table.ConsumeID = consumeid
			}

			pcmd := &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(uid)),
					MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
					Result: proto.Int32(0),
				},
				PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
				Simple: &protobuf.RespSimple{
					Tag:      protobuf.RespSimple_PRESENT_SUCC.Enum(),
					IntValue: proto.Int32(int32(table.ID)),
				},
			}
			writeResponse(pcmd)

			//创建房间 ------------------------------------------------------------------------
		} else {
			table.JoinTable(_player)

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
			//发送创建桌子成功
			pcmd = &protobuf.ResponseCmd{
				Head: &protobuf.RespHead{
					Uid:    proto.Int32(int32(uid)),
					MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
					Result: proto.Int32(0),
				},
				PlayerInfo: protobuf.Helper.GetPlayerInfo(_player),
			}
			writeResponse(pcmd)

		}
	})

}

//根据客户端选择的参数创建桌子
func CreateByType(TableMap *table.Tables, _tableCfg *config.TableCfg) *table.Table {

	_tableType := _tableCfg.TableType

	if !config.Opts().CloseAutoRobot {
		_tableCfg.RobotCt = 3
	}

	if _tableType == mjcomn.TableType_HFMJ { //合肥麻将
		_tableCfg.MaxCardColorIndex = 3 //万筒条2-8
		table := TableMap.CreateHFTable(_tableCfg)
		return table.Table

	} else if _tableType == mjcomn.TableType_HZMJ { //红中麻将
		_tableCfg.MaxCardColorIndex = 5  //万筒条1-9 中
		_tableCfg.KehuQidui = consts.Yes //七对默认
		_tableCfg.MenQing = consts.Yes   //门清默认
		table := TableMap.CreateHZTable(_tableCfg)
		return table.Table

	} else if _tableType == mjcomn.TableType_FYMJ { //阜阳麻将
		_tableCfg.MaxCardColorIndex = 5 //万筒条1-9 风 中
		//_tableCfg.KehuQidui = consts.No  //不带七对
		_tableCfg.QiangGang = consts.Yes //可抢杠
		table := TableMap.CreateFYTable(_tableCfg)
		return table.Table
	} else if _tableType == mjcomn.TableType_XGMJ { //香港麻将

		_tableCfg.MaxCardColorIndex = 6  //万筒条1-9 风 中
		_tableCfg.KehuQidui = consts.Yes //不带七对
		_tableCfg.QiangGang = consts.Yes //可抢杠
		_tableCfg.DianpaoHu = consts.Yes
		table := TableMap.CreateXGTable(_tableCfg)
		return table.Table
	} else if _tableType == mjcomn.TableType_HYMJ { //怀远麻将

		_tableCfg.MaxCardColorIndex = 6  //万筒条1-9 风 中 花
		_tableCfg.KehuQidui = consts.Yes //不带七对
		_tableCfg.QiangGang = consts.Yes //可抢杠
		table := TableMap.CreateHYTable(_tableCfg)
		return table.Table
	}

	return nil
}
