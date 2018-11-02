//
// Author: leafsoar
// Date: 2017-08-10 15:53:58
//

package micros

import "qianuuu.com/player/domain"

// Handler 接口
type Handler interface {
	// 数据存储
	UserInfo(int) (*domain.User, error)
	AddRecord(tid, inning int, uids, scores string, data []byte, appname, path string) (*domain.TableRecord, error)
	AutoConsumeFangka(uid int, count int, detail string) (int, error)
	ConsumeGoodsLoss(changeid int) error
	EarnGoods(uid int, goodtype string, count int, detail string) (int, error)
	ChargeCoin(uid int, detail string) (int, error)
	WinLossCoin(uid int, detail string) error
	MultiWinLossCoin(detail string) (int, error)
	ConsumeCoin(uid int, goodtype string, count int, detail string) (int, error)
	// OSS 存储
	EnableOSS() bool
	PutTableRecord(content []byte) (string, error)
}
