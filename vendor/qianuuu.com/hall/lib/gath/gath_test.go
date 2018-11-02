//
// Author: leafsoar
// Date: 2017-12-19 10:33:33
//

package gath

import (
	"fmt"
	"testing"
	"time"

	"qianuuu.com/lib/values"
)

type GameHandler struct {
}

func (gh *GameHandler) DismissTable(tid int) error {
	fmt.Println("解散牌桌 ", tid)
	return nil
}

func TestGath(t *testing.T) {
	fmt.Println("test gath ...")

	gh := New("testgame", "test.qianuuu.cn:7599")
	go func() {
		gh.ListeningGameCommand(&GameHandler{})
	}()
	vm := values.ValueMap{
		"command_type": "dismiss_table",
		"table_id":     10001,
	}
	go func() {
		for i := 0; i < 1; i++ {
			gh.writedata("_list:podk_t1:cmd", vm.ToJSON())
			time.Sleep(time.Second * 1)
		}
	}()

	gh.SetGameState(values.ValueMap{
		"count": 80,
	})
	time.Sleep(time.Second * 100)
}
