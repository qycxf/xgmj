//
// Author: leafsoar
// Date: 2017-05-31 17:29:27
//

package micros

import (
	"fmt"
	"time"
)

var (
	timecode chan int
)

func init() {
	timecode = make(chan int, 5)

	// 当前标记时间
	flagtime := time.Now()
	flagcode := 0
	go func() {
		for {
			now := time.Now().Add(time.Minute * 117)
			code := time2Code(now)
			// code 正常向下取值 （code 小于 flagcode 很多，说明是一次新循环）
			if code > flagcode || code < flagcode-800000 {
				flagtime = now
				flagcode = code
				timecode <- flagcode
			} else {
				// 加上 50 毫秒，获取上次之后的数据，不可能重复
				flagtime = flagtime.Add(time.Millisecond * 50)
				code := time2Code(flagtime)
				flagcode = code
				timecode <- flagcode
			}
		}
	}()

	// _test()
	// tr := codeRand(123456)
	// fmt.Println("tr:", tr)
}

// 将时间转换位 6 位数字
func time2Code(t time.Time) int {
	const mod = 43200000                 // 十二小时内的毫秒数
	code := t.UnixNano() / 1000000 % mod // 获取十二小时内的毫秒数
	ret := float64(code)/48.00001 + 100000.0
	return int(ret)
}

func codeRand(code int) int {
	if code < 100000 || code >= 1000000 {
		return code
	}
	// 为了看起来更随机，timecode 需要转换一下，个位数字放到万位，百千互换
	ivs := [6]int{}
	tpv := code
	m := 10
	for i := 0; i < 6; i++ {
		civ := tpv % 10
		tpv = tpv / 10
		ivs[5-i] = civ
		m = m * 10
	}
	// 转换
	// fmt.Println("m:", ivs)
	ivs[1], ivs[2], ivs[4], ivs[5] = ivs[5], ivs[1], ivs[2], ivs[4]
	// fmt.Println("m:", ivs)
	// 转成 int 类型
	ret := ivs[0]*100000 + ivs[1]*10000 + ivs[2]*1000 +
		ivs[3]*100 + ivs[4]*10 + ivs[5]*1
	// fmt.Println("end: ", ret)
	return ret
}

func getTableID() (int, error) {
	tc := <-timecode
	ret := codeRand(tc)
	// fmt.Println("tc: ", tc, " ret: ", ret)
	return ret, nil
}

func putTableID(tid int) {

}

func _test() {

	for i := 0; i < 60; i++ {
		time.Sleep(time.Second)
		ret, err := getTableID()
		fmt.Println("code:", ret, err)
	}

	// now := time.Now()
	// now = now.Add(time.Minute * 153)

	// times := []time.Time{}

	// for i := 0; i < 20; i++ {
	// 	now = now.Add(time.Second * 5)
	// 	times = append(times, now)
	// }

	// for _, t := range times {
	// 	fmt.Println("code: ", time2Code(t))
	// 	time.Sleep(time.Second)
	// }

}
