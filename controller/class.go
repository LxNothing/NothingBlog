package controller

import (
	"NothingBlog/logic"
	"NothingBlog/models"
	"NothingBlog/package/utils"
	"NothingBlog/settings"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateClassHandler(ctx *gin.Context) {
	classParam := new(models.ClassCreateFormParams)
	if err := ctx.ShouldBindJSON(classParam); err != nil {
		zap.L().Debug("解析创建Tag的参数失败", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	cls := new(models.Class)
	cls.Name = classParam.Name
	cls.Desc = classParam.Desc

	if err := logic.CreateNewClass(cls); err != nil {
		zap.L().Debug("创建Tag失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccessWithMsg(ctx, "创建Class成功", nil)
}

type category struct {
	CurClassName string
	CurTagName   string
	AllClasses   []models.ClassBriefReturn
	CurTagList   []models.TagBriefReturn
}

// category?class=all
func GetAllClassHandler(ctx *gin.Context) {
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
			zap.L().Debug("页号传递错误，解析失败")
			ResponseError(ctx, CodeParameterInvalid)
			return
		}
		page = int(pg)
	}

	// 判断class是否已经存在
	if clsStr != "all" {
		cls, _ := logic.GetClassByName(clsStr) // 查询对应的class是否存在
		if cls == nil {
			zap.L().Debug("传递的查询文章种类参数错误")
			ResponseError(ctx, CodeParameterInvalid)
			return
		}
		clsId = cls.ClassId
	}

	// 判断tag是否已经存在
	if tagStr != "all" {
		tag, _ := logic.GetTagByName(tagStr) // 查询对应的class是否存在
		if tag == nil {
			zap.L().Debug("传递的查询tag种类参数错误")
			ResponseError(ctx, CodeParameterInvalid)
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
	category := &category{
		CurClassName: clsStr,
		CurTagName:   tagStr,
		AllClasses:   make([]models.ClassBriefReturn, 0, len(classes)),
	}
	for _, class := range classes {
		category.AllClasses = append(category.AllClasses, models.ClassBriefReturn{
			ClassId:  class.ClassId,
			AtcCount: class.AtcCount,
			Name:     class.Name,
		})
	}

	// 根据class名称，获取对应章的所有tag
	var tags []models.Tag
	if clsStr == "all" {
		tags, _ = logic.GetAllTags()
	} else {
		cls, err := logic.GetClassByName(clsStr)
		if err != nil {
			zap.L().Debug("查询class ID失败", zap.Error(err))
			ResponseError(ctx, CodeServerBusy)
			return
		}

		tags, err = logic.GetTagByClassId(cls.ClassId)
		if err != nil {
			zap.L().Debug("查询tag失败", zap.Error(err))
			ResponseError(ctx, CodeServerBusy)
			return
		}
	}
	category.CurTagList = make([]models.TagBriefReturn, 0, len(tags))
	for _, tag := range tags {
		category.CurTagList = append(category.CurTagList, models.TagBriefReturn{
			TagId:        tag.TagId,
			Name:         tag.Name,
			ArticleCount: tag.ArticleCount,
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
