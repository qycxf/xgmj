package protoapi

import (
	"qianuuu.com/xgmj/internal/protobuf"
	"qianuuu.com/player"
)

//请求牌局回放数据
func TableInfoRec(cmd *protobuf.RequestCmd, playerSvr player.Server) {

	// uid := cmd.Head.GetUid()
	// result := consts.Success
	// errfunc := func(_result int) {
	// 	pcmd := &protobuf.ResponseCmd{
	// 		Head: &protobuf.RespHead{
	// 			Uid:    proto.Int32(uid),
	// 			MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
	// 			Result: proto.Int32(int32(result)),
	// 			Tip:    proto.String(consts.GetErrTip(result)),
	// 		},
	// 	}
	// 	writeResponse(pcmd)
	// }

	// _player := playerSvr.GetPlayer(int(uid))
	// if _player == nil {
	// 	result = consts.ErrorNotLoginInGame
	// 	errfunc(result)
	// 	return
	// }

	// recId := cmd.GetSimple().GetIntValue() //数据id
	// tableRecord, err1 := playerSvr.GetZhanjiByID(int(recId))

	// if err1 != nil || tableRecord == nil {
	// 	result = consts.ErrorTableRecNotFound
	// 	errfunc(result)
	// 	return
	// }

	// tableInfo := make([]*protobuf.TableInfo, 0)
	// err2 := json.Unmarshal([]byte(tableRecord.RecordInfo), &tableInfo)
	// //logs.Info("------------------>tableRecord.RecordInfo, err:%v ", err2)
	// if err2 != nil {
	// 	result = consts.ErrorTableInfoData
	// 	errfunc(result)
	// 	return
	// }

	// tableInfoRec := &protobuf.RespTableInfoRec{
	// 	TableId:      proto.Int32(int32(tableRecord.TableID)),
	// 	Sequence:     proto.Int32(int32(tableRecord.Inning)),
	// 	TableInfoArr: tableInfo,
	// }
	// //请求成功
	// pcmd := &protobuf.ResponseCmd{
	// 	Head: &protobuf.RespHead{
	// 		Uid:    proto.Int32(int32(uid)),
	// 		MsgID:  proto.Int32(cmd.GetHead().GetMsgID()),
	// 		Result: proto.Int32(int32(result)),
	// 	},
	// 	TableInfoRec: tableInfoRec,
	// }
	// writeResponse(pcmd)

}
