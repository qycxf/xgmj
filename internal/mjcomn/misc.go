package mjcomn

//麻将配置对象 ------------------------------------------
type MjCfg struct {
	TableType   int   //牌桌类型
	MaxColorCt  int   //所使用的麻将牌花色最大数(最大6)
	Check7Dui   bool  //是否检测七对
	LaiziData   int   //癞子牌值(如有)
	LZExceptArr []int //癞子不可代替的牌
}

func NewMjCfg() MjCfg {
	ret := MjCfg{
		TableType:   0,
		MaxColorCt:  DefaultIndex,
		Check7Dui:   false,
		LaiziData:   DefaultIndex,
		LZExceptArr: make([]int, 0),
	}
	return ret
}

//出牌提示数据对象 ------------------------------------------
type SendTip struct {
	SendCard int         //打出的牌
	HuCards  []int       //如果打出这张 对应能胡的牌
	HuInfos  []CalHuInfo //对应能胡的牌 胡牌型信息
}

func NewSendTip() *SendTip {
	ret := SendTip{
		SendCard: DefaultIndex,
		HuCards:  make([]int, 0),
		HuInfos:  make([]CalHuInfo, 0),
	}
	return &ret
}

func IntSliceCopy(src []int) []int {
	ret := make([]int, len(src))
	copy(ret, src)
	return ret
}
