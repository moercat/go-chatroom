package enum

type Area string

const (
	PublicScreen Area = "public_screen" // 公屏
	GroupArea    Area = "group_chat"    //群聊
	PrivateArea  Area = "private_chat"  //私聊
)

var areaMap = map[Area]string{
	PublicScreen: "【公屏】",
	GroupArea:    "【群聊】",
	PrivateArea:  "【私聊】",
}

// AreaCheck 检查是否是合法区域，默认在公屏聊天
func (a Area) AreaCheck() string {
	if area, exist := areaMap[a]; exist {
		return area
	}
	return areaMap[PublicScreen]
}
