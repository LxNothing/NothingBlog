package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
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

func CreateNewArticle(article *models.Article, tagList []models.TagFormsParams) error {
	// 根据文章标题查询文章
	if err := mysql.QueryArticleByTitle(article.Title); err != nil {
		zap.L().Debug("文章名称重复", zap.Error(err))
		return ErrArticleNameExisted
	}

	// 文章分类ID，文章tag 由 数据库插入的时候进行自动维护，因为建表的时候就关联了对应的键
	if article.Image == "" {
		article.Image = "to do, default image"
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

func GetArticleByClassAndTagWithPage(clsId, tagId int64, page int, size int) ([]models.ArticleBriefReturn, int64, error) {
	atcs_tmp, total, err := mysql.QueryArticleByClassAndTagWithPage(clsId, tagId, page, size)

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
	atcs_tmp, total, err := mysql.QueryArticleByClassWithPage(classId, page, size)
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
	atcs_tmp, total, err := mysql.QueryArticleByTagWithPage(tagId, page, size)
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
