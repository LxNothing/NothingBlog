package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
	"NothingBlog/settings"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidUserId      = errors.New("用户ID无效")
	ErrArticleNameExisted = errors.New("文章名称已经存在")
	ErrArticleNotExisted  = errors.New("查询的文章不存在")
	ErrArticleQueryFailed = errors.New("查询数据库出错")
)

// generateReturnBriefArticle 将完整的文章数据结构映射到简略的文章结构中
func generateReturnBriefArticle(atc *models.Article) *models.ArticleBriefReturn {
	tmp := &models.ArticleBriefReturn{
		ArticleId:  atc.ArticleId,
		CreatedAt:  atc.CreatedAt,
		UpdatedAt:  atc.UpdatedAt,
		AuthorId:   atc.AuthorId,
		AuthorName: atc.User.UserName,
		ClassId:    atc.ClassId,
		ClassName:  atc.Class.Name,
		Title:      atc.Title,
		Privilege:  atc.Privilege,
		EnComment:  atc.EnComment,
		TopFlag:    atc.TopFlag,
		Image:      atc.Image,
		Summary:    atc.Summary,
	}
	tmp.TagId = make([]int64, 0, len(atc.TagList))
	tmp.TagName = make([]string, 0, len(atc.TagList))
	for _, tag := range atc.TagList {
		tmp.TagId = append(tmp.TagId, tag.TagId)
		tmp.TagName = append(tmp.TagName, tag.Name)
	}
	return tmp
}

func generateReturnEntireArticle(atc *models.Article) *models.ArticleEntireReturn {
	tmp := generateReturnBriefArticle(atc)
	return &models.ArticleEntireReturn{
		ArticleBriefReturn: *tmp,
		Content:            atc.Content,
	}
}

// CreateNewArticle 创建新文章
func CreateNewArticle(article *models.Article, tagList []models.TagFormsParams) error {
	// 根据文章标题查询文章
	if atc, err := mysql.QueryArticleByTitle(article.Title); atc != nil {
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
	return mysql.CreateArticle(article, tagList)
}

// UpdateArticle 更新已经存在的文章
func UpdateArticle(atcId int64, newAtc *models.Article, tagList []models.TagFormsParams) error {
	var oldAtc *models.Article
	var err error
	// 根据文章标题查询文章 - 文章标题不允许重复
	if oldAtc, err = mysql.QueryArticleByTitle(newAtc.Title); oldAtc != nil && oldAtc.ArticleId != atcId {
		zap.L().Debug("修改的文章名称不允许重复", zap.Error(err))
		return ErrArticleNameExisted
	}

	// 根据文章ID查询对应的文章是否存在
	if oldAtc, err = mysql.QueryArticleById(atcId); oldAtc == nil {
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
	return mysql.UpdateArticle(newAtc, oldAtc, tagList)
}

// GetArticleById 根据文章ID获取文章
func GetArticleById(id int64) (*models.ArticleEntireReturn, error) {
	article, err := mysql.QueryArticleById(id)
	if article == nil {
		if err == nil {
			return nil, ErrArticleNotExisted
		} else {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}

	return generateReturnEntireArticle(article), nil
}

func GetAllArticle() ([]models.Article, error) {
	atcs, err := mysql.QueryArticleAll()
	if atcs == nil {
		if err == nil {
			return nil, ErrArticleNotExisted
		} else {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	return atcs, nil
}

func GetArticleByClassAndTagWithPage(clsId, tagId int64, page int, size int, state models.StatusType, privilege models.PrivilegeType) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := mysql.QueryArticleByClassAndTagWithPage(clsId, tagId, page, size, state, privilege)

	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	for _, atc := range atcs_tmp {
		atcs = append(atcs, *generateReturnBriefArticle(&atc))
	}

	if err != nil {
		return nil, 0, err
	}

	return atcs, total, nil
}

func GetArticleByClassWithPage(classId int64, page int, size int, state models.StatusType, privilege models.PrivilegeType) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := mysql.QueryArticleByClassWithPage(classId, page, size, state, privilege)
	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	for _, atc := range atcs_tmp {
		atcs = append(atcs, *generateReturnBriefArticle(&atc))
	}

	if err != nil {
		return nil, 0, err
	}

	return atcs, total, nil
}

func GetArticleByTagWithPage(tagId int64, page int, size int, state models.StatusType, privilege models.PrivilegeType) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := mysql.QueryArticleByTagWithPage(tagId, page, size, state, privilege)
	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	for _, atc := range atcs_tmp {
		atcs = append(atcs, *generateReturnBriefArticle(&atc))
	}

	if err != nil {
		return nil, 0, err
	}

	return atcs, total, nil
}

func GetArticleWithPage(page int, size int, state models.StatusType, privilege models.PrivilegeType, class int64) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := mysql.QueryArticleWithPage(page, size, state, privilege, class)
	atcs := make([]models.ArticleBriefReturn, 0, len(atcs_tmp))

	for _, atc := range atcs_tmp {
		atcs = append(atcs, *generateReturnBriefArticle(&atc))
	}

	if err != nil {
		return nil, 0, err
	}

	return atcs, total, nil
}

func DeleteArticleById(id int64) error {
	//return mysql.DeleteArticleById(id)
	ids := make([]int64, 1)
	ids[0] = id
	return mysql.DeleteMultiArticleById(ids)
}

func DeleteMultiArticleById(ids []int64) error {
	return mysql.DeleteMultiArticleById(ids)
}

func UpdateArticleStatusById(ids []int64, del bool) error {
	var stus = models.StatusDelete
	if !del {
		stus = models.StatusCommit
	}
	return mysql.UpdateArticleStatusById(ids, stus)
}

// UpdateVisitCount 更新访问量
func UpdateArticleVisitCount(id int64, newCount uint32) error {
	return mysql.UpdateArticleVisitCountById(id, newCount)
}
