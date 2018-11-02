//
// Author: leafsoar
// Date: 2016-09-02 14:27:02
//

package qo

import (
	"fmt"
	"time"

	"qianuuu.com/lib/logs"
)

// TimeAnalyse 耗时统计
type TimeAnalyse struct {
	startTime time.Time
	timeOut   time.Duration
}

// NewTimeAnalyse 添加一个耗时统计
func NewTimeAnalyse() *TimeAnalyse {
	return &TimeAnalyse{
		startTime: time.Now(),
		timeOut:   time.Millisecond * 500, // 默认 500 毫秒以上的操作统计
	}
}

// SetTimeOut 设置超时时间
func (ta *TimeAnalyse) SetTimeOut(timeout time.Duration) {
	ta.timeOut = timeout
}

// TimeOut 超时
func (ta *TimeAnalyse) TimeOut(format string, a ...interface{}) {
	offset := time.Now().Sub(ta.startTime)
	if offset > ta.timeOut {
		prestr := fmt.Sprintf("[time] %d ms ", offset/time.Millisecond)
		logs.Warning(prestr+format, a...)
	}
}
