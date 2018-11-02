package table

import (
	"qianuuu.com/ahmj/internal/game/seat"
	. "qianuuu.com/ahmj/internal/mjcomn"
	"qianuuu.com/player"
)

// MsgHandler 发送消息接口
type Handler interface {
	RobotEnterInfo(_table *Table, _robot *player.Player)     //机器人进入
	SendTableInfo(_table *Table)                             //发送TableInfo
	KickPlayerInfo(_table *Table, _uID int)                  //踢出玩家消息
	DisMissTable(_table *Table, _seatID int)                 //解散桌子消息
	DestroyTable(_table *Table)                              //桌子销毁回调
	ConsumeFangka(_uid int, _tid, count int) (int, error)    //消耗房卡
	ConsumeFangkaLoss(_uid int, consumeid int)               //返还房卡
	MultiConsumeFangka(uids []int, tid, count int) error     // 多人消耗房卡
	SendTopTip(_table *Table, tipStr string)                 //发送TopTip
	AddZhanji(_table *Table, _seq int, score []player.Score) //添加战绩
}

//桌子重写方法接口
type TableInter interface {
	MakeDSeat()             //确定庄家
	HuPai(_seat *seat.Seat) //玩家胡牌
	ThinkSelfPai(_specialPxId int, ispeng bool) bool  //思考自己牌
	ThinkOtherPai(_tkCard *MCard) bool            //思考他人牌
	CheckHuPai(_seatId int) bool                  //检测胡牌
	Check_AnGang(_seatId int) bool                //暗杠检测
	DPlayerFirstThink()                           //庄家第一次思考
	FetcherStartThink(_specialPxId int)           //拿牌者拿牌后开始思考
	LastHuThinkerCancer()                         //最后一个可胡玩家选择取消胡
	SendCheckTing()                               //等待出牌->听牌检测
	UpdateTingCards(_seatId int)                  //更新位置当前听牌
	Gang(_seatId int, _paiData int, _isSave bool) //玩家杠牌
	CalSingle()                                   //小局结算
	Peng(_seatId int)                             //玩家碰牌
	GameStart()                                   //游戏开始
	ChkOver() bool                                //检测单局游戏结束

}
