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
// @Router /article/all [get]
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

// 根据文章ID获取文章
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

// 根据Tag获取文章
func GetArticleWithTagHandler(ctx *gin.Context) {

}

// 根据文章分类获取文章
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

// 创建文章
func CreateArticleHandler(ctx *gin.Context) {
	newAtcParams := new(models.NewArticleFormsParams)
	if err := ctx.ShouldBindJSON(newAtcParams); err != nil {
		zap.L().Error("解析创建文章的新参数出错", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	uid, err := getCurrentUserId(ctx)
	if err != nil {
		zap.L().Error("获取用户id失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

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

// 文章更新

// 更新文章访问量

// 删除单篇文章
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

func DeleteMultiArticleHandler(ctx *gin.Context) {
	param := new(models.DeleteMultiArticleParams)
	if err := ctx.ShouldBindJSON(param); err != nil {
		zap.L().Debug("查询文章的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	if err := logic.DeleteMultiArticleById(param.Ids); err != nil {
		zap.L().Error("删除文章失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "删除文章成功", nil)
}
