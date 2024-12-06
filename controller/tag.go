package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"
	"strconv"

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

func GetTagByIdHandler(ctx *gin.Context) {
	tagId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		zap.L().Debug("Tag的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	tag, err := logic.GetTagById(tagId)
	if err != nil {
		zap.L().Debug("查询文章种类失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
	}
	ResponseSuccess(ctx, tag)
}

// 删除单个tag
func DeleteTagHandler(ctx *gin.Context) {
	tagId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		zap.L().Debug("Tag的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	err = logic.DeleteTagById(tagId)
	if err != nil {
		zap.L().Error("删除tag失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "删除Tag成功", nil)
}

func DeleteMultiTagHandler(ctx *gin.Context) {
	param := new(models.DeleteMultiTagParams)
	if err := ctx.ShouldBindJSON(param); err != nil {
		zap.L().Debug("Tag的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	if err := logic.DeleteMultiTagById(param.Ids); err != nil {
		zap.L().Error("删除Tag失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "删除Tag成功", nil)
}
