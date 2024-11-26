package controller

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/logic"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 获取社区的列表
func CommunityListHandler(ctx *gin.Context) {
	// 参数校验 - 无参数需要校验
	// 调用逻辑层的社区获取接口 - 以community_id 以community_name的形式进行返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("获取社区信息失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy) // 返回这个的目的是向前端隐藏真实的错误信息
		return
	}
	// 返回消息
	ResponseSuccess(ctx, data)
}

// 获取社区详情
func CommunityDetailHandler(ctx *gin.Context) {
	// 参数处理 - 这里使用路径参数 - 获取社区id
	cid, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 调用逻辑层的社区获取接口 - 以community_id 以community_name的形式进行返回
	data, err := logic.GetCommunityInfoById(cid)
	if err != nil {
		if errors.Is(err, mysql.ErrInvalidCommunityId) {
			zap.L().Error("社区参数无效", zap.Error(err))
			ResponseError(ctx, CodeCommunityIdInvalid) // 返回这个的目的是向前端隐藏真实的
			return
		}
		zap.L().Error("获取社区信息失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy) // 返回这个的目的是向前端隐藏真实的错误信息
		return
	}
	// 返回消息
	ResponseSuccess(ctx, data)
}
