package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")
var ErrorUserIdInvalid = errors.New("用户ID格式错误")

// 这个不能放到middleware中的auth.go中，会造成循环引用的问题
const ContextUserIdKey = "UserId"

// 从gin的上下文中取用户id
func getCurrentUserId(ctx *gin.Context) (int64, error) {
	val, ok := ctx.Get(ContextUserIdKey)
	if !ok {
		return -1, ErrorUserNotLogin
	}

	uid, ok := val.(int64)
	if !ok {
		return -1, ErrorUserIdInvalid
	}
	return uid, nil
}

// 从参数获取要查询的条数和
// func getBlogSizeAndPage(ctx *gin.Context) (page int64, size int64, err error) {
// 	pageStr := ctx.Query("page")
// 	sizeStr := ctx.Query("size")

// 	page, err = strconv.ParseInt(pageStr, 10, 64)
// 	if err != nil {
// 		page = 1
// 	}

// 	size, err = strconv.ParseInt(sizeStr, 10, 64)
// 	if err != nil {
// 		size = 10
// 	}
// 	err = nil
// 	return
// }
