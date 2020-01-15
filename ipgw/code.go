package ipgw

var (
	codeMap = map[int]string{
		3: "与ipgw版本不兼容，请升级",
		4: "读取本地配置失败",
		5: "网络错误，请检查网络",

		11: "未指定密码",
		12: "无已保存账号，请指定账号密码",
		13: "账号或密码错误",
		14: "Cookie已失效",
		15: "账户被禁或服务未授权",
		16: "登陆失败，请重试",
	}
)

func CodeText(code int) string {
	return codeMap[code]
}
