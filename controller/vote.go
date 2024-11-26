package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VotedHandler(ctx *gin.Context) {
	// 获取参数
	voteDate := new(models.VoteDateParams)
	if err := ctx.ShouldBindJSON(voteDate); err != nil {
		zap.L().Error("获取投票参数失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	// 参数校验 - 已经使用validater这个库做了字段和传递的值的判断 - 还需要判断传递文章id是否存在

	// 获取投票的用户id
	uid, err := getCurrentUserId(ctx)
	if err != nil {
		zap.L().Error("投票时获取用户ID失败", zap.Error(err))
		ResponseError(ctx, CodeTokenInvaild)
		return
	}
	// 逻辑层处理
	if err := logic.VoteToBlog(uid, voteDate); err != nil {
		zap.L().Error("投票信息写入失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回结果
	ResponseSuccess(ctx, nil)
}
