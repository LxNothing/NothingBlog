package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var logicComment = logic.LogicComment{}

func GetCommentForClientHandler(ctx *gin.Context) {
	atcIdStr := ctx.Query("atcid")

	if atcIdStr == "" {
		ResponseErrorWithMsg(ctx, CodeParameterInvalid, "文章id不能为空")
		return
	}

	commens, err := logicComment.GetWithPageForClient(atcIdStr)
	fmt.Println(commens)
	if err != nil && !errors.Is(err, logic.ErrCommentNotFound) {
		ResponseErrorWithMsg(ctx, CodeServerBusy, "评论查找失败")
		return
	}
	ResponseSuccessWithMsg(ctx, "评论查询成功", commens)
}

// GetCommentWithPageHandler 分页获取评论
// @Summary 分页获取评论的接口
// @Description 通过该接口可以分页获取评论
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param keyword query string false "关键字"
// @Param type query uint8 false "评论类型 1-文章"
// @Param atc_id query string false "对应的文章id"
// @Param status query uint8 false "状态: 1-审核中 2-审核通过 3-审核未通过"
// @Param page_idx query  uint false "页号"
// @Param size query  uint false "每页大小"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /comments [get]
func GetCommentWithPageHandler(ctx *gin.Context) {
	var param models.CommentWithPageParams

	if err := ctx.ShouldBindQuery(&param); err != nil {
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 调用logic层的方法进行数据分页查询
	res, cur_page, total_page, err := logicComment.GetWithPageForAdmin(&param)
	if err != nil {
		ResponseErrorWithMsg(ctx, CodeServerBusy, err.Error())
		return
	}

	resList := make([]models.ResponseCommentListForAdmin, len(res))
	for idx, v := range res {
		resList[idx] = *(v.BindToResponseForAdmin())
	}

	ResponseSuccessWithMsg(ctx, "评论查询成功", gin.H{
		"cur_page":   cur_page,
		"total_page": total_page,
		"comments":   resList,
	})
}

// CreateCommentHandler 创建一条新的评论
// @Summary 创建评论Comment的接口
// @Description 通过该接口可以创建评论
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param object body models.CommentCreateFormParams true "创建评论的参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /comment [post]
func CreateCommentHandler(ctx *gin.Context) {
	cmtParam := new(models.CommentCreateFormParams)
	if err := ctx.ShouldBindJSON(cmtParam); err != nil {
		ResponseErrorWithMsg(ctx, CodeParameterInvalid, "评论参数传递错误")
		return
	}

	// 数据类型转换
	comment := cmtParam.ParamToDbModel()

	// 调用逻辑层的创建评论的方法
	if err := logicComment.CreateComment(comment); err != nil {
		ResponseErrorWithMsg(ctx, CodeServerBusy, err.Error())
		return
	}
	// 返回结果
	ResponseSuccessWithMsg(ctx, "创建评论成功", nil)
}

// DeleteCommentById 删除指定ID的评论
// @Summary 删除单条评论Comment的接口
// @Description 通过该接口可以删除单条评论
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param id path uint true "删除的评论的ID"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /comment/:id [delete]
func DeleteCommentById(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := logicComment.DeleteCommentById(id); err != nil {
		ResponseErrorWithMsg(ctx, CodeServerBusy, err.Error())
		return
	}
	ResponseSuccessWithMsg(ctx, "删除评论成功", nil)
}

// DeleteCommentsByIds 删除指定ID的评论
// @Summary 删除多条评论Comment的接口
// @Description 通过该接口可以删除多条评论
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.CommentDeleteFormParams true "删除的评论的ID"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /comments [delete]
func DeleteCommentsByIds(ctx *gin.Context) {
	idsParams := new(models.CommentDeleteFormParams)

	if err := ctx.ShouldBindJSON(idsParams); err != nil {
		ResponseErrorWithMsg(ctx, CodeParameterInvalid, "参数传递无效")
		return
	}

	if err := logicComment.DeleteCommentByIds(idsParams.Ids); err != nil {
		ResponseErrorWithMsg(ctx, CodeServerBusy, err.Error())
		return
	}
	ResponseSuccessWithMsg(ctx, "删除评论成功", nil)
}

// UpdateCommentStateHandler 修改指定评论的状态
// @Summary 修改指定评论的状态的接口
// @Description 通过该接口可以修改评论的状态
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.CommentDeleteFormParams true "删除的评论的ID"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /comments [delete]
func UpdateCommentStateHandler(ctx *gin.Context) {
	stParam := new(models.CommentUpdateStateParams)
	if err := ctx.ShouldBindJSON(stParam); err != nil {
		ResponseErrorWithMsg(ctx, CodeParameterInvalid, "参数传递无效")
		return
	}

	err := logicComment.UpdateCommentStatus(stParam.Id, models.CommentStatusType(stParam.Value))
	if err != nil {
		if errors.Is(err, logic.ErrCommentParamInvalid) {
			ResponseErrorWithMsg(ctx, CodeParameterInvalid, err.Error())
			return
		}
		ResponseErrorWithMsg(ctx, CodeServerBusy, "更新失败,服务器繁忙")
		return
	}

	ResponseSuccessWithMsg(ctx, "更新评论状态成功", nil)
}
