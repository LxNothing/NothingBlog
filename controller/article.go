package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAllArticleHandler 获取所有文章
// @Summary 获取所有文章的接口
// @Description 通过该接口可以获得当前的所有文章
// @Tags 文章相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData
// @Router /articles [get]
func GetAllArticleHandler(ctx *gin.Context) {
	aticles, err := logic.GetAllArticle()

	if errors.Is(err, logic.ErrArticleNotExisted) {
		zap.L().Debug("查询的文章不存在", zap.Error(err))
		ResponseError(ctx, CodeArticleNotExisted)
		return
	}

	if errors.Is(err, logic.ErrArticleQueryFailed) {
		zap.L().Debug("数据库查询出错", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccess(ctx, aticles)
}

// GetArticleWithIdHandler 通过文章ID查询文章详细信息
// @Summary 通过文章ID查询文章详细信息的接口
// @Description 通过该接口可以查询文章详细信息
// @Tags 文章相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData "code=200表示成功其余失败"
// @Router /article/:id [get]
func GetArticleWithIdHandler(ctx *gin.Context) {
	// 获取参数
	atcId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		zap.L().Debug("查询文章的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 调用逻辑层方法
	article, err := logic.GetArticleById(atcId)
	if errors.Is(err, logic.ErrArticleNotExisted) {
		zap.L().Debug("查询的文章不存在", zap.Error(err))
		ResponseError(ctx, CodeArticleNotExisted)
		return
	}

	if errors.Is(err, logic.ErrArticleQueryFailed) {
		zap.L().Debug("数据库查询出错", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 更新文章的访问量 - 如果未登录进行访问则代表是客户端访问，admin访问不更新
	if uid, _ := getCurrentUserId(ctx); uid == -1 {
		_ = logic.UpdateArticleVisitCount(atcId, article.VisitCount+1)
	}

	ResponseSuccess(ctx, article)
}

// 根据Tag获取文章列表
func GetArticleWithTagHandler(ctx *gin.Context) {

}

// 根据文章分类获取文章列表
func GetArticleWithClassHandler(ctx *gin.Context) {

}

// 获取前一篇文章
func GetPreArticleHandler(ctx *gin.Context) {

}

// 获取后一篇文章
func GetAfterArticleHandler(ctx *gin.Context) {

}

// 分页获取文章
func GetArticleWithPageHandler(ctx *gin.Context) {

}

// CreateArticleHandler 创建文章
// @Summary 创建文章的接口
// @Description 通过该接口可以创建文章
// @Tags 文章相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.NewArticleFormsParams true "创建文章的参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData "code=200表示成功其余失败"
// @Router /article [post]
func CreateArticleHandler(ctx *gin.Context) {
	newAtcParams := new(models.NewArticleFormsParams)
	if err := ctx.ShouldBindJSON(newAtcParams); err != nil {
		zap.L().Error("解析创建文章的新参数出错", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	// uid, err := getCurrentUserId(ctx)
	// if err != nil {
	// 	zap.L().Error("获取用户id失败", zap.Error(err))
	// 	ResponseError(ctx, CodeNeedReLogin)
	// 	return
	// }
	uid := int64(4027674162892800) //for test

	article := &models.Article{
		AuthorId:  uid,
		ClassId:   newAtcParams.ClassId,
		TopFlag:   newAtcParams.TopFlag,
		EnComment: newAtcParams.EnComment,
		Status:    newAtcParams.Status,
		Privilege: newAtcParams.Privilege,
		Title:     newAtcParams.Title,
		Content:   newAtcParams.Content,
		Summary:   newAtcParams.Summary,
		Image:     newAtcParams.Image,
	}

	// 调用logic层的方法创建文章
	if err := logic.CreateNewArticle(article, newAtcParams.TagIdList); err != nil {
		zap.L().Error("创建文章失败", zap.Error(err))
		if errors.Is(err, logic.ErrArticleNameExisted) {
			ResponseError(ctx, CodeArticleTitleExisted)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccessWithMsg(ctx, "创建文章成功", nil)
}

// DeleteMultiArticleHandler 删除单篇文章
// @Summary 删除单篇文章的接口
// @Description 通过该接口可以删除指定的文章
// @Tags 文章相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData "code=200表示成功其余失败"
// @Router /article/:id [delete]
func DeleteArticleHandler(ctx *gin.Context) {
	atcId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		zap.L().Debug("查询文章的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	err = logic.DeleteArticleById(atcId)
	if err != nil {
		zap.L().Error("删除文章失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "删除文章成功", nil)
}

// DeleteMultiArticleHandler 删除多篇文章 - 硬删除
// @Summary 删除多篇文章的接口 - 硬删除，会直接删除数据库中的数据
// @Description 通过该接口可以删除指定的文章
// @Tags 文章相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.DeleteMultiArticleParams true "删除文章的参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData "code=200表示成功其余失败"
// @Router /articles [delete]
func DeleteMultiArticleHandler(ctx *gin.Context) {
	param := new(models.DeleteMultiArticleParams)
	if err := ctx.ShouldBindJSON(param); err != nil {
		zap.L().Debug("删除文章的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	if err := logic.DeleteMultiArticleById(param.Ids); err != nil {
		zap.L().Debug("删除文章失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "删除文章成功", nil)
}

// UpdateArticleStatusHandler 删除多篇文章 - 软删除 或者恢复
// @Summary 删除时,只会更新其状态为删除状态,不会删除数据库中的数据，恢复时将其恢复为指定状态
// @Description 通过该接口可以删除指定的文章 - 软删除
// @Tags 文章相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.SoftDeleteArticleParams true "删除文章的参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData "code=200表示成功,其余失败"
// @Router /articles [delete]
func UpdateArticleStatusHandler(ctx *gin.Context) {
	param := new(models.SoftDeleteArticleParams)
	if err := ctx.ShouldBindJSON(param); err != nil {
		zap.L().Debug("软删除文章的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	if err := logic.UpdateArticleStatusById(param.Ids, param.DelFalg); err != nil {
		zap.L().Debug("软删除文章失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "软删除文章成功", nil)
}

// UpdateArticleHandler 更新文章
// @Summary 更新文章的接口
// @Description 通过该接口可以更新指定的文章
// @Tags 文章相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.UpdateArticleFormsParams true "更新文章的参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData "code=200表示成功其余失败"
// @Router /article [put]
func UpdateArticleHandler(ctx *gin.Context) {
	updateAtcParams := new(models.UpdateArticleFormsParams)
	if err := ctx.ShouldBindJSON(updateAtcParams); err != nil {
		zap.L().Debug("解析更新文章的新参数出错", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	uid, err := getCurrentUserId(ctx)
	if err != nil {
		zap.L().Debug("获取用户id失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	article := &models.Article{
		AuthorId:  uid,
		ClassId:   updateAtcParams.ClassId,
		TopFlag:   updateAtcParams.TopFlag,
		EnComment: updateAtcParams.EnComment,
		Status:    updateAtcParams.Status,
		Privilege: updateAtcParams.Privilege,
		Title:     updateAtcParams.Title,
		Content:   updateAtcParams.Content,
		Summary:   updateAtcParams.Summary,
		Image:     updateAtcParams.Image,
	}

	// 调用logic层的方法更新文章
	if err := logic.UpdateArticle(updateAtcParams.ArticleId, article, updateAtcParams.TagIdList); err != nil {
		zap.L().Debug("更新文章失败", zap.Error(err))
		if errors.Is(err, logic.ErrArticleNameExisted) {
			ResponseError(ctx, CodeArticleTitleExisted)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccessWithMsg(ctx, "更新文章成功", nil)
}
