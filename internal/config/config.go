//
// Author: leafsoar
// Date: 2016-05-14 10:31:32
//

package config

import (
	"bytes"
	"io/ioutil"
	"os"

	"qianuuu.com/lib/logs"
	"qianuuu.com/xgmj/internal/mjcomn"

	"github.com/BurntSushi/toml"
	"qianuuu.com/xgmj/internal/consts"
)

// Config 配置类型
type Config struct {

	//system-------------------
	AppName     string // 游戏名称
	RPCUrl      string
	ReceiveRdds string // 接受消息地址 redis
	HallRdds    string // 大厅地址 redis
	LogPath     string // 日志文件目录
	LogLevel    string // 日志等级

	OpenHeart bool // 开启心跳检测
	CloseTest bool // 关闭测试模式
	OpenDb    bool // 开启数据库

	//game-------------------
	OpenHu         bool // 开启胡牌检测
	OpenPeng       bool // 开启碰牌检测
	OpenGang       bool // 开启杠牌检测
	CloseTestPx    bool //关闭测试牌型
	OpenReConnect  bool //开启断线重连
	OpenCharge     bool //计费开关
	CloseAutoRobot bool //机器人自动进入

	Rooms   []Room // 房间配置
	FangKas []FangKa

	SelectScoreCt   []int32    //牌开分数选择
	SelectTdqZuiCt  []int32    //天地胡\清一色分数选择
	RoomCards       []RoomCard //房间局数房卡配置
	RewardFangKaNum int        //输入邀请码增房卡数
	NoticeStr       string     //大厅公告
}

// Opts Config 默认配置
var opts = Config{
	AppName:        "hzmj",
	OpenHeart:      true,
	CloseTest:      true,
	OpenDb:         true,
	OpenHu:         true,
	OpenPeng:       true,
	OpenGang:       true,
	CloseTestPx:    true,
	OpenReConnect:  true,
	OpenCharge:     true,
	CloseAutoRobot: true,

	Rooms:           []Room{},
	FangKas:         []FangKa{},
	SelectScoreCt:   []int32{},
	SelectTdqZuiCt:  []int32{},
	RoomCards:       []RoomCard{},
	RewardFangKaNum: 0,
	NoticeStr:       "",
}

// Opts 获取配置
func Opts() Config {
	return opts
}

//创建房间房卡
type RoomCard struct {
	GameCt     int32 // 创建房间局数
	NeedCard   int32 // 所需房卡数
	AvgPerCard int32 // 平均房卡数
}

// Room 房间
type Room struct {
	Name       string   // 房间名称
	BasMoney   int64    // 底注
	MinMoney   int64    // 最小准入
	Tip        int      // 台费 （百分比 0 - 100）
	RaiseChips [4]int64 // 固定加注额度
}

// 房卡
type FangKa struct {
	ID     int
	Name   string
	Count  int
	Price  int
	Reward int
}

// ParseToml 解析配置文件
func ParseToml(file string) error {
	logs.Info("读取配置文件 ...")
	// 如果配置文件不存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(Opts()); err != nil {
			return err
		}
		logs.Info("没有找到配置文件，创建新文件 ...")
		// logs.Info(buf.String())
		return ioutil.WriteFile(file, buf.Bytes(), 0644)
	}
	var conf Config
	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return err
	}
	opts = conf
	//logs.Info("----------->，config.Opts().FangKas:%v", opts.FangKas)

	return nil
}

//创建牌桌 返回平台对应的字段
func (tc *TableCfg) GetItype() int {
	//1:房卡场(预留字段)  10:好友场  11:菜鸟场12:平民场13:进阶场14:欢乐场15:龙虎场16:高手场

	if tc.TableType == mjcomn.TableType_XGMJ {
		return 1
	}
	//} else if tc.TableClass == pkcomn.TableClass_Coin {
	//	_room := GetRoomData(tc.RoomId)
	//	return _room.Id + 10
	//}
	return 0
}

//桌子参数配置
type TableCfg struct {
	MaxCardColorIndex int //所使用的麻将牌花色最大下标

	TableId    int //桌子id
	TestSeatId int

	TableType     int // 桌子类型
	PlayerCt      int // 玩家数
	GameCt        int // 游戏局数
	BaseScore     int // 牌局底分
	RobotCt       int // 机器人数
	DianpaoHu     int // 点炮胡
	ZimoHu        int // 自摸胡
	TiandiHu      int // 天地胡
	KehuQidui     int // 可胡七对
	QiangGang     int // 可抢杠
	ZhuaNiaoCt    int // 抓鸟数
	YiMaQuanZh    int // 一码全中
	MenQing       int // 门清
	Present       int // 赠送房间
	TdqZuiCt      int // 合肥麻将,天地胡\清一色嘴数
	KePengGang    int // 红中麻将,可碰杠
	KaiHuSuanGang int // 阜阳麻将,开胡算杠
	YouGangYouFen int // 阜阳麻将,有杠有分(没荒庄)
	CreaterId     int //创建房间者id
	DaiHua        int //蚌埠麻将是否带花
	FengLing      int //怀远麻将风令
	BaoTing       int //怀远麻将报听
	WuHuaGuo      int //怀远麻将无花果

	FangKa     int //所需房卡数
	AvgPerCard int //平均房卡数
	PayWay     int //付费方式

}

func NewTableCfg() *TableCfg {
	ret := &TableCfg{

		MaxCardColorIndex: 6,
		TableId:           0,
		TestSeatId:        consts.DefaultIndex,

		TableType: 0,         //桌子类型
		PlayerCt:  4,         //玩家数
		GameCt:    0,         //游戏局数
		BaseScore: 0,         //牌局底分
		RobotCt:   0,         //机器人数
		DianpaoHu: consts.No, //点炮胡
		ZimoHu:    consts.No, //自摸胡
		DaiHua:    consts.No, //是否带花
		FengLing:  consts.No, //怀远麻将风令
		BaoTing:   consts.No, //怀远麻将风令
		WuHuaGuo:  consts.No, //怀远麻将无花果

		TiandiHu:      consts.No, //天地胡
		KehuQidui:     consts.No, //可胡七对
		QiangGang:     consts.No, //可抢杠
		ZhuaNiaoCt:    0,         //抓鸟数
		YiMaQuanZh:    consts.No, //一码全中
		MenQing:       consts.No, //门清
		Present:       consts.No,
		TdqZuiCt:      0,
		KePengGang:    consts.No,
		KaiHuSuanGang: consts.No,
		YouGangYouFen: consts.No,
		CreaterId:     0,

		FangKa:     0, //所需房卡数
		AvgPerCard: 0,
	}

	return ret
}

func (ret *TableCfg) Print() {

	logs.Info("tableId:%v==============================>创建牌桌信息[TableCfg] --> 桌子类型:%v,天地胡:%v,可胡七对:%v,门清:%v,可抢杠:%v",
		ret.TableId, ret.TableType, ret.TiandiHu, ret.KehuQidui, ret.MenQing, ret.QiangGang)
}
