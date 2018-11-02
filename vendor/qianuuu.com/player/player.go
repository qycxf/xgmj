// 玩家对象

package player

import (
	"strconv"
	"time"

	"qianuuu.com/lib/qo"
	"qianuuu.com/player/micros/pb"
)

// Player 玩家对象
type Player struct {
	// User *domain.User
	user *pb.User

	tableID         int
	isFangZhu       bool  //是否房主
	lastMsgRecvTime int64 //消息接收时间
	isOffline       bool  //指定的时间内未收到玩家消息,认为离线

	//游戏模块信息
	ql *qo.Qo
}

// NewRobot 创建一个机器人
func NewRobot(id int, name string) *Player {
	ret := &Player{
		user: &pb.User{
			Uid:      int32(id),
			Nickname: name,
		},
		tableID:         id,
		isFangZhu:       false,
		lastMsgRecvTime: time.Now().Unix(),
		isOffline:       false,
		ql:              qo.New(),
	}
	return ret
}

// Go 异步调用 (线性，有顺序)
func (p *Player) Go(fn func()) {
	if p != nil && p.ql != nil {
		p.ql.Go(fn)
	} else {
		go fn()
		// logs.Info("player go func is nil")
	}
}

// 常用属性

// ID 玩家 ID
func (p *Player) ID() int {
	return int(p.user.Uid)
}

// NickName 昵称
func (p *Player) NickName() string {
	return p.user.Nickname
}

// SetTableID 设置桌子 ID
func (p *Player) SetTableID(tid int) {
	p.tableID = tid
}

// SetFangZhu 设置房主
func (p *Player) SetFangZhu(v bool) {
	p.isFangZhu = v
}

// SetOffline 离线
func (p *Player) SetOffline(v bool) {
	p.isOffline = v
}

// GetTableID 获取 tid
func (p *Player) GetTableID() int {
	return p.tableID
}

// Reset 重置
func (p *Player) Reset() {
	p.tableID = 0
	p.isFangZhu = false
}

func (p *Player) String() string {
	return "[" + p.user.Nickname + " " + strconv.Itoa(int(p.user.Uid)) + "]"
}

// IsRobot 是否机器人
func (p *Player) IsRobot() bool {
	if p.user.Uid <= 10000 {
		return true
	}
	return false
}

// IsFangZhu 是否房主
func (p *Player) IsFangZhu() bool {
	return p.isFangZhu
}

// IsOffline 是否离线
func (p *Player) IsOffline() bool {
	return p.isOffline
}

// UpdateMsgRecevTime 更新接收消息时间
func (p *Player) UpdateMsgRecevTime() {
	p.lastMsgRecvTime = time.Now().Unix()
}
