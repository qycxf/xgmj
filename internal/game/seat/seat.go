// 玩家对象

package seat

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player"
)

// Seat 玩家
type Seat struct {
	Id     int //座位id
	player *player.Player
	state  int //玩家游戏状态

	RemainTime int //剩余操作时间
	NotReadyCt int //连续不准备次数
}

// NewSeat 创建玩家
func NewSeat(_seatID int, _player *player.Player) *Seat {
	ret := &Seat{
		Id:     _seatID,
		player: _player,
	}
	ret.Init()
	return ret
}

//初始化信息
func (s *Seat) Init() {
	s.RemainTime = 0
	s.NotReadyCt = 0
	s.SetState(consts.SeatStateIdle)
}

//重置座位游戏信息
func (s *Seat) Reset() {
	//logs.Info("-------------->重置座位状态 seat:%v", s)
	s.RemainTime = consts.DefaultIndex

	s.SetState(consts.SeatStateIdle)

}

//清空操作时间
func (s *Seat) ClearRemainTime() {
	s.RemainTime = 0 //清除位置操作时间
}

// Seats 多个玩家
type Seats map[int]*Seat

func (p Seats) String() string {
	values := make([]string, 0, len(p))
	for _, seat := range p {
		if seat != nil {
			values = append(values, fmt.Sprintf("%v", seat))
		}
	}
	return strings.Join(values, "   ")
}

//// ForEach 遍历
//func (p Seats) ForEach(call func(*Seat)) {
//	for _, seat := range p {
//		call(seat)
//	}
//}

// Range 遍历
func (ss Seats) Foreach(f func(key int, s *Seat)) {
	keys := make([]int, 0, len(ss))
	for k, _ := range ss {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, key := range keys {
		f(key, ss[key])
	}
}

// GetSeatByUID 根据uid 返回seat
func (ss Seats) GetSeatByUID(_uid int) *Seat {
	for _, _seat := range ss {
		if _seat.GetPlayer().ID() == _uid {
			return _seat
		}
	}
	logs.Info("************************* game error --->call GetSeatByUID,not found _uid:%v", _uid)
	return nil
}

// GetSeatBySeatId 根据SeatId 返回seat
func (ss Seats) GetSeatBySeatId(_seatId int) *Seat {
	for _, seat := range ss {
		if seat != nil && seat.GetId() == _seatId {
			return seat
		}
	}
	logs.Info("************************* game error ----->call GetSeatBySeatId,not found _seatId:%v", _seatId)
	return nil
}

// RemoveSeat 删除seat
func (ss Seats) RemoveSeat(_seatId int) {
	len1 := ss.Count()
	delete(ss, _seatId)
	len2 := ss.Count()
	if len1-len2 != 1 {
		logs.Info("************************* game error RemoveSeat ,_seatId:%v, len1:%v, len2:%v", _seatId, len1, len2)
		return
	}
	logs.Info("RemoveSeat success ,_seatId:%v, ", _seatId)
}

// GetUIDBySeatID 根据座位号返回uid
func (p Seats) GetUIDBySeatID(seatID int) int {
	for _, seat := range p {
		if seat != nil && seat.GetId() == seatID {
			return seat.GetPlayer().ID()
		}
	}
	return -1
}

// Count 返回牌桌在座玩家数
func (ss Seats) Count() int {
	count := 0
	for _, seat := range ss {
		if seat != nil {
			count++
		}
	}
	return count
}

//位置信息
func (s *Seat) String() string {
	return "[位置" + strconv.Itoa(s.GetId()) + "]"
}

//---------------------------------------

func (s *Seat) GetPlayer() *player.Player {
	return s.player
}

func (s *Seat) GetState() int {
	return s.state
}
func (s *Seat) SetState(_state int) {
	s.state = _state
}

func (s *Seat) GetId() int {
	return s.Id
}
