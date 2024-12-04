package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateTagHandler(ctx *gin.Context) {
	tagParam := new(models.TagCreateFormParams)
	if err := ctx.ShouldBindJSON(tagParam); err != nil {
		zap.L().Debug("解析创建Tag的参数失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	tag := new(models.Tag)
	tag.Name = tagParam.Name
	tag.Desc = tagParam.Desc

	if err := logic.CreateNewTag(tag); err != nil {
		zap.L().Debug("创建Tag失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "创建Tag成功", nil)
}

func GetAllTagsHandler(ctx *gin.Context) {
	tag, err := logic.GetAllTags()
	if err != nil {
		zap.L().Debug("查询文章种类失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
	}
	ResponseSuccess(ctx, tag)
}
