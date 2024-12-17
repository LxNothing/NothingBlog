package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
	"NothingBlog/package/utils"
	"NothingBlog/settings"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidUserId      = errors.New("用户ID无效")
	ErrArticleNameExisted = errors.New("文章名称已经存在")
	ErrArticleNotExisted  = errors.New("查询的文章不存在")
	ErrArticleQueryFailed = errors.New("查询数据库出错")
	ErrParamInvalid       = errors.New("参数无效")
)

var articleDb mysql.DaoArticle

type LogicArticle struct {
}

// CreateNewArticle 创建新文章
func (la LogicArticle) CreateNewArticle(article *models.Article, tagList []models.TagFormsParams) error {
	// 根据文章标题查询文章
	if atc, err := articleDb.QueryArticleByTitle(article.Title); atc != nil {
		zap.L().Debug("文章名称重复", zap.Error(err))
		return ErrArticleNameExisted
	}

	// 文章分类ID，文章tag 由 数据库插入的时候进行自动维护，因为建表的时候就关联了对应的键
	if article.Image == "" {
		article.Image = settings.Confg.DefaultAtcImg
	}

	// 生成文章id
	article.ArticleId = snowflake.GetNextId().Int64()

	// 如果文章摘要为空 - 可以使用文章的前多少个字符作为摘要
	if article.Summary == "" {
		article.Summary = "to do, default summary"
	}

	// 生成时间
	article.CreatedAt = time.Now()
	article.UpdatedAt = article.CreatedAt

	// 访问数据库 - 进行文章写入操作
	return articleDb.CreateArticle(article, tagList)
}

// UpdateArticle 更新已经存在的文章
func (la LogicArticle) UpdateArticle(atcId int64, newAtc *models.Article, tagList []models.TagFormsParams) error {
	var oldAtc *models.Article
	var err error
	// 根据文章标题查询文章 - 文章标题不允许重复
	if oldAtc, err = articleDb.QueryArticleByTitle(newAtc.Title); oldAtc != nil && oldAtc.ArticleId != atcId {
		zap.L().Debug("修改的文章名称不允许重复", zap.Error(err))
		return ErrArticleNameExisted
	}

	// 根据文章ID查询对应的文章是否存在
	if oldAtc, err = articleDb.QueryArticleById(atcId); oldAtc == nil {
		zap.L().Debug("更新文章时,传递文章ID不存在", zap.Error(err))
		return ErrArticleNotExisted
	}

	// 文章分类ID，文章tag 由 数据库插入的时候进行自动维护，因为建表的时候就关联了对应的键
	if newAtc.Image == "" {
		newAtc.Image = "to do, default image"
	}

	newAtc.ArticleId = atcId

	// 如果文章摘要为空 - 可以使用文章的前多少个字符作为摘要
	if newAtc.Summary == "" {
		newAtc.Summary = "to do, default summary"
	}

	// 访问数据库 - 对文章进行更新
	return articleDb.UpdateArticle(newAtc, oldAtc, tagList)
}

// GetArticleById 根据文章ID获取文章
func (la LogicArticle) GetArticleById(id int64) (*models.ArticleEntireReturn, error) {
	article, err := articleDb.QueryArticleById(id)
	if article == nil {
		if err == nil {
			return nil, ErrArticleNotExisted
		} else {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}

	return article.BindToEntireArticle(), nil
}

// 通过指定的参数查找对应的文章列表
func (la LogicArticle) GetAllWithParams(param *models.ArticleWithPageParams) ([]models.ArticleBriefReturn, int, int, error) {
	var class_id int64
	var tag_ids []int64
	var stus uint8
	var priv uint8
	var err error

	// 根据类别名称查类别ID
	if param.Class != "" {
		class, err := daoClass.QueryClassesByName(param.Class)
		if err != nil {
			return nil, 0, 0, ErrArticleQueryFailed
		}
		if class != nil {
			class_id = class.ClassId
		}
	}

	// 根据状态名称查找对应的编号
	if param.Status != "" {
		stus, err = models.StatusStringToNumber(param.Status)
		if err != nil {
			return nil, 0, 0, ErrParamInvalid
		}
	}

	// 根据权限名称查找对应的编号
	if param.Privilege != "" {
		priv, err = models.PrivilegeStringToNumber(param.Privilege)
		if err != nil {
			return nil, 0, 0, ErrParamInvalid
		}
	}

	// 获取tag的列表
	if param.Tag != "" {
		tag_ids, err = daoTag.QueryTagIdsByName(strings.Split(param.Tag, ","))
		if errors.Is(err, mysql.ErrTagOtherReason) {
			return nil, 0, 0, ErrArticleQueryFailed
		}
	}

	// 查询文章
	atcs, total, err := articleDb.QueryArticleWithParams(int(param.Page), int(param.Size), class_id, tag_ids, stus, priv, param.Keyword)
	if err != nil {
		zap.L().Warn("查询数据库出错", zap.Error(err))
		return nil, 0, 0, ErrArticleQueryFailed
	}

	briefAtc := make([]models.ArticleBriefReturn, len(atcs))
	for k, v := range atcs {
		briefAtc[k] = *v.BindToBriefArticle()
	}

	var cur_page = 1
	var total_page = 1
	if param.Page > 0 && param.Size > 0 && total != 0 {
		cur_page = int(param.Page)
		total_page = utils.GetTotalPage(int(param.Size), int(total))
	}

	return briefAtc, cur_page, total_page, nil
}

// func (la LogicArticle) GetAllArticle() ([]models.Article, error) {
// 	atcs, err := articleDb.QueryArticleAll()
// 	if atcs == nil {
// 		if err == nil {
// 			return nil, ErrArticleNotExisted
// 		} else {
// 			zap.L().Warn("查询数据库出错", zap.Error(err))
// 			return nil, ErrArticleQueryFailed
// 		}
// 	}
// 	return atcs, nil
// }

/*func (la LogicArticle) GetArticleByClassAndTagWithPage(clsId, tagId int64, page int, size int, state, privilege uint8) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, _, err := articleDb.QueryArticleByClassAndTagWithPage(clsId, tagId, page, size, state, privilege)

	tagIds := make([]int64, 1)
	tagIds[0] = tagId

	//atcs_tmp, err := articleDb.QueryArticleWithParams(page, size, clsId, tagIds, state, privilege, "")
	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	// for _, atc := range atcs_tmp {
	// 	atcs = append(atcs, *generateReturnBriefArticle(&atc))
	// }

	if err != nil {
		return nil, 0, err
	}

	return atcs, 0, nil
}

func (la LogicArticle) GetArticleByClassWithPage(classId int64, page, size int, state, privilege uint8) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := articleDb.QueryArticleByClassWithPage(classId, page, size, state, privilege)
	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	// for _, atc := range atcs_tmp {
	// 	atcs = append(atcs, *generateReturnBriefArticle(&atc))
	// }

	if err != nil {
		return nil, 0, err
	}

	return atcs, total, nil
}

func (la LogicArticle) GetArticleByTagWithPage(tagId int64, page, size int, state, privilege uint8) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := articleDb.QueryArticleByTagWithPage(tagId, page, size, state, privilege)
	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	// for _, atc := range atcs_tmp {
	// 	atcs = append(atcs, *generateReturnBriefArticle(&atc))
	// }

	if err != nil {
		return nil, 0, err
	}

	return atcs, total, nil
}

func (la LogicArticle) GetArticleWithPage(page, size int, state, privilege uint8, class int64) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := articleDb.QueryArticleWithPage(page, size, state, privilege, class)
	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	// for _, atc := range atcs_tmp {
	// 	atcs = append(atcs, *generateReturnBriefArticle(&atc))
	// }

	if err != nil {
		return nil, 0, err
	}

	return atcs, total, nil
}*/

func (la LogicArticle) DeleteArticleById(id int64) error {
	//return mysql.DeleteArticleById(id)
	ids := make([]int64, 1)
	ids[0] = id
	return articleDb.DeleteMultiArticleById(ids)
}

func (la LogicArticle) DeleteMultiArticleById(ids []int64) error {
	return articleDb.DeleteMultiArticleById(ids)
}

func (la LogicArticle) UpdateArticleStatusById(ids []int64, del bool) error {
	var stus = models.Recycle
	if !del {
		stus = models.Commit
	}
	num, _ := models.StatusStringToNumber(stus)
	return articleDb.UpdateArticleStatusById(ids, num)
}

// UpdateVisitCount 更新访问量
func (la LogicArticle) UpdateArticleVisitCount(id int64, newCount uint32) error {
	return articleDb.UpdateArticleVisitCountById(id, newCount)
}
