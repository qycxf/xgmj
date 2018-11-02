//
// Author: leafsoar
// Date: 2016-10-09 18:36:58
//

package player

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONTime 时间格式
type JSONTime time.Time

// MarshalJSON 时间 format
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// Score 记录分数
type Score struct {
	UID      int    `json:"uid"`
	NickName string `json:"nickname"`
	Score    int    `json:"score"`
}

// Scores 分数集合
type Scores []Score

// Uids 用户
func (ss Scores) Uids() string {
	if len(ss) <= 0 {
		return ""
	}
	uids := make([]int, 0, 4)
	for _, v := range ss {
		uids = append(uids, v.UID)
		// uids = append(uids, strconv.Itoa(v.UID))
	}
	ret, _ := json.Marshal(uids)
	// ret := "[" + strings.Join(uids, ",") + "]"
	return string(ret)
}

// String 分数
func (ss Scores) String() string {
	if len(ss) <= 0 {
		return ""
	}
	ret, _ := json.Marshal(&ss)
	return string(ret)
}

// FromJSON 从 json 解析
func (ss Scores) FromJSON(data []byte) (Scores, error) {
	ret := json.Unmarshal(data, &ss)
	return ss, ret
}

// GetScore 获取指定用户的分数
func (ss Scores) GetScore(uid int) *Score {
	for _, item := range ss {
		if item.UID == uid {
			return &item
		}
	}
	return nil
}
