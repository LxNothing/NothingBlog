package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BlogCreaterHandler(ctx *gin.Context) {
	// 验证post的数据
	blog := new(models.BlogsArch)
	if err := ctx.ShouldBindJSON(blog); err != nil {
		zap.L().Error("解析创建文章的参数有误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 验证社区id是否符合要求
	if !logic.QueryCommunityExistedById(blog.CommmunityId) {
		zap.L().Error("社区id传递错误")
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 获取用户的uid
	uid, err := getCurrentUserId(ctx)
	if err != nil {
		zap.L().Error("获取用户id失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	blog.AuthorId = uid

	// 调用logic层的方法处理请求
	if err := logic.CreateNewBlog(blog); err != nil {
		zap.L().Error("logic.CreateNewBlog 创建文章出错", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回结果
	ResponseSuccess(ctx, nil)
}

func BlogDetailHandler(ctx *gin.Context) {
	// 获取url id参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("博文ID解析失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}
	// 调用logic层方法获取数据
	blog, err := logic.GetBlogDetailById(id)
	if err != nil {
		if errors.Is(err, logic.ErrBlogIdNotExisted) {
			zap.L().Error("文章ID不存在", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		if errors.Is(err, logic.ErrInvalidCommId) {
			zap.L().Error("社区ID无效", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		if errors.Is(err, logic.ErrInvalidUserId) {
			zap.L().Error("用户ID无效", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		zap.L().Error("查询博文数据库失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回数据
	ResponseSuccess(ctx, blog)
}

func BlogListHandler(ctx *gin.Context) {
	page, size, _ := getBlogSizeAndPage(ctx)

	// 调用逻辑层的函数
	blogs, err := logic.GetBlogList(page, size)
	if err != nil {
		if errors.Is(err, logic.ErrBlogIdNotExisted) {
			zap.L().Error("文章ID不存在", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		if errors.Is(err, logic.ErrInvalidCommId) {
			zap.L().Error("社区ID无效", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		if errors.Is(err, logic.ErrInvalidUserId) {
			zap.L().Error("用户ID无效", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		zap.L().Error("查询博文数据库失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回数据
	ResponseSuccess(ctx, blogs)
}

// 支持前端传递排序参数，通过排序参数进行查询
// 这里利用了在 redis中的两个zset - time和score，这两个zset中存储的是依据创建时间和当前分数排序后的
// 文章id
func BlogOrderListHandler(ctx *gin.Context) {
	bgParams, _ := getBlogListParams(ctx) // 执行出错则返回默认值

	// 调用逻辑层的函数
	blogs, err := logic.GetBlogOrderList(bgParams.Page, bgParams.Szie, bgParams.Order)
	if err != nil {
		if errors.Is(err, logic.ErrBlogIdNotExisted) {
			zap.L().Error("文章ID不存在", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		if errors.Is(err, logic.ErrInvalidCommId) {
			zap.L().Error("社区ID无效", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		if errors.Is(err, logic.ErrInvalidUserId) {
			zap.L().Error("用户ID无效", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}

		zap.L().Error("查询博文数据库失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回数据
	ResponseSuccess(ctx, blogs)
}
