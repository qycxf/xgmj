package protoapi

import (
	"math/rand"
	"time"

	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/player"
)

//游戏世界循环
func GameCircle(playerSvr player.Server) {

	//logs.Info("~~~~~~~~~~~GameCircle~~~~~~~~~~~")

	playerSvr.ReadRange(func(_uid int, _player *player.Player) {
		if _player.IsRobot() {
			return
		}

		//发送聊天内容
		_fromUID := 0
		_fromNickName := "系统"
		_strVale := []string{"抵制不良游戏,拒绝赌博！"}
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		ranIndex := random.Intn(len(_strVale))
		pcmd := &protobuf.ResponseCmd{
			Head: &protobuf.RespHead{
				Uid:    proto.Int32(int32(_player.ID())),
				MsgID:  proto.Int32(0),
				Result: proto.Int32(0),
			},
			WordChat: &protobuf.WordChat{
				UID:      proto.Int32(int32(_fromUID)),
				NickName: proto.String(_fromNickName),
				Content:  proto.String(_strVale[ranIndex]),
				Type:     proto.Int32(int32(consts.WorldChatTypeSystem)),
			},
		}
		writeResponse(pcmd)
	})

}
