package controller

import (
	"NothingBlog/logic"
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
	page_size := int(settings.Confg.PageSize)
	atc, total, err := logic.GetArticleWithPage(int(page), page_size, models.StatusDraft, models.PrivilegePrivte, 0)

	if err != nil {
		zap.L().Debug("分页查询文章出错", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 计算总的页数
	total_page := total / int64(page_size)
	if total%int64(page_size) != 0 {
		total_page++
	}

	// 按页返回文章
	ResponseSuccess(ctx, gin.H{
		"article":    atc,
		"cur_page":   page,       // 当前返回的页数
		"total_page": total_page, // 总共有多少页
	})
}
