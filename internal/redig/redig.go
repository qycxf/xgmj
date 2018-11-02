//
// Author: leafsoar
// Date: 2017-05-18 14:05:25
//

// redis 网关，这个网关替代 gate 的功能
// 可以从 redis 接受来自大厅的消息，可以写回消息

package redig

import (
	"golang.org/x/protobuf/proto"
	"qianuuu.com/ahmj/internal/protobuf"
	"qianuuu.com/hall/lib/redic"
	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/util"
)

// Gate redis 网关
type Gate struct {
	hallchs  *util.Map                 // 大厅通道，可以有多个
	defchs   *Channel                  // 默认游戏通道
	readchan chan *protobuf.RequestCmd // 默认读消息通道
}

// NewGate 创建网关
func NewGate() *Gate {
	gate := &Gate{
		readchan: make(chan *protobuf.RequestCmd, 20),
	}
	return gate
}

// AddHallAddr 添加大厅地址
func (g *Gate) AddHallAddr(addr, name string) {
	logs.Info("[redig] add hall addr %s name %s", addr, name)
	// 将当前地址设置为默认地址
	// TODO: 注意如果是动态改变的，需要重构
	ch := NewChannel(addr, name)
	g.defchs = ch
}

// ReadMessage 读取消息
func (g *Gate) ReadMessage() *protobuf.RequestCmd {
	cmd := <-g.readchan
	//logs.Debug("[redig] read message %v", cmd.GetHead())
	return cmd
}

// WriteMessage 写入消息
func (g *Gate) WriteMessage(cmd *protobuf.ResponseCmd) {
	if g.defchs == nil {
		logs.Error("[redig] not found hall chal to write message %v", cmd.GetHead())
	}
	SetRespIDs(cmd)
	//logs.Debug("[redig] write message : %v", RespCmdInfo(cmd))
	data, err := proto.Marshal(cmd)
	if err == nil {
		err = g.defchs.WriteMessage(data)
	}
	if err != nil {
		logs.Error("[redig] write message err %s", err.Error())
	}

}

// TestRequestCmd 测试消息
func (g *Gate) TestRequestCmd(cmd *protobuf.RequestCmd) {
	g.readchan <- cmd
}

// 从 redis 接受 request 消息
func (g *Gate) handlerRedisMessage(addr string, name string) {
	logs.Info("[redig] listener redis game %s %s", addr, name)
	rc := redic.NewRedisClient(addr)
	if _, err := rc.Do("del", name); err != nil {
		logs.Info("[redig] client %s err %s", name, err.Error())
	}
	rc.PopByte(name, func(data []byte, err error) {
		// 读取到 redis 消息
		if err != nil {
			logs.Error("[redig] read redis data : %s", err.Error())
		} else {
			// 转换消息为 proto 对象
			qcmd := &protobuf.RequestCmd{}
			if err := proto.Unmarshal(data, qcmd); err != nil {
				logs.Error("[redig] read redis not unmarshal err : %s", err.Error())
			} else {
				g.readchan <- qcmd
			}
		}
	})
	logs.Info("[redig] read message from hall %v end", addr)
}

// Serve 开始服务
func (g *Gate) Serve(addr string, name string) {
	// name = _list:game:1
	g.handlerRedisMessage(addr, name)
}
