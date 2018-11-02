package timer

import (
	"sort"
	"time"

	"qianuuu.com/player"
)

//排行数据
type Rank struct {
	Player player.Player
}

//按金币排序--------------------------------------------------
var RankCoinSorter = CoinSorter{}

type CoinSorter []Rank

func (ms CoinSorter) Len() int {
	return len(ms)
}

func (ms CoinSorter) Less(i, j int) bool {
	//return ms[i].Player.Coin > ms[j].Player.Coin // 按金币降序
	return true
}

func (ms CoinSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms CoinSorter) GetRankArr() []Rank {
	return ms
}

//刷新排行数据
func tickRank(playerSvr player.Server) {

	//logs.Info("tickRank")
	//TODO 查询数据库所有数据

	//取出内存数据

	rankCoinSorter := make(CoinSorter, 0)
	playerSvr.ReadRange(func(_uid int, _player *player.Player) {
		rankCoinSorter = append(rankCoinSorter, Rank{*_player})
	})

	//logs.Debug("memMap.Len():%v",memMap.Len())

	//TODO 合并数据后排序
	sort.Sort(rankCoinSorter)

	//logs.Info("len(RankCoinSorter):%v",len(RankCoinSorter))
	//for i:=0;i<len(RankCoinSorter) ;i++  {
	//	_player:=RankCoinSorter[i].Player
	//	logs.Info(" nick:%v, coin:%v",_player.User.NickName,_player.Coin)
	//}
}

func Serve(playerSvr player.Server) {
	for {
		select {
		//定时执行
		//case <-time.After(time.Second * 60*10):
		case <-time.After(time.Second * 1): //刷新排行
			tickRank(playerSvr)
			break
		}
	}
}
