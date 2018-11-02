// 玩家管理

package player

import (
	"time"

	"google.golang.org/grpc"

	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/qo"
	"qianuuu.com/lib/util"
	"qianuuu.com/lib/values"
	mc "qianuuu.com/player/micros/client"
)

// Server 提供的玩家数据服务模块，完全独立的模块
type Server interface {
	Init(string, string) error

	CheckToken(int, string) error
	Login(int) (*Player, error)
	Logout(int)       // 登出游戏
	IsLogin(int) bool // 是否登录

	GetPlayer(int) *Player        //从内存中取玩家数据
	ReadRange(func(int, *Player)) // 玩家读取遍历

	// 物品相关
	// GetFangka(int) int                                         // 获取用户房卡数
	GetUsableFangka(uid int) int                               // 获取用户可用房卡数 (房卡 + 钻石)
	ConsumeFangka(uid, tid, count int) (int, error)            // 消耗房卡
	ConsumeFangkaLoss(uid, consumeid int) error                // 消耗房卡失败
	MultiConsumeGoods(uids []int, tid, count int) (int, error) // 多人消耗房卡
	// 添加牌局记录 (桌子 id, 牌局, 用户到当前局的总分 第一个是房主， record 为 json 记录)
	AddTableRecord(tableid, inning int, scores Scores, record []byte) (int, error)

	//金币场接口
	GetCoin(int) int                                                                        // 获取用户金币
	ConsumeCoin(uid, tid, count int, parms values.ValueMap) (int, error)                    // 消耗金币
	EarnCoin(uid, tid, count int, parms values.ValueMap) (int, error)                       // 赚得金币
	ChargeCoin(uids []int, tid, count int, parms values.ValueMap) (int, error)              // 多人消耗服务费
	WinLossCoin(src, des, tid, count int) error                                             // 双人输赢,src输掉牌桌的玩家id，des赢得牌桌的玩家id
	MultiWinLossCoin(uids []int, counts []int, tid int, parms values.ValueMap) (int, error) // 多人消耗赢得金币

	GetTableID(parms values.ValueMap) (int, error)
	PutTableID(int) error

	IsTimeout(int, int64) bool // 玩家是否操作超时
	PlayerCount() int          //玩家数量
}

// Players 玩家集合
type service struct {
	players *util.Map
	client  *mc.Client
	appname string
}

// NewServer 创建一个用户服务模块
func NewServer() Server {
	ret := &service{
		players: &util.Map{},
	}
	return ret
}

func (svr *service) Init(msurl, appname string) error {
	svr.appname = appname
	svr.client = mc.NewClient(msurl, grpc.WithInsecure())
	if err := svr.client.SayHello(); err != nil {
		logs.Error("[service] rpc err: %v", err)
	}
	return nil
}

// GetPlayer 根据uid获取内存玩家
func (svr *service) GetPlayer(_uid int) *Player {
	p := svr.players.Get(_uid)
	if p != nil {
		return p.(*Player)
	}
	return nil
}

func (svr *service) ReadRange(fn func(int, *Player)) {
	ps := make([]*Player, 0, svr.players.Len())
	svr.players.LockRange(func(k interface{}, v interface{}) {
		player := v.(*Player)
		ps = append(ps, player)
	})
	// 在这里回调，防止 RLock 锁住
	for _, p := range ps {
		fn(p.ID(), p)
	}
}

// CheckToken 验证 token
func (svr *service) CheckToken(uid int, token string) error {
	_, err := svr.client.CheckToken(uid, token)
	return err
}

// 创建一个内存玩家
func (svr *service) Login(uid int) (*Player, error) {
	user, err := svr.client.UserInfo(uid)
	if err != nil {
		return nil, err
	}
	player := &Player{
		user:            user,
		tableID:         0,
		isFangZhu:       false,
		lastMsgRecvTime: time.Now().Unix(),
		isOffline:       false,
		ql:              qo.New(),
	}
	svr.players.Set(player.ID(), player)
	logs.Custom(logs.PlayerTag, "玩家登录: %v", player)
	return player, err
}

// 从内存中删除
func (svr *service) Logout(uid int) {
	logs.Info("玩家登出: %d", uid)
	//登出之前保存游戏数据
	svr.players.Del(uid)
}

func (svr *service) IsLogin(uid int) bool {
	player := svr.players.Get(uid)
	return player != nil
}

func (svr *service) AddTableRecord(tableid, inning int, scores Scores, record []byte) (int, error) {
	records, errs := values.NewValuesFromJSON(record)
	appname := records.GetString("appName")
	if errs != nil {
		appname = svr.appname
	}
	ret, err := svr.client.AddTableRecord(appname, tableid, inning, scores.Uids(), scores.String(), record)
	if err != nil {
		return 0, err
	}
	if ret == nil {
		return 0, err
	}
	return int(ret.IntVal), nil
}

func (svr *service) PlayerCount() int {
	return svr.players.Len()
}

// GetCoin 获取用户金币
func (svr *service) GetCoin(uid int) int {
	user, err := svr.client.UserInfo(uid)
	if err != nil {
		return 0
	}
	return getGoods(user.Goods, "coin")
}

// ConsumeCoin 消耗金币
func (svr *service) ConsumeCoin(uid, tid, count int, parms values.ValueMap) (int, error) {
	id, err := svr.client.ConsumeGoods(uid, "coin", count, svr.appname, tid, parms)
	return id, err
}

// EarnCoin 赚得金币
func (svr *service) EarnCoin(uid, tid, count int, parms values.ValueMap) (int, error) {
	id, err := svr.client.EarnGoods(uid, "coin", count, svr.appname, tid, parms)
	return id, err
}
func (svr *service) ChargeCoin(uids []int, tid, count int, parms values.ValueMap) (int, error) {
	cid := 0
	var err error
	for i := 0; i < len(uids); i++ {
		cid, err = svr.client.ChargeCoin(uids[i], "coin", count, svr.appname, tid, parms)
	}
	return cid, err
}

func (svr *service) WinLossCoin(src, des, tid, count int) error {
	mp := map[string]interface{}{
		"type": "winloss",
	}
	err := svr.client.WinLossCoin(src, "coin", des, tid, count, svr.appname, mp)
	return err
}

//MultiConsumeEarnCoin
func (svr *service) MultiWinLossCoin(uids []int, counts []int, tid int, parms values.ValueMap) (int, error) {
	cid, err := svr.client.MultiWinLossCoin(uids, "coin", counts, svr.appname, tid, parms)
	return cid, err
}

// func (svr *service) GetFangka(uid int) int {
// 	user, err := svr.client.UserInfo(uid)
// 	if err != nil {
// 		return 0
// 	}
// 	return getGoods(user.Goods, "fangka")
// }

func (svr *service) GetUsableFangka(uid int) int {
	user, err := svr.client.UserInfo(uid)
	if err != nil {
		return 0
	}
	// 可用的房卡数等于房卡加钻石
	ret := getGoods(user.Goods, "fangka")
	return ret
}

func (svr *service) ConsumeFangka(uid, tid, count int) (int, error) {
	id, err := svr.client.ConsumeGoods(uid, "fangka", count, svr.appname, tid, nil)
	return id, err
}

func (svr *service) ConsumeFangkaLoss(uid, consumeid int) error {
	err := svr.client.ConsumeGoodsLoss(consumeid)
	return err
}

func (svr *service) MultiConsumeGoods(uids []int, tid, count int) (int, error) {
	cid, err := svr.client.MultiConsumeGoods(uids, "fangka", count, svr.appname, tid)
	return cid, err
}

func (svr *service) GetTableID(params values.ValueMap) (int, error) {
	ps := params
	if ps == nil {
		ps = make(values.ValueMap)
	}
	// return svr.client.GetTableID(svr.appname)
	if ps.GetString("appname") == "" {
		ps["appname"] = svr.appname
	}
	return svr.client.GetTableID(ps)
}

func (svr *service) PutTableID(tid int) error {
	return svr.client.PutTableID(tid)
}

func (svr *service) IsTimeout(uid int, outtime int64) bool {
	p := svr.GetPlayer(uid)
	if p == nil {
		return false
	}
	// 玩家最后更新时间是否超市
	timeoffset := time.Now().Unix() - p.lastMsgRecvTime
	return timeoffset >= outtime
}

// getGoods 获取物品个数
func getGoods(goods, name string) int {
	vm, err := values.NewValuesFromJSON([]byte(goods))
	if err != nil {
		return 0
	}
	return vm.GetInt(name)
}
