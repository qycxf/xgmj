// 机器人对象

package table

import (
	"strconv"

	"qianuuu.com/ahmj/internal/consts"
	"qianuuu.com/lib/util"
)

// Robots 机器人集合
type Robots struct {
	robots   *util.Map
	robotIds int //机器人id 1~99 ,不可重用
}

// Robot 机器人
type Robot struct {
	robotID   int
	robotName string //昵称
}

// 创建机器人管理对象
func NewRobots() *Robots {
	ret := &Robots{
		robots:   &util.Map{},
		robotIds: 1,
	}
	return ret
}

// NewRobot 创建机器人
func (rs *Robots) NewRobot() *Robot {
	_robotID := rs.robotIds
	rs.robotIds++
	if rs.robotIds > consts.RobotMaxUid {
		rs.robotIds = 1
	}
	ret := &Robot{
		robotID:   _robotID,
		robotName: "我是机器人_" + strconv.Itoa(_robotID),
	}

	rs.robots.Set(_robotID, ret)
	return ret
}

// 删除机器人
func (rs *Robots) RemoveRobot(_robotID int) {
	rs.robots.Del(_robotID)
}

// ThinkOpt 机器人操作思考 optArr:当前机器人可以进行的操作 依次为: 等待 跟注 加注 看牌 比牌 放弃 开牌
func (r *Robots) ThinkOpt(optArr [7]int) int {

	//canFollow := optArr[consts.SeatOptFollow] > 0
	//canRaise := optArr[consts.SeatOptRaise] > 0
	//canCompare := optArr[consts.SeatOptCompare] > 0
	//canOpen := optArr[consts.SeatOptOpen] > 0
	//
	////能比牌
	//if canCompare {
	//	return consts.SeatOptCompare
	//}
	//
	////能开牌
	//if canOpen {
	//	return consts.SeatOptOpen
	//}
	////能跟注
	//if canFollow {
	//	//能加注
	//	if canRaise {
	//		//raiseRate := 100 //加注概率
	//		//random := rand.New(rand.NewSource(time.Now().UnixNano()))
	//		//r = random.Intn(100)
	//		//
	//		//if r < raiseRate {
	//		//	return consts.SEAT_OPT_RAISE
	//		//}
	//		return consts.SeatOptRaise
	//
	//	}
	//	return consts.SeatOptFollow
	//}
	//
	//return consts.SeatOptGiveup
	return 1

}
