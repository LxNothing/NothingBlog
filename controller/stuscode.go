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

	// 文章
	CodeArticleTitleExisted
	CodeArticleNotExisted
	CodeHaveArticleInClass

	// tag
	CodeTagExisted
	CodeTagNotExisted

	// class
	CodeClassNotExisted
	CodeClassNameExisted
)

var codeMsgMap = map[ResponseCodeType]string{
	CodeSuccess:             "success",
	CodeParameterInvalid:    "参数错误",
	CodeUserNotExist:        "用户名不存在",
	CodeUserExist:           "用户名已存在",
	CodePasswordError:       "用户名或密码错误",
	CodeServerBusy:          "服务器繁忙",
	CodeTokenInvaild:        "Token 无效",
	CodeTokenEmpty:          "Token 为空",
	CodeNeedReLogin:         "登录信息无效，重新登录",
	CodeCommunityIdInvalid:  "社区ID无效",
	CodeVerifyCodeInvaild:   "验证码错误",
	CodeArticleTitleExisted: "文章名重复",
	CodeArticleNotExisted:   "文章不存在",
	CodeHaveArticleInClass:  "类别下存在文章",
	CodeTagExisted:          "Tag名称重复",
	CodeTagNotExisted:       "Tag不存在",
	CodeClassNotExisted:     "Class不存在",
	CodeClassNameExisted:    "Class名称重复",
}

func (rct ResponseCodeType) Msg() string {
	msg, ok := codeMsgMap[rct]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
