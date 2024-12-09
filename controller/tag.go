package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateTagHandler 创建tag
// @Summary 创建标签(tag)的接口
// @Description 通过该接口可以创建标签
// @Tags 标签相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.TagCreateFormParams true "创建标签的参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /tag [post]
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

// GetAllTagsHandler 获取所有的标签
// @Summary 获取所有标签（简略信息）的接口
// @Description 通过该接口可以获得当前的所有标签
// @Tags 标签相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseAllTagList
// @Router /tags [get]
func GetAllTagsHandler(ctx *gin.Context) {
	tag, err := logic.GetAllTags()
	if err != nil {
		zap.L().Debug("查询文章种类失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
	}
	ResponseSuccess(ctx, tag)
}

// GetTagByIdHandler 根据ID获取标签的详细信息
// @Summary 根据ID获取标签的详细信息的接口
// @Description 通过该接口可以获得当前标签的详细信息
// @Tags 标签相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseTagDetailList
// @Router /tag/:id [get]
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

	resp := models.ResponseTagDetail{
		ResponseTagBrief: models.ResponseTagBrief{
			TagId:    tag.TagId,
			Name:     tag.Name,
			AtcCount: int32(tag.ArticleCount),
		},
		Desc:      tag.Desc,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}

	ResponseSuccess(ctx, resp)
}

// DeleteTagHandler 删除单个Tag
// @Summary 删除单个标签(Tag)的接口
// @Description 通过该接口可以删除单个Tag
// @Tags 标签相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /tag/:id [delete]
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

// DeleteMultiTagHandler 删除多个Tag
// @Summary 删除多个标签(Tag)的接口
// @Description 通过该接口可以删除多个Tag
// @Tags 标签相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.DeleteMultiTagParams true "删除多个标签的参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /tag [delete]
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

// UpdateTagHandler 更新Tag
// @Summary 更新标签(Tag)的接口
// @Description 通过该接口可以更新Tag
// @Tags 标签相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.UpdateTagParams true "更新标签的参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /tag [put]
func UpdateTagHandler(ctx *gin.Context) {
	newTag := new(models.UpdateTagParams)

	if err := ctx.ShouldBindJSON(newTag); err != nil {
		zap.L().Debug("解析修改Tag的参数错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// // 查找ID是否存在
	_, err := logic.GetTagById(newTag.TagId)
	if err != nil {
		zap.L().Debug("ID错误", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeParameterInvalid, "要修改的Tag不存在")
		return
	}

	// 查找名称是否重复
	oldTag, err := logic.GetTagByName(newTag.Name)
	if err == nil && oldTag.TagId > 0 && oldTag.TagId != newTag.TagId {
		zap.L().Debug("Tag名称重复")
		ResponseError(ctx, CodeTagExisted) // 名称重复
		return
	}

	var tag = &models.Tag{
		TagId: newTag.TagId,
		Name:  newTag.Name,
		Desc:  newTag.Desc,
	}

	if err := logic.UpdateTag(tag); err != nil {
		zap.L().Debug("更新Tag失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "修改Tag成功", nil)
}
