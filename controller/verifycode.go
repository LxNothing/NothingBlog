package controller

import (
	"NothingBlog/package/verifycode"
	"NothingBlog/settings"

	"github.com/gin-gonic/gin"
)

func VerifyCodeHandler(ctx *gin.Context) {
	id, code, _ := verifycode.GenerateVerifyCode()
	// if err != nil {
	// 	zap.L().Error("生成验证码出错", zap.Error(err))
	// 	ResponseError(ctx, CodeServerBusy)
	// 	return
	// }
	ResponseSuccess(ctx, gin.H{
		"id":      id,
		"code":    code,
		"expired": settings.Confg.ExpiredTime, // 验证码的过期时间
	})
}
