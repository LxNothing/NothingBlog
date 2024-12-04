package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateClassHandler(ctx *gin.Context) {
	classParam := new(models.ClassCreateFormParams)
	if err := ctx.ShouldBindJSON(classParam); err != nil {
		zap.L().Debug("解析创建Tag的参数失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	cls := new(models.Class)
	cls.Name = classParam.Name
	cls.Desc = classParam.Desc

	if err := logic.CreateNewClass(cls); err != nil {
		zap.L().Debug("创建Tag失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "创建Class成功", nil)
}
