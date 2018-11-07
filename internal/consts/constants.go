// 游戏常量

package consts

import "strconv"

const (
	Success int = 0  // Success 成功
	Failure     = -1 // Failure 失败

	Yes int = 1 // Yes 是
	No  int = 0 // No 否

	RobotMaxUid  int = 999 //机器人uid 最大值
	DefaultIndex int = 999 //默认下标值,用该值替代-1

)

// TableState 牌桌状态
type TableState int

const (
	TableStateIdle       = 0 // 空闲
	TableStateWaiteReady = 1 // 等待准备
	TableStateDealCard   = 2 // 发牌状态
	TableStateChangeCard = 3 // 换三张状态
	TableStateSelectQue  = 4 // 选缺状态
	TableStateWaiteSend  = 5 // 等待玩家出牌
	TableStateWaiteThink = 6 // 等待玩家思考
	TableStateZhuaNiao   = 7 // 抓鸟状态
	TableStateShowResult = 8 // 结算状态
	TableStateForTest    = 9 // 测试状态
)

func (t TableState) String() string {
	return [10]string{"空闲", "等待准备", "发牌", "换三张状态", "选缺状态", "等待玩家出牌",
		"等待玩家思考", "抓鸟状态", "结算状态", "测试状态"}[t]
}

// 解散牌桌状态值
const (
	DismissTableStateWaite    = 0 // 等待
	DismissTableStateAgree    = 1 // 同意
	DismissTableStateDisAgree = 2 // 不同意
)

// 合肥麻将胡牌牌型\分数 ----------------------------------------------------------------------------
const (
	PXID_HFMJ_PINGHU          = 1 //平胡
	PXID_HFMJ_QYS             = 2 //清一色
	PXID_HFMJ_7DUI            = 3 //七对
	PXID_HFMJ_HAOHUA7DUI      = 4 //豪华七对
	PXID_HFMJ_SUPERHAOHUA7DUI = 5 //超豪华七对
	PXID_HFMJ_HAI_DI_LAO_YUE  = 6 //海底捞月
	PXID_HFMJ_GANG_SHANG_HUA  = 7 //杠上开花
	PXID_HFMJ_TIANHU          = 8 //天胡
	PXID_HFMJ_DIHU            = 9 //地胡
)

func GetHuPxName_HFMJ(_pxIndex int) string {
	return []string{"", "平胡", "清一色", "七对", "豪华七对", "超豪华七对", "海底捞月", "杠上开花", "天胡", "地胡"}[_pxIndex]
}

//合肥麻将胡牌牌型 "嘴" 数
func GetHuPxScore_HFMJ(_pxIndex int) int {
	//"杠上开花":总嘴数*2
	return []int{0, 15, 20, 10, 50, 100, 15, 0, 20, 20}[_pxIndex]
}

//额外加成的牌型
const (
	//合肥麻将
	EXTPXID_HFMJ_ZHI            = iota + 1 //［支］ ：    某一门有 8 张是胡牌的的基本要求，如有 9 张则加 1 支， 10 张加 2 支，以此类推，多支多分 ( 一支+1嘴) 。
	EXTPXID_HFMJ_KA                        //［卡］ ：    只胡一张牌的牌型所胡的牌称为卡。+1嘴 23胡4 78胡6也算卡
	EXTPXID_HFMJ_QUE_MENG                  //［缺门］：   胡牌时所有牌只存在有2门花色。+2嘴
	EXTPXID_HFMJ_TONG                      //［同］：     所有牌中数字一样的牌从 4 张起数， 每多一张多+1嘴（如， 3 张四万、 2 张四条、 3 张四饼的牌得 4倍每基础分）
	EXTPXID_HFMJ_SHI_TONG                  //［10同］：   同超过10张 +10嘴
	EXTPXID_HFMJ_KAN                       //［坎］：     三张一样的牌在手，且符合基本胡牌牌型中的刻 ，即三张牌未分开叫做一坎，每一坎+1嘴
	EXTPXID_HFMJ_SI_AN_KE                  //［四暗刻］：  4坎加一个搭牌，坎牌必须是自己摸的才算，点炮不算 +10嘴
	EXTPXID_HFMJ_3LIAN_KAN                 //［3连坎］：   同花色连在一起的3个坎 +10嘴
	EXTPXID_HFMJ_ANGANG                    //［暗杠］：    必须杠出来 +4嘴
	EXTPXID_HFMJ_SI_HUO                    //［四活］：    某种牌有四张且没有开杠，只要有四张的就算活（豪华七对也算活）。+4嘴
	EXTPXID_HFMJ_MING_SHUANG_PU            //［明双铺］：  两个一样的顺，如345万 345万，若胡牌是胡其中的一张算明双铺，+2嘴 若不是胡其中的牌则为暗双铺+4嘴 自摸+10嘴
	EXTPXID_HFMJ_AN_SHUANG_PU              //［暗双铺］：  两个一样的顺，如345万 345万，若胡牌是胡其中的一张算明双铺，+2嘴 若不是胡其中的牌则为暗双铺+4嘴 自摸+10嘴
	EXTPXID_HFMJ_2SHUANG_PU_ZI             //［双暗双铺］：有2个暗双铺子，双暗双铺子必须是两个暗双铺，自摸算，点炮不算。 +10嘴
)

func GetExtPxName_HFMJ(_pxIndex int) string {
	return []string{"", "支", "卡", "缺门", "同", "10同", "坎", "四暗刻", "3连坎", "暗杠", "四活", "明双铺", "暗双铺", "双暗双铺"}[_pxIndex]
}

//额外牌型分数
func GetExtPxScore_HFMJ(_pxIndex int) int {
	return []int{0, 1, 1, 2, 1, 10, 1, 10, 10, 4, 4, 2, 4, 10, 1, 1}[_pxIndex]
}

// 红中麻将胡牌牌型\分数  ----------------------------------------------------------------------------
const (
	PXID_HZMJ_PINGHU = 1 //平胡
	PXID_HZMJ_DDH    = 2 //对对胡
	PXID_HZMJ_7DUI   = 3 //七对
	PXID_HZMJ_QYS    = 4 //清一色
	PXID_HZMJ_TIANHU = 5 //天胡
	PXID_HZMJ_DIHU   = 6 //地胡
)

func GetHuPxName_HZMJ(_pxIndex int) string {
	return []string{"", "平胡", "对对胡", "七对", "清一色", "天胡", "地胡"}[_pxIndex]
}

func GetHuPxScore_HZMJ(_pxIndex int) int {
	return []int{0, 0, 6, 6, 6, 8, 8}[_pxIndex]
}

//额外加成牌型
const (
	EXTPXID_HZMJ_MENQING = 1 //［门清］：无碰\杠\吃
	EXTPXID_HZMJ_SUPAI   = 2 //［素牌］：无红中
)

func GetExtPxName_HZMJ(_pxIndex int) string {
	return []string{"", "门清", "素牌"}[_pxIndex]
}

//额外牌型分数
func GetExtPxScore_HZMJ(_pxIndex int) int {
	return []int{0, 2, 2}[_pxIndex]
}

// 阜阳麻将胡牌牌型\分数  ----------------------------------------------------------------------------
const (
	PXID_FYMJ_PINGHU         = 1 //平胡
	PXID_FYMJ_GANG_SHANG_HUA = 2 //杠上花
)

func GetHuPxName_FYMJ(_pxIndex int) string {
	return []string{"", "平胡", "杠上花"}[_pxIndex]
}

func GetHuPxScore_FYMJ(_pxIndex int) int {
	return []int{0, 0, 1}[_pxIndex]
}

const (
	PXID_XGMJ_PINGHU  = 1  //平胡
	PXID_XGMJ_HAIDIHU = 30 //海底胡
	PXID_XGMJ_QQR     = 31 //全求人
	PXID_XGMJ_MENQING = 32 //门清

)

const (
	DongFeng = 1 //东
	NanFeng  = 2 //南
	XiFeng   = 3 //西
	BeiFeng  = 4 //北
)

// 怀远麻将胡牌牌型\分数  ----------------------------------------------------------------------------
const (
	PXID_HYMJ_PINGHU     = 1  //平胡
	PXID_HYMJ_TIANHU     = 2  //天胡
	PXID_HYMJ_YINGZI     = 3  //印子（天听）
	PXID_HYMJ_DANDIAO    = 4  //	单吊
	PXID_HYMJ_DUIDUIHU   = 5  //	对对胡
	PXID_HYMJ_QIDUI      = 6  //	七对
	PXID_HYMJ_YIBANGAO   = 7  //一般高
	PXID_HYMJ_HAILAO     = 8  //	海捞
	PXID_HYMJ_SIGUIYI    = 9  //  四归一
	PXID_HYMJ_WUHUAGUO   = 10 // 无花果
	PXID_HYMJ_QINGYISE   = 11 //  清一色
	PXID_HYMJ_YITIAOLONG = 12 //  一条龙
	PXID_HYMJ_QIANGGANG  = 13 // 抢杠
	PXID_HYMJ_GANGKAI    = 14 // 杠开
	PXID_HYMJ_BAOTING    = 15 // 报听

)

//特殊牌型
func GetHuPxName_HYMJ(_pxIndex int) string {
	return []string{"", "平胡", "天胡", "印子", "单吊", "对对胡", "七对", "一般高", "海捞", "四归一", "无花果", "清一色", "一条龙", "抢杠", "杠开", "报听"}[_pxIndex]
}

//特殊牌型分数
func GetHuPxScore_HYMJ(_pxIndex int) int {
	return []int{0, 0, 40, 20, 20, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}[_pxIndex]
}

// TimerType 定时器类型
type TimerType int

const (
	TimerAutoChangeState = 1 //  自动切换牌桌状态
	TimerAutoDismisTable = 2 //  超时自动解散牌桌
)

// SeatState 玩家状态
type SeatState int

const (
	SeatStateNull            = iota //  没有玩家
	SeatStateIdle            = 1    //  牌桌中 空闲,未准备
	SeatStateGameReady       = 2    //  游戏中 点击了准备
	SeatStateGameChangeThree = 3    //  游戏中,已换三张
	SeatStateGameSelectQue   = 4    //  游戏中,已选缺
	SeatStateGameHasHu       = 5    //  游戏中,已胡牌
)

func (t SeatState) String() string {
	return []string{"没有玩家", "空闲,未参与游戏", "游戏中,点击了准备", "游戏中,已换三张", "游戏中,已选缺", "游戏中,已胡牌"}[t]
}

const (
	TableChatTypeExpress = 1 //  牌桌聊天-表情
	TableChatTypeText    = 2 //  牌桌聊天-文字
)

const (
	WorldChatTypeSystem = 1 //  世界喇叭-系统
	WorldChatTypePlayer = 2 //  世界喇叭-玩家
)

const (

	//游戏时间控制 -------------------------------------------------------
	GameTableMaxLiveTime      = 60 * 60 * 3 //桌子最大存贮时间
	GameMaxSavePlayerDataTime = 60 * 30     //最大保存玩家内存数据时间
	GameCheckOffLineTime      = 9           //掉线检测时间,该段时间内未收到消息,认为离线
	GameDismissTableTime      = 60 * 5      //解散房间超时同意时间
	GameWaiteReadyTime        = 15          //等待游戏开始时间
	GameWaiteDealCardTime     = 2           //第一次发牌状态
	GameChangeCardTime        = 10          //换三张时间
	GameSelectQueTime         = 10          //选缺时间
	GameWaiteSendTime         = 15          //等待出牌时间
	GameWaiteThinkTime        = 15          //等待思考时间
	GameZhuaNiaoTime          = 3           //抓鸟时间
	GameShowResultTime        = 2           //结算界面展示时间

	GameWaiteRobotOptPengTime  = 2 //机器人碰牌操作时间
	GameWaiteRobotOptGangTime  = 2 //机器人杠牌操作时间
	GameWaiteRobotOptBUHUATime = 1 //机器人补牌操作时间
	GameWaiteRobotOptHuTime    = 2 //机器人胡牌操作时间
	GameWaiteRobotOptTingTime  = 2 //机器人胡牌操作时间

	//err code ----------------------------------------------------------
	ErrorRoomNotExist             int = -1001 //房间不存在
	ErrorMoneyInsufficient        int = -1002 //金币不足
	ErrorSeatHasPlayer            int = -1003 //座位有人
	ErrorPlayerHasReady           int = -1004 //已准备
	ErrorTableNoPlayer            int = -1006 //桌子没有该玩家
	ErrorTableNotExist            int = -1007 //桌子不存在
	ErrorPlayerNotExist           int = -1008 //玩家不存在
	ErrorChatContentIsNull        int = -1009 //聊天内容不能为空
	ErrorNotLoginInGame           int = -1010 //未登陆游戏
	ErrorTableStateReady          int = -1011 //准备状态错误
	ErrorNotInGame                int = -1012 //不在牌桌游戏中
	ErrorWatchState               int = -1013 //看牌状态错误
	ErrorGiveUpState              int = -1014 //弃牌状态错误,已经放弃
	ErrorNotCurtSpeaker           int = -1015 //没有轮到操作
	ErrorHasInTable               int = -1016 //已经在桌子中
	ErrorChatTypeWrong            int = -1017 //聊天类型错误
	ErrorDiamondInsufficient      int = -1018 //钻石不足
	ErrorDiamondGoodsID           int = -1019 //钻石商品ID错误
	ErrorCoinGoodsID              int = -1020 //金币商品ID错误
	ErrorRankType                 int = -1021 //排行类型错误
	ErrorGiveUpTableState         int = -1022 //弃牌 牌桌状态错误
	ErrorCreateTableType          int = -1023 //创建牌桌类型错误
	ErrorCreateTableSelectParam   int = -1024 //是否型参数错误
	ErrorCreateTableFanCt         int = -1025 //番数错误
	ErrorCreateTableGameCt        int = -1026 //局数错误
	ErrorNotHasThisCard           int = -1027 //没有这张牌
	ErrorSelectQueValue           int = -1028 //选缺参数错误
	ErrorTableStateSend           int = -1029 //出牌状态错误
	ErrorNotCurrentSenderIndex    int = -1030 //不是当前出牌玩家位置
	ErrorTableOptVal              int = -1031 //牌桌操作参数错误
	ErrorTableStateThinkOpt       int = -1032 //操作状态错误
	ErrorNotCurrentThinkIndex     int = -1033 //不是当前思考者
	ErrorWrongTableID             int = -1034 //错误的TableId
	ErrorTableIsFull              int = -1035 //桌子已满
	ErrorHasSaveThinkOpt          int = -1036 //重复操作
	ErrorSendNoThisCard           int = -1037 //出牌错误,手上无此牌
	ErrorSeatOptNoThisOpt         int = -1038 //操作非法,已操作过,或无此操作
	ErrorExitTableGameStart       int = -1039 //游戏已开始,不能退出
	ErrorCreateTableXZWrongParam  int = -1040 //血战到底 ,结算参数选择错误
	ErrorCreateTableDDWrongParam  int = -1041 //血战到底 ,胡牌参数选择错误
	ErrorSendHandHasQueCard       int = -1042 //出牌错误,选缺牌未出完
	ErrorChangeThreeState         int = -1043 //换三张牌桌状态错误
	ErrorSelectQueState           int = -1044 //选缺牌桌状态错误
	ErrorChangeThreeSeatState     int = -1045 //换三张座位状态错误
	ErrorSelectQueSeatState       int = -1046 //选缺座位状态错误
	ErrorDismissTableState        int = -1047 //正在等待解散确认
	ErrorChangeThreeCardColor     int = -1048 //换三张选牌花色错误
	ErrorChangeThreeCardCount     int = -1049 //换三张选牌数量错误
	ErrorDismissNotFangSeatID     int = -1050 //第一局未开始之前只能房主解散房间
	ErrorExitTable                int = -1051 //离桌状态错误
	ErrorCanHuMustHu              int = -1052 //操作非法,最后4张 如果能胡牌则必须胡
	ErrorNotEnoughFangKa          int = -1053 //房卡数量不足
	ErrorSendFangKaFailure        int = -1054 //赠送房卡失败
	ErrorFangKaNumWrong           int = -1055 //房卡数量错误
	ErrorTableRecNotFound         int = -1056 //牌局数据不存在
	ErrorGetZhanJiData            int = -1057 //获取战绩数据失败
	ErrorGetZhanJiDetailData      int = -1058 //获取战绩详细数据失败
	ErrorTableInfoData            int = -1059 //TableInfo 回放数据转换错误
	ErrorCreateTableDianGangParam int = -1060 //点杠花 参数选择错误
	ErrorCreateTableDDHCalType    int = -1061 //倒到胡,结算方式选择错误
	ErrorCreateTableTDQZuiCt      int = -1062 //天地胡\清一色嘴数错误
	ErrorDismissGameNotStart      int = -1063 //游戏未开始,不能解散房间

)

//错误提示键值
var ErrTipArr = [][]string{

	{"-1090", ""},
	{"-1089", ""},
	{"-1088", ""},
	{"-1087", ""},
	{"-1086", ""},
	{"-1085", ""},
	{"-1084", ""},
	{"-1083", ""},
	{"-1082", ""},
	{"-1081", ""},

	{"-1080", ""},
	{"-1079", ""},
	{"-1078", ""},
	{"-1077", ""},
	{"-1076", ""},
	{"-1075", ""},
	{"-1074", ""},
	{"-1073", ""},
	{"-1072", ""},
	{"-1071", ""},

	{"-1070", ""},
	{"-1069", ""},
	{"-1068", ""},
	{"-1067", ""},
	{"-1066", ""},
	{"-1065", ""},
	{"-1064", ""},
	{"-1063", "游戏未开始,不能解散房间"},
	{"-1062", "天地胡清一色嘴数错误"},
	{"-1061", "倒到胡,结算方式选择错误"},

	{"-1060", "点杠花,参数选择错误"},
	{"-1059", "TableInfo回放数据转换错误"},
	{"-1058", "获取战绩详细数据失败"},
	{"-1057", "获取战绩数据失败"},
	{"-1056", "牌局数据不存在"},
	{"-1055", "房卡数量错误"},
	{"-1054", "赠送房卡失败"},
	{"-1053", "房卡数量不足"},
	{"-1052", "操作非法,最后4张 如果能胡牌则必须胡"},
	{"-1051", "离桌状态错误"},
	{"-1050", "第一局未开始之前只能房主解散房间"},
	{"-1049", "换三张选牌数量错误"},
	{"-1048", "换三张选牌花色错误"},
	{"-1047", "正在等待解散确认"},
	{"-1046", "选缺座位状态错误"},
	{"-1045", "换三张座位状态错误"},
	{"-1044", "选缺牌桌状态错误"},
	{"-1043", "换三张牌桌状态错误"},
	{"-1042", "出牌错误,选缺牌未出完"},
	{"-1041", "血战到底 ,胡牌参数选择错误"},

	{"-1040", "血战到底 ,结算参数选择错误"},
	{"-1039", "游戏已开始,不能退出"},
	{"-1038", "已操作过,或无此操作"},
	{"-1037", "出牌错误,手上无此牌"},
	{"-1036", "重复操作"},
	{"-1035", "桌子已满"},
	{"-1034", "错误的TableId"},
	{"-1033", "不是当前思考者"},
	{"-1032", "操作状态错误"},
	{"-1031", "牌桌操作参数错误"},
	{"-1030", "不是当前出牌玩家位置"},
	{"-1029", "出牌状态错误"},
	{"-1028", "选缺参数错误"},
	{"-1027", "没有这张牌"},
	{"-1026", "局数错误"},
	{"-1025", "番数错误"},
	{"-1024", "是否型参数错误"},
	{"-1023", "创建牌桌类型错误"},
	{"-1022", "弃牌牌桌状态错误"},
	{"-1021", "排行类型错误"},

	{"-1020", "金币商品ID错误"},
	{"-1019", "钻石商品ID错误"},
	{"-1018", "钻石不足"},
	{"-1017", "聊天类型错误"},
	{"-1016", "已经在桌子中"},
	{"-1015", "没有轮到操作"},
	{"-1014", "弃牌状态错误,已经放弃"},
	{"-1013", "看牌状态错误"},
	{"-1012", "不在牌桌游戏中"},
	{"-1011", "准备状态错误"},
	{"-1010", "未登陆游戏"},
	{"-1009", "聊天内容不能为空"},
	{"-1008", "玩家不存在"},
	{"-1007", "桌子不存在"},
	{"-1006", "桌子没有该玩家"},
	{"-1005", ""},
	{"-1004", "已准备"},
	{"-1003", "座位有人"},
	{"-1002", "金币不足"},
	{"-1001", "房间不存在"},
	{"0", "success"},
}

//返回错误提示
func GetErrTip(_errCode int) string {
	for _, v := range ErrTipArr {
		if strconv.Itoa(_errCode) == v[0] {
			return v[1]
		}
	}
	return "unknow error"
}
