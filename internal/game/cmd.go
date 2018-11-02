//
// Author: leafsoar
// Date: 2017-12-19 11:40:25
//

package game

import (
	"qianuuu.com/lib/logs"
)

type GameCmdHandler struct {
	g *Game
}

func (gh *GameCmdHandler) DismissTable(tid int) error {
	logs.Info("解散牌桌 :%v", tid)
	_table := gh.g.TableMap.GetTable(tid) //根据id获取牌桌
	if _table == nil {
		logs.Info("解散牌桌 %v,未找到牌桌 ", tid)
		return nil
	} else {
		go func() {
			//_table.Majhong.GameCt = 0
			//_table.SendTableInfo()
			gh.g.DoDisMiss(_table, "强制解散牌桌", true)
			logs.Info("解散牌桌 %v,销毁牌桌 ", tid)
		}()

	}
	return nil
}
