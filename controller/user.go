package controller

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/logic"
	"NothingBlog/models"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 处理用户注册的请求函数
func SignUpHandler(ctx *gin.Context) {
	// 获取参数和参数校验
	u := new(models.SignUpParams)
	if err := ctx.ShouldBindJSON(u); err != nil {
		zap.L().Error("解析用户注册信息失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 业务处理 - 位于logic层
	if err := logic.UserSignup(u); err != nil {
		zap.L().Error("用户参数校验成功，但是注册失败", zap.Error(err))
		// 用户已存在
		if errors.Is(err, mysql.ErrUserExisted) {
			ResponseError(ctx, CodeUserExist)
			return
		}
		// 验证码错误
		if errors.Is(err, logic.ErrVerifyCode) {
			ResponseError(ctx, CodeVerifyCodeInvaild)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回请求
	ResponseSuccessWithMsg(ctx, "注册成功", nil)
}

// 处理用户登录的请求函数
func LoginHandler(ctx *gin.Context) {
	// 获取参数和参数校验
	u := new(models.LoginParams)
	if err := ctx.ShouldBindJSON(u); err != nil {
		zap.L().Error("解析用户登录数据失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 业务处理 - 位于logic层
	token, err := logic.UserLogin(u)
	if err != nil {
		zap.L().Error("用户参数校验成功，但是登录失败", zap.Error(err))
		// 验证码错误
		if errors.Is(err, logic.ErrVerifyCode) {
			ResponseError(ctx, CodeVerifyCodeInvaild)
			return
		}
		// 用户不存在
		if errors.Is(err, logic.ErrUserNotExisted) {
			ResponseError(ctx, CodeUserNotExist)
			return
		}
		// 密码错误
		if errors.Is(err, logic.ErrUserPassword) {
			ResponseError(ctx, CodePasswordError)
			return
		}
		// 其他错误返回服务器繁忙 - 隐藏真正的后端错误
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回请求 - 并且附带token
	ResponseSuccessWithMsg(ctx, "登录成功", gin.H{
		"jwt": token,
	})
}
