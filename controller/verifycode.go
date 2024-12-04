package controller

import (
	"NothingBlog/package/verifycode"
	"NothingBlog/settings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VerifyCodeHandler 获取验证码接口
// @Summary 获取数字验证码的接口
// @Description 通过该接口可以获得基于数字的验证码，目前只支持数字验证码，后续可以更改
// @Tags 认证相关接口
// @Accept application/json
// @Produce application/json
// @Security No
// @Success 200 {object} ResponseData
// @Router /auth/verifycode [get]
func VerifyCodeHandler(ctx *gin.Context) {
	id, code, _, err := verifycode.GenerateVerifyCode()
	if err != nil {
		zap.L().Error("生成验证码出错", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, gin.H{
		"id":      id,
		"code":    code,
		"expired": settings.Confg.ExpiredTime, // 验证码的过期时间
	})
}
