package mjcomn

const (
	Success int = 0  // Success 成功
	Failure     = -1 // Failure 失败

	DefaultIndex int = 999 //默认下标值,用该值替代-1
)

// 牌桌分类
const (
	TableClass_FangKa = 1 // 房卡
	TableClass_Coin   = 2 // 金币
)

//返回牌桌分类名称
func GetTableClassName(_ttype int) string {
	return []string{"", "房卡", "金币"}[_ttype]

}

// 牌桌类型
const (
	TableType_HFMJ       = 1  // 合肥麻将
	TableType_HZMJ       = 2  // 红中麻将
	TableType_YC_XLCH    = 3  // 宜昌血流成河
	TableType_XY_KA5XING = 4  // 襄阳卡五星
	TableType_SC_XZDD    = 5  // 四川血战到底
	TableType_FJ_FZMJ    = 6  // 福建福州麻将
	TableType_NJMJ       = 7  // 南京麻将
	TableType_FYMJ       = 8  // 阜阳麻将
	TableType_ASMJ       = 9  // 鞍山麻将
	TableType_GBMJ       = 10 // 国标麻将
	TableType_AQMJ       = 11 // 安庆麻将
	TableType_BBMJ       = 12 // 蚌埠麻将
	TableType_HYMJ       = 13 // 怀远麻将
)

//返回牌桌名称
func GetTableName(_ttype int) string {
	return []string{"", "合肥麻将", "红中麻将", "宜昌血流成河", "襄阳卡五星", "四川血战到底", "福州麻将",
		"南京麻将", "阜阳麻将", "鞍山麻将", "国标麻将", "安庆麻将", "蚌埠麻将", "怀远麻将"}[_ttype]

}

//麻将花色
const (
	Color_Wan  int = 0 // 万
	Color_Tong int = 1 // 筒
	Color_Tiao int = 2 // 条
	Color_Feng int = 3 // 东南西北风
	Color_Zfb  int = 4 // 中发白
	Color_Hua  int = 5 // 花
)

//玩家操作类型
const (
	OptTypePeng   int = 1 //  碰
	OptTypeChi    int = 2 //  吃
	OptTypeGang   int = 3 //  杠
	OptTypeHu     int = 4 //  胡
	OptTypeBu     int = 5 //  补张
	OptTypeFetch  int = 6 //  拿牌
	OptTypeSend   int = 7 //  出牌
	OptTypeCancel int = 8 //  取消
	OptTypeTing   int = 9 //  听牌
)

//胡牌类型
const (
	PXTYPE_UNKNOW = 0 //未知
	PXTYPE_PINGHU = 1 //平胡
	PXTYPE_7DUI   = 2 //七对
	PXTYPE_SSY    = 3 //十三幺胡牌
	PXTYPE_SSL    = 4 //十三烂胡牌
	PXTYPE_TIANHU = 3 //天胡
)

//返回胡牌类型名称
func GetPxTypeName(_opt int) string {
	return []string{"未知", "平胡", "七对", "天胡"}[_opt]

}

//返回操作名称
func GetOptName(_opt int) string {
	arr := []string{"", "碰牌", "吃牌", "杠牌", "胡牌", "补张", "拿牌", "出牌", "取消", "听牌"}
	return arr[_opt]

}

//碰牌\吃牌\杠牌 用于展示手牌[碰杠吃]牌部分
const (
	PGCTypePeng     int = 1 //  碰
	PGCTypeChi      int = 2 //  吃
	PGCTypeMingGang int = 3 //  明杠
	PGCTypeAnGang   int = 4 //  暗杠
)

// 杠牌类型
const (
	GANGTYPE_AN   = 1 //暗杠
	GANGTYPE_ZHI  = 2 //直杠(明杠)
	GANGTYPE_MIAN = 3 //面杠(明杠)
)

//胡牌类型
const (
	HUTYPE_JIEPAO = 1 //接炮
	HUTYPE_ZIMO   = 2 //自摸
)

func GetHuTypeName(index int) string {
	return []string{"", "接炮", "自摸"}[index]
}

//胡牌详细类型
const (
	HUTYPE_DETAIL_COMMN          = 1 //正常胡牌(接炮)
	HUTYPE_DETAIL_QIANGGANG      = 2 //抢杠胡
	HUTYPE_DETAIL_GANG_SHANG_PAO = 3 //杠上炮
	HUTYPE_DETAIL_GHH_AG         = 4 //杠后花 - 暗杠
	HUTYPE_DETAIL_GHH_MG         = 5 //杠后花 - 面杠
	HUTYPE_DETAIL_GHH_ZG         = 6 //杠后花 - 直杠
	HUTYPE_DETAIL_GANG_SHANG_HUA = 7 //杠上花 - 杠、补花都算
)

func GetHuDetailName(index int) string {
	return []string{"", "正常胡", "抢杠胡", "杠上炮", "杠后花(暗杠)", "杠后花(面杠)", "杠后花(直杠)", "杠上花"}[index]
}

func GetGangTypeName(index int) string {
	return []string{"", "暗杠", "直杠", "面杠"}[index]
}
