//
// Author: leafsoar
// Date: 2017-12-19 10:29:28
//

// 大厅网关，游戏中不需要网关 gate fro hall

package gath

import (
	"fmt"

	"qianuuu.com/hall/lib/redic"
	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/values"
)

type Gath struct {
	gname string
	haddr string
	chall *redic.Client
}

// New 游戏名称 大厅地址
func New(name, addr string) *Gath {
	return &Gath{
		gname: name,
		haddr: addr,
		chall: redic.NewRedisClient(addr),
	}
}

// SetGameState 设置当前游戏状态
func (g *Gath) SetGameState(params values.ValueMap) error {
	name := fmt.Sprintf("_state:%s", g.gname)
	g.chall.Set(name, string(params.ToJSON()))
	return nil
}

func (g *Gath) writedata(key string, data []byte) {
	err := g.chall.Push(key, data)
	fmt.Println(err)
}

// ListeningGameCommand 监听游戏消息，传入消息处理
func (g *Gath) ListeningGameCommand(handler GameCommand) {
	name := fmt.Sprintf("_list:%s:cmd", g.gname)
	logs.Info("[gath] listener game command %s %s", name, g.haddr)
	rc := redic.NewRedisClient(g.haddr)
	if _, err := rc.Do("del", name); err != nil {
		logs.Info("[gath] client %s err %s", name, err.Error())
	}
	rc.PopByte(name, func(data []byte, err error) {
		// 读取到 redis 消息
		if err != nil {
			logs.Error("[gath] read cmd data : %s", err.Error())
			return
		}
		vm, err := values.NewValuesFromJSON(data)
		if err != nil {
			logs.Error("[redig] read cmd data : %s", err.Error())
			return
		}
		// 根据内容处理消息

		if vm.GetString("command_type") == "dismiss_table" {
			tid := vm.GetInt("table_id")
			if tid > 0 {
				handler.DismissTable(tid)
			}
		}

	})
	logs.Info("[redig] read cmd from hall %v end", g.haddr)
}
