package controller

type ResponseCodeType uint32

const (
	CodeSuccess ResponseCodeType = 1000 + iota
	CodeParameterInvalid
	CodeUserNotExist
	CodeUserExist
	CodePasswordError
	CodeServerBusy

	// jwt token 相关状态定义
	CodeTokenInvaild
	CodeTokenEmpty

	// 重新登录
	CodeNeedReLogin

	// 社区
	CodeCommunityIdInvalid

	// 验证码错误
	CodeVerifyCodeInvaild
)

var codeMsgMap = map[ResponseCodeType]string{
	CodeSuccess:            "success",
	CodeParameterInvalid:   "参数错误",
	CodeUserNotExist:       "用户名不存在",
	CodeUserExist:          "用户名已存在",
	CodePasswordError:      "用户名或密码错误",
	CodeServerBusy:         "服务器繁忙",
	CodeTokenInvaild:       "Token 无效",
	CodeTokenEmpty:         "Token 为空",
	CodeNeedReLogin:        "登录信息无效，重新登录",
	CodeCommunityIdInvalid: "社区ID无效",
	CodeVerifyCodeInvaild:  "验证码错误",
}

func (rct ResponseCodeType) Msg() string {
	msg, ok := codeMsgMap[rct]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
