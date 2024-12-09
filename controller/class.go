package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"
	"NothingBlog/package/utils"
	"NothingBlog/settings"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAllClassesHandler 获取所有的类别
// @Summary 获取所有类别（简略信息）的接口
// @Description 通过该接口可以获得当前的所有文章
// @Tags 类别相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseAllClassesList
// @Router /classes [get]
func GetAllClassesHandler(ctx *gin.Context) {
	cls, err := logic.GetAllClasses()
	if err != nil {
		zap.L().Debug("获取所有的class失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccess(ctx, cls)
}

// GetClassByIdHandler 根据类别ID查询类的信息
// @Summary 根据类别ID查询类的信息（详细信息）
// @Description 通过该接口可以获得指定ID的类别
// @Tags 类别相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseClassDetailList "code字段为1000表示执行成功，其余表示出错"
// @Router /class/:id [get]
func GetClassByIdHandler(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		zap.L().Debug("通过ID获取class失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	cls, err := logic.GetClassById(id)
	if err != nil {
		zap.L().Debug("通过ID获取class失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccess(ctx, cls)
}

// CreateClassHandler 创建class
// @Summary 创建类别(Class)的接口
// @Description 通过该接口可以创建类别
// @Tags 类别相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.ClassCreateFormParams true "创建类别的参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCreateClass "code=1000表示成功其余失败"
// @Router /class [post]
func CreateClassHandler(ctx *gin.Context) {
	var clsId int64
	var err error
	classParam := new(models.ClassCreateFormParams)
	if err = ctx.ShouldBindJSON(classParam); err != nil {
		zap.L().Debug("解析创建Tag的参数失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	cls := new(models.Class)
	cls.Name = classParam.Name
	cls.Desc = classParam.Desc

	if clsId, err = logic.CreateNewClass(cls); err != nil {
		zap.L().Debug("创建Tag失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccess(ctx, models.ResponseClassBrief{
		ClassId:  clsId,
		Name:     cls.Name,
		AtcCount: 0,
	})
}

// category?class=all 这个接口是供客户端使用的，不需要进行jwt认证
func GetAllClassClientHandler(ctx *gin.Context) {
	var clsId int64
	var tagId int64
	var page int = 1
	// 获取参数
	clsStr := ctx.Query("class")
	tagStr := ctx.Query("tag")
	pageStr := ctx.Query("page")

	if clsStr == "" {
		clsStr = "all"
	}

	if tagStr == "" {
		tagStr = "all"
	}

	if pageStr != "" {
		pg, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			zap.L().Debug("页号传递错误，解析失败", zap.Error(err))
			ResponseError(ctx, CodeParameterInvalid)
			return
		}
		page = int(pg)
	}

	// 如果客户端指定了class，则判断class是否已经存在
	if clsStr != "all" {
		cls, _ := logic.GetClassByName(clsStr) // 查询对应的class是否存在
		if cls == nil {
			zap.L().Debug("传递的查询文章种类参数错误")
			ResponseErrorWithMsg(ctx, CodeParameterInvalid, "class名错误")
			return
		}
		clsId = cls.ClassId
	}

	// 判断tag是否已经存在
	if tagStr != "all" {
		tag, _ := logic.GetTagByName(tagStr) // 查询对应的class是否存在
		if tag == nil {
			zap.L().Debug("传递的查询tag种类参数错误")
			ResponseErrorWithMsg(ctx, CodeParameterInvalid, "tag名错误")
			return
		}
		tagId = tag.TagId
	}

	// 获取所有的类别
	classes, err := logic.GetAllClasses()
	if err != nil {
		zap.L().Debug("查询文章种类失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	category := &models.ResponseClassAllForClient{
		CurClassName: clsStr,
		CurTagName:   tagStr,
		BriefClasses: make([]models.ResponseClassBrief, len(classes)),
	}
	for idx, class := range classes {
		category.BriefClasses[idx] = models.ResponseClassBrief{
			ClassId:  class.ClassId,
			AtcCount: class.AtcCount,
			Name:     class.Name,
		}
	}

	// 根据class名称，获取对应章的所有tag
	var tags []models.ResponseTagBrief
	if clsStr == "all" {
		tags, _ = logic.GetAllTags() // class = all, 则获取所有的tag信息
	} else {
		// class != all, 则获取该类别下的所有文章的不重复tag
		tags, err = logic.GetTagByClassId(clsId)
		if err != nil {
			zap.L().Debug("查询tag失败", zap.Error(err))
			ResponseError(ctx, CodeServerBusy)
			return
		}
	}
	category.CurTagList = make([]models.ResponseTagBrief, 0, len(tags))
	for _, tag := range tags {
		category.CurTagList = append(category.CurTagList, models.ResponseTagBrief{
			TagId:    tag.TagId,
			Name:     tag.Name,
			AtcCount: tag.AtcCount,
		})
	}

	// 获取其中所属的文章列表 - 根据class 和 tag进行筛选
	// class = all tag = all - 即不包含这两个筛选条件，通过时间或者阅读量进行排序
	// class = all tag = 具体的tag - 不包含类别条件，使用tag进行查询，所得结果通过时间或者阅读量进行排序
	// class = 具体的class tag = all
	// class = 具体的class tag = 具体的tag - 使用两者进行排序
	// 都是 all 则 不包含这两个筛选条件，通过时间或者阅读量进行排序

	var pagesize = int(settings.Confg.PageSize)
	var atcs []models.ArticleBriefReturn
	var total int64

	if clsStr == "all" {
		if tagStr == "all" {
			atcs, total, err = logic.GetArticleWithPage(page, pagesize, models.StatusDraft, models.PrivilegePrivte, 0)

		} else {
			atcs, total, err = logic.GetArticleByTagWithPage(tagId, page, pagesize, models.StatusDraft, models.PrivilegePrivte)
		}
	} else {
		if tagStr == "all" {
			atcs, total, err = logic.GetArticleByClassWithPage(clsId, page, pagesize, models.StatusDraft, models.PrivilegePrivte)

		} else {
			atcs, total, err = logic.GetArticleByClassAndTagWithPage(clsId, tagId, page, pagesize, models.StatusDraft, models.PrivilegePrivte)
		}
	}

	if err != nil {
		zap.L().Debug("分页查询文章出错", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 计算总的页数
	total_page := utils.GetTotalPage(int64(pagesize), total)
	// 返回结果
	ResponseSuccess(ctx, gin.H{
		"category":   category,
		"cur_page":   min(page, int(total_page)), // 当前返回的页数-超出总页数返回最后一页
		"total_page": total_page,                 // 总共有多少页
		"article":    atcs,
	})
}

// DeleteClassHandler 删除class
// @Summary 删除类别(Class)的接口
// @Description 通过该接口可以删除类别
// @Tags 类别相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功"
// @Router /class/:id [delete]
func DeleteClassHandler(ctx *gin.Context) {
	clsId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		zap.L().Debug("Class的ID参数传递错误", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeParameterInvalid, "class id指定错误")
		return
	}

	// 类别下存在文章则不允许删除
	err = logic.DeleteOneClassById(clsId)
	if err != nil {
		zap.L().Debug("删除Class失败", zap.Error(err))
		if errors.Is(err, logic.ErrDeleteClassByIds) {
			ResponseError(ctx, CodeHaveArticleInClass)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "删除Class成功", nil)
}

// DeleteMultiClassHandler 删除多个class
// @Summary 删除多个类别(Class)的接口
// @Description 通过该接口可以删除多个类别
// @Tags 类别相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.DeleteMultiTagParams true "待删除的class ID列表"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseDeleteClass "code=1000成功，,返回data域中不为空代表对应的class删除失败"
// @Router /classes [delete]
func DeleteMultiClassHandler(ctx *gin.Context) {
	param := new(models.DeleteMultiTagParams)
	if err := ctx.ShouldBindJSON(param); err != nil {
		zap.L().Debug("Class的ID参数传递错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	ids, err := logic.DeleteMultiClassById(param.Ids)
	if err != nil {
		zap.L().Error("删除Class失败", zap.Error(err))
		if errors.Is(err, logic.ErrDeleteClassByIds) {
			ResponseErrorWithDataMsg(ctx, CodeHaveArticleInClass, "类别下存在文章不允许删除", ids)
			return
		}

		ResponseErrorWithDataMsg(ctx, CodeServerBusy, "服务器繁忙，删除失败", ids)
		return
	}
	if len(ids) == 0 {
		ids = nil
	}
	ResponseSuccessWithMsg(ctx, "删除Class成功", ids)
}

// UpdateClassHandler 更新class
// @Summary 更新类别(Class)的接口
// @Description 通过该接口可以更新类别
// @Tags 类别相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token(jwt)"
// @Param object body models.UpdateClassParams true "更新类别的参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseNoDataArea "code=1000表示成功其余失败"
// @Router /class [put]
func UpdateClassHandler(ctx *gin.Context) {
	newClass := new(models.UpdateClassParams)

	if err := ctx.ShouldBindJSON(newClass); err != nil {
		zap.L().Debug("解析修改Class的参数错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// // 查找ID是否存在
	_, err := logic.GetClassById(newClass.ClassId)
	if err != nil {
		zap.L().Debug("ID错误", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeParameterInvalid, "要修改的Class不存在")
		return
	}

	// 查找名称是否重复
	oldClass, err := logic.GetClassByName(newClass.Name)
	if err == nil && oldClass.ClassId > 0 && oldClass.ClassId != newClass.ClassId {
		zap.L().Debug("Class名称重复")
		ResponseError(ctx, CodeClassNameExisted) // 名称重复
		return
	}

	var class = &models.Class{
		ClassId: newClass.ClassId,
		Name:    newClass.Name,
		Desc:    newClass.Desc,
	}

	if err := logic.UpdateClass(class); err != nil {
		zap.L().Debug("更新Class失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "修改类别成功", nil)
}
