package errMsg

var MessageMap map[int]string

func init() {
	MessageMap = make(map[int]string, 0)
	MessageMap[100] = "操作成功"
	MessageMap[101] = "签名错误"
	MessageMap[102] = "参数错误"
	MessageMap[103] = "系统繁忙"
	MessageMap[104] = "数据存储出错"
	MessageMap[105] = "未支付"
	MessageMap[106] = "获取超时"
	MessageMap[107] = "用户名或密码错误"
	MessageMap[108] = "设备不在线"
	MessageMap[109] = "操作失败"
	MessageMap[110] = "请求超时"
	MessageMap[111] = "数据格式错误"
	MessageMap[112] = "任务不存在"
	MessageMap[113] = "认证失败"
}

func GetErrMessage(code int) string {
	if value, ok := MessageMap[code]; ok {
		return value
	}
	return ""
}
