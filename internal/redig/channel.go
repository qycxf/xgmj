//
// Author: leafsoar
// Date: 2017-05-17 16:06:29
//

// 需要检测当前游戏通道是否正常

package redig

import "qianuuu.com/hall/lib/redic"

// Channel 游戏通道
type Channel struct {
	name   string
	client *redic.Client
}

// NewChannel 创建一个通道 (地址，名称)
func NewChannel(addr string, name string) *Channel {
	return &Channel{
		name:   name,
		client: redic.NewRedisClient(addr),
	}
}

// WriteMessage 写入消息
func (c *Channel) WriteMessage(data []byte) error {
	return c.client.Push(c.name, data)
}

// TODO 删除
func (c *Channel) TODO() {
	// 定期删除缓存中的数据，检测是否有积压的消息
}
