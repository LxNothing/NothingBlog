package controller

import (
	"NothingBlog/models"
	"NothingBlog/settings"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetIndexHandler(ctx *gin.Context) {
	// 获取参数
	pageStr := ctx.Param("page")
	var page int64 = 1
	var err error
	if pageStr != "" {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			zap.L().Debug("页数转换为int失败", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}
	}
	// page_size := int(settings.Confg.PageSize)
	// stus, _ := models.StatusStringToNumber(models.Drift)
	// priv, _ := models.PrivilegeStringToNumber(models.Private)
	// atc, total, err := logicArticle.GetArticleWithPage(int(page), page_size, stus, priv, 0)

	param := &models.ArticleWithPageParams{
		Keyword:   "",
		Privilege: models.Public,
		Status:    models.Commit,
		Page:      uint(page),
		Size:      uint(settings.Confg.PageSize),
		Tag:       "",
		Class:     "",
	}

	atc, cur_page, total_page, err := logicArticle.GetAllWithParams(param)
	if err != nil {
		zap.L().Debug("分页查询文章出错", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 按页返回文章
	ResponseSuccess(ctx, gin.H{
		"article":    atc,
		"cur_page":   cur_page,   // 当前返回的页数
		"total_page": total_page, // 总共有多少页
	})
}
