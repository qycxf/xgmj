//
// Author: leafsoar
// Date: 2017-12-19 11:06:03
//

package gath

// GameCommand 游戏命令
type GameCommand interface {
	DismissTable(tid int) error
}
