package mysql

import (
	"NothingBlog/models"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrArticleExisted    = errors.New("该文章已经存在")
	ErrArticleNotExist   = errors.New("该文章不存在")
	ErrCreateTransaction = errors.New("创建事务出错")
	ErrCreateArticle     = errors.New("创建文章出错")
	ErrQueryArticle      = errors.New("查询文章出错")
)

type DaoArticle struct {
}

func (da DaoArticle) CreateArticle(a *models.Article, tagList []models.TagFormsParams) error {
	// 开启事务
	tx := Db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// 创建事务错误
	if tx.Error != nil {
		return ErrCreateTransaction
	}

	// 插入数据
	if err := tx.Exec("INSERT INTO `articles`(`article_id`, `author_id`, `class_id`, `top_flag`, `en_comment`, `status`, `privilege`, `title`, `summary`, `image`, `content`) VALUES(?,?,?,?,?,?,?,?,?,?,?)",
		a.ArticleId, a.AuthorId, a.ClassId, a.TopFlag, a.EnComment,
		a.Status, a.Privilege, a.Title, a.Summary, a.Image, a.Content).Error; err != nil {
		tx.Rollback()
		zap.L().Error("向文章表插入数据出错", zap.Error(err))
		return ErrCreateArticle
	}

	// 更新文章所属分类的文章数量
	if err := tx.Exec("UPDATE classes SET atc_count = atc_count + 1 WHERE class_id = ?", a.ClassId).Error; err != nil {
		tx.Rollback()
		zap.L().Error("向类别表插入数据出错", zap.Error(err))
		return ErrCreateArticle
	}

	// 更新tag和文章关联的表 - 文章允许没有标签
	for _, v := range tagList {
		if err := tx.Exec("INSERT INTO tag_article (article_id, tag_id) values (?,?)", a.ArticleId, v.Id).Error; err != nil {
			tx.Rollback()
			zap.L().Error("更新文章和tag相关联的表时出错", zap.Error(err))
			return ErrCreateArticle
		}

		if err := tx.Exec("UPDATE tags SET article_count = article_count + 1 WHERE tag_id = ?", v.Id).Error; err != nil {
			tx.Rollback()
			zap.L().Error("更新tag表时出错", zap.Error(err))
			return ErrCreateArticle
		}
	}
	// 提交事务
	return tx.Commit().Error
}

func (da DaoArticle) UpdateArticle(new *models.Article, old *models.Article, tagList []models.TagFormsParams) error {
	tx := Db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		zap.L().Debug("更新文章时,开启事务出错")
		return tx.Error
	}

	// 更新文章表 - 使用map更新避免gorm默认不更新0值
	err := tx.Model(&models.Article{}).Where("`article_id` = ?", new.ArticleId).
		Updates(map[string]interface{}{
			"class_id":   new.ClassId,
			"top_flag":   new.TopFlag,
			"en_comment": new.EnComment,
			"status":     new.Status,
			"privilege":  new.Privilege,
			"title":      new.Title,
			"summary":    new.Summary,
			"image":      new.Image,
			"content":    new.Content,
		}).Error
	if err != nil {
		zap.L().Debug("更新文章时,更新文章表错误", zap.Error(err))
		tx.Rollback()
		return err
	}

	// 更新classes表 - 文章的class改变时才更新
	if new.ClassId != old.ClassId {
		// 旧类别的文章数量 - 1
		err = tx.Exec(`UPDATE classes 
						SET atc_count = atc_count - 1 
						WHERE (class_id = ? AND atc_count > 0)`, old.ClassId).Error
		if err != nil {
			zap.L().Debug("更新文章时,更新文章表错误", zap.Error(err))
			tx.Rollback()
			return err
		}
		// 新类别的文章数量 + 1
		err = tx.Exec(`UPDATE classes 
						SET atc_count = atc_count + 1 
						WHERE class_id = ?`, new.ClassId).Error
		if err != nil {
			zap.L().Debug("更新文章时,更新文章表错误", zap.Error(err))
			tx.Rollback()
			return err
		}
	}

	// 更新tag表 - 对tag表中的文章数量 - 1
	err = tx.Exec(`UPDATE tags 
					SET article_count = article_count - 1 
					WHERE tag_id IN (SELECT tag_id FROM tag_article WHERE article_id = ?) AND article_count > 0`, new.ArticleId).Error
	if err != nil {
		zap.L().Debug("更新文章时,更新Tag表错误", zap.Error(err))
		tx.Rollback()
		return err
	}

	// 更新tag表 - 删除标签文章表中关联记录
	err = tx.Exec("DELETE FROM tag_article WHERE article_id = ?", new.ArticleId).Error
	if err != nil {
		tx.Rollback()
		zap.L().Debug("更新文章时,更新Tag-article表错误")
		return err
	}

	// 重新建立tag表和文章的关联
	for _, v := range tagList {
		if err := tx.Exec("INSERT INTO tag_article (article_id, tag_id) values (?,?)", new.ArticleId, v.Id).Error; err != nil {
			tx.Rollback()
			zap.L().Debug("更新文章和tag相关联的表时出错", zap.Error(err))
			return ErrCreateArticle
		}

		if err := tx.Exec("UPDATE tags SET article_count = article_count + 1 WHERE tag_id = ?", v.Id).Error; err != nil {
			tx.Rollback()
			zap.L().Debug("更新tag表时出错", zap.Error(err))
			return ErrCreateArticle
		}
	}
	return tx.Commit().Error
}

// 根据文章名称查询文章
func (da DaoArticle) QueryArticleByTitle(title string) (*models.Article, error) {
	article := new(models.Article)
	err := Db.Preload("Class").Preload("TagList").Where("title=?", title).Take(article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return article, err
}

func (da DaoArticle) QueryArticleById(id int64) (*models.Article, error) {
	article := new(models.Article)

	err := Db.Preload("Class").Preload("TagList").Where("article_id=?", id).Take(article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return article, err
}

// QueryArticleAll 查询所有的文章
func (da DaoArticle) QueryArticleAll() ([]models.Article, error) {
	var articles []models.Article
	if err := Db.Preload("Class").Preload("TagList").Find(&articles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return articles, nil
}

// QueryArticleWithParams 根据指定的参数从数据库查询文章
// page 查询的页号, 当page=0，size=0时默认查询所有
// size 每页的大小
// class_id 文章所属类别的id, <0 时 查询所有类别
// tagId 文章所属包含的Tag列表，为nil时 查询所有tag
// stus 文章状态 =0时查询所有
// priv 文章权限 =0时查询所有
// key 模糊匹配的关键字，=""时查询所有
func (da DaoArticle) QueryArticleWithParams(page, size int, class_id int64, tagId []int64, stus, priv uint8, key string) ([]models.Article, int64, error) {
	var articles []models.Article

	query := Db.Preload("User").Preload("Class").Preload("TagList").Model(&models.Article{}).Order("top_flag desc, article_id desc")

	if stus != 0 {
		zap.L().Debug("query with stus", zap.Uint8("stus", stus))
		query = query.Where("status = ?", stus)
	}

	if priv != 0 {
		zap.L().Debug("query with priv")
		query = query.Where("privilege = ?", priv)
	}

	if class_id > 0 {
		zap.L().Debug("query with class_id")
		query = query.Where("class_id = ?", class_id)
	}

	if key != "" { // 按关键字模糊匹配 - 仅限标题
		zap.L().Debug("query with key", zap.String("key", key))
		query = query.Where("title like concat('%',?,'%')", key)
	}

	if tagId != nil { // 查询tag
		joinStr := `INNER JOIN tag_article 
				ON articles.article_id = tag_article.article_id 
				INNER JOIN tags 
				ON tag_article.tag_id = tags.tag_id`

		query = query.Joins(joinStr).Where("tags.tag_id in (?)", tagId)
	}

	// 只查询部分字段
	query = query.Select([]string{"articles.article_id", "articles.created_at",
		"articles.updated_at", "author_id", "class_id", "top_flag", "Status",
		"privilege", "title", "summary", "image", "comment_count", "visit_count"})

	// 查询满足要求的记录条数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		zap.L().Debug("查询文章总数出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}

	if page > 0 && size > 0 {
		query = query.Limit(size).Offset(size * (page - 1))
	}

	if err := query.Find(&articles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 未找到记录不应该返回错误
			return make([]models.Article, 0), 0, nil
		}
		zap.L().Debug("按页查找文章出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}

	return articles, total, nil
}

// QueryArticleWithPage 根据页和每页文章数查询文章信息，不包含文章具体内容信息
// 返回参数：
//  1. 文章信息
//  2. 满足指定条件的记录总数
//  3. 错误信息
/*func (da DaoArticle) QueryArticleWithPage(page, size int, state, privilege uint8, class_id int64) ([]models.Article, int64, error) {
	var articles []models.Article
	query := Db.Preload("User").Preload("Class").Preload("TagList").Model(&models.Article{}).Order("top_flag desc, article_id desc")
	query = query.Where("privilege = ? and status = ?", privilege, state)

	if class_id > 0 {
		query = query.Where("class_id = ?", class_id)
	}

	// 查询部分字段
	query = query.Select([]string{"article_id", "created_at", "updated_at", "author_id", "class_id", "top_flag",
		"Status", "privilege", "title", "summary", "image", "comment_count", "visit_count"})

	// 按页查找
	var total int64
	if err := query.Count(&total).Error; err != nil {
		zap.L().Debug("查询文章总数出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}

	utils.GetVaildPageAndSize(&page, &size, int(total))
	if err := query.Limit(size).Offset(size * (page - 1)).Find(&articles).Error; err != nil {
		zap.L().Debug("按页查找文章出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}

	return articles, total, nil
}*/

// QueryArticleByClass 获取指定类别下的文章
func (da DaoArticle) QueryArticleByClass(classId int64) ([]models.Article, error) {
	var articles []models.Article
	if err := Db.Preload("Class").Model(&models.Article{}).Where("class_id = ?", classId).Find(&articles).Error; err != nil {
		zap.L().Debug("通过class id查询文章出错", zap.Error(err))
		return nil, ErrQueryArticle
	}

	return articles, nil
}

func (da DaoArticle) QueryArticleNumberWithClass(classId int64) (int64, error) {
	var counter int64
	if err := Db.Model(&models.Article{}).Where("class_id = ?", classId).Count(&counter).Error; err != nil {
		zap.L().Debug("通过class id查询文章出错", zap.Error(err))
		return -1, ErrQueryArticle
	}
	return counter, nil
}

// QueryArticleByClass 获取指定类别下的文章
/*func (da DaoArticle) QueryArticleByClassWithPage(classId int64, page int, size int, state, privilege uint8) ([]models.Article, int64, error) {
	var articles []models.Article
	query := Db.Preload("Class").Preload("TagList").Preload("User").Model(&models.Article{}).Where("class_id = ?", classId)
	query = query.Where("privilege = ? and status = ?", privilege, state)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		zap.L().Debug("查询文章总数出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}

	utils.GetVaildPageAndSize(&page, &size, int(total))
	if err := query.Limit(size).Offset(size * (page - 1)).Find(&articles).Error; err != nil {
		zap.L().Debug("按页查找文章出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}
	return articles, total, nil
}

// QueryArticleByTag 获取指定Tag下的文章
func (da DaoArticle) QueryArticleByTag(tagId int64) ([]models.Article, error) {
	var articles []models.Article

	err := Db.Raw("select * from articles where 'article_id' in (select article_id from tag_article where tag_id = ?)", tagId).Scan(&articles).Error
	if err != nil {
		zap.L().Debug("通过Tag查询文章出错", zap.Error(err))
		return nil, ErrQueryArticle
	}
	return articles, nil
}

func (da DaoArticle) QueryArticleByTagWithPage(tagId int64, page, size int, state, privilege uint8) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64
	//Db.Preload("TagList").Joins("INNER JOIN tag_article ON articles.article_id = tag_article.article_id INNER JOIN tags ON tag_article.tag_id = tags.tag_id").Where("tags.tag_id = ?", tagId).Find(&articles)
	//fmt.Println(articles)

	queryStr := `select count(*)
				 from articles
				 where privilege = ? and status = ? and article_id in (select article_id from tag_article where tag_id = ?)`

	if err := Db.Raw(queryStr, privilege, state, tagId).Count(&total).Error; err != nil {
		zap.L().Debug("查询文章总数出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}

	//utils.GetVaildPageAndSize(&page, &size, int(total))
	// queryStr = `select *
	// 			from articles
	// 			where article_id in (select article_id from tag_article where tag_id = ?) limit ? offset ?`
	// if err := Db.Raw(queryStr, tagId, size, size*(page-1)).Scan(&articles).Error; err != nil {
	// 	zap.L().Debug("通过Tag按页查找文章出错", zap.Error(err))
	// 	return nil, 0, ErrQueryArticle
	// }
	joinStr := `INNER JOIN tag_article
				ON articles.article_id = tag_article.article_id
				INNER JOIN tags
				ON tag_article.tag_id = tags.tag_id`
	query := Db.Preload("TagList").Preload("User").Preload("Class").Joins(joinStr)
	query = query.Where("tags.tag_id = ?", tagId)
	query = query.Where("privilege = ? and status = ?", privilege, state)
	utils.GetVaildPageAndSize(&page, &size, int(total))
	if err := query.Limit(size).Offset(size * (page - 1)).Find(&articles).Error; err != nil {
		zap.L().Debug("通过Tag按页查找文章出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}
	return articles, total, nil
}

func (da DaoArticle) QueryArticleByClassAndTagWithPage(clsId int64, tagId int64, page int, size int, state, privilege uint8) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	queryStr := `select count(*) from articles
				where class_id = ? and privilege = ? and status = ? and article_id in (select article_id from tag_article where tag_id = ?)`

	if err := Db.Raw(queryStr, clsId, privilege, state, tagId).Count(&total).Error; err != nil {
		zap.L().Debug("查询文章总数出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}

	// queryStr = `select * from articles
	// 			where class_id = ? and article_id in (select article_id from tag_article where tag_id = ?)
	// 			limit ? offset ?`
	// utils.GetVaildPageAndSize(&page, &size, int(total))
	// if err := Db.Raw(queryStr, clsId, tagId, size, size*(page-1)).Scan(&articles).Error; err != nil {
	// 	zap.L().Debug("通过Tag和class按页查找文章出错", zap.Error(err))
	// 	return nil, 0, ErrQueryArticle
	// }

	// return articles, total, nil

	joinStr := `INNER JOIN tag_article
				ON articles.article_id = tag_article.article_id
				INNER JOIN tags
				ON tag_article.tag_id = tags.tag_id`
	query := Db.Preload("TagList").Preload("User").Preload("Class").Joins(joinStr)
	query = query.Where("tags.tag_id = ? and articles.class_id = ?", tagId, clsId)
	query = query.Where("privilege = ? and status = ?", privilege, state)
	utils.GetVaildPageAndSize(&page, &size, int(total))
	if err := query.Limit(size).Offset(size * (page - 1)).Find(&articles).Error; err != nil {
		zap.L().Debug("通过Tag按页查找文章出错", zap.Error(err))
		return nil, 0, ErrQueryArticle
	}
	return articles, total, nil
}*/

/*func DeleteArticleById(id int64) error {
	tx := Db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// 更新种类表（classes表）中的文章统计
	res := tx.Exec(`update classes
					set atc_count = atc_count - 1
					where class_id in (select class_id from articles where article_id = ?)
					and atc_count > 0`, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// 更新对应tag所对应的文章数量
	res = tx.Exec(`update tags
					set article_count = article_count - 1
					where tag_id in (select tag_id from tag_article where article_id = ?)
					and article_count > 0`, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// 删除关联文章和tag的中间表
	res = tx.Exec("delete from `tag_article` where `article_id` = ?", id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// 删除文章表中的文章
	res = tx.Exec("delete from `articles` where `article_id` = ?", id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return tx.Commit().Error
}*/

func (da DaoArticle) DeleteMultiArticleById(ids []int64) error {
	tx := Db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// 更新种类表（classes表）中的文章统计
	res := tx.Exec(`update classes 
					set atc_count = atc_count - 1 
					where class_id in (select class_id from articles where article_id in (?)) 
					and atc_count > 0`, ids)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// 删除对应文章的所有评论
	res = tx.Exec("delete from `comments` where `article_id` in (?)", ids)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// 更新对应tag所对应的文章数量
	res = tx.Exec(`update tags 
					set article_count = article_count - 1 
					where tag_id in (select tag_id from tag_article where article_id in (?)) 
					and article_count > 0`, ids)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// 删除关联文章和tag的中间表
	res = tx.Exec("delete from `tag_article` where `article_id` in (?)", ids)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// 删除文章表中的文章
	res = tx.Exec("delete from `articles` where `article_id` in (?)", ids)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return tx.Commit().Error
}

// UpdateArticleStatusById 更新文章的状态
func (da DaoArticle) UpdateArticleStatusById(ids []int64, tp uint8) error {
	return Db.Model(&models.Article{}).Where("article_id in (?)", ids).UpdateColumn("status", tp).Error
}

// UpdateArticleVisitCountById 更新文章的访问量
func (da DaoArticle) UpdateArticleVisitCountById(id int64, newCount uint32) error {
	return Db.Model(&models.Article{}).Where("article_id = ?", id).UpdateColumn("visit_count", newCount).Error
}
