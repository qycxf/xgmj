//
// Author: leafsoar
// Date: 2016-09-01 17:48:06
//

package qo

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"

	"qianuuu.com/lib/logs"
)

// Qo 线性 gorouting 队列
type Qo struct {
	linearQo       *list.List
	mutexLinearQo  sync.Mutex
	mutexExecution sync.Mutex
}

// LinearQo 线性
type LinearQo struct {
	fn func()
}

// New 创建 Qo
func New() *Qo {
	q := &Qo{
		linearQo: list.New(),
	}
	return q
}

// Go 调用
func (q *Qo) Go(fn func()) {
	q.mutexLinearQo.Lock()
	q.linearQo.PushBack(&LinearQo{fn: fn})
	q.mutexLinearQo.Unlock()

	go func() {
		q.mutexExecution.Lock()
		defer q.mutexExecution.Unlock()

		q.mutexLinearQo.Lock()
		lq := q.linearQo.Remove(q.linearQo.Front()).(*LinearQo)
		q.mutexLinearQo.Unlock()

		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logs.Info("%v", fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
					err, n, trace[:n]))
			}
		}()
		lq.fn()
	}()
}
