package controller

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/logic"
	"NothingBlog/models"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpHandler 用户注册接口
// @Summary 用户注册的接口
// @Description 用户注册的接口，需要接收参数
// @Tags 认证相关接口
// @Accept application/json
// @Produce application/json
// @Security No
// @Param object body models.SignUpParams true "注册参数"
// @Success 200 {object} ResponseData
// @Router /auth/signup [post]
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

// LoginHandler 用户登录接口
// @Summary 用户登录接口
// @Description 用于用户登录
// @Tags 认证相关接口
// @Accept application/json
// @Produce application/json
// @Security No
// @Param object body models.LoginParams true "登录参数"
// @Success 200 {object} ResponseData
// @Router /auth/login [post]
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

// PwdResetHandler 重置密码接口
// @Summary 重置密码接口
// @Description 用于用户重置密码，需要使用邮箱 - 功能未实现
// @Tags 认证相关接口
// @Accept application/json
// @Produce application/json
// @Security No
// @Param object body models.ResetPasswordParams true "重置密码的参数"
// @Success 200 {object} ResponseData
// @Router /auth/password/reset [post]
func PwdResetHandler(ctx *gin.Context) {
	// 获取前端传递的账号信息
	var rstPara = new(models.ResetPasswordParams)
	if err := ctx.ShouldBindJSON(rstPara); err != nil {
		zap.L().Error("解析用户重置密码参数失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	// 逻辑层处理
	_ = logic.ResetPassword(rstPara)
	// 校验账号信息

	// 获取邮箱并发送重置验证码

	// 缓存验证码

	// 重置用户密码为 go_ngblog

	// 返回信息告知
}

// PwdModifyHandler 修改用户密码
// @Summary 修改用户密码
// @Description 用于修改用户密码，不需要登录，需要验证旧密码，账户和验证码
// @Tags 认证相关接口
// @Accept application/json
// @Produce application/json
// @Security No
// @Param object body models.ModifyPasswordParams true "修改密码的参数"
// @Success 200 {object} ResponseData
// @Router /auth/password/modify [post]
func PwdModifyHandler(ctx *gin.Context) {
	// 获取密码参数
	var rstPara = new(models.ModifyPasswordParams)
	if err := ctx.ShouldBindJSON(rstPara); err != nil {
		zap.L().Error("解析用户修改密码参数失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	// 校验验证码
	if err := logic.ModifyPassword(rstPara); err != nil {
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
	// 返回结果
	ResponseSuccessWithMsg(ctx, "修改密码成功", nil)
}
