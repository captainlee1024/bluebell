package controller

type ResCode int64

// 定义状态码
const (
	CodeSuccess         ResCode = 1000 + iota
	CodeInvalidParam            // 请求参数有误
	CodeUserExist               // 用户已存在
	CodeUserNotExist            // 用户不存在
	CodeInvalidPassword         // 密码错误
	CodeServerBusy              // 服务器繁忙，例如数据库连接错误的时候，不需要吧具体的信息返回给前端用户

	CodeNeedLogin
	CodeInvalidToken

	CodeLoginElsewhere
)

// 定义状态码对应信息
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数有误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeNeedLogin:    "请登录",
	CodeInvalidToken: "无效的token",

	CodeLoginElsewhere: "账号已在其它客户端登录，重新登录",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
