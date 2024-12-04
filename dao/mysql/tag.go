package mysql

import (
	"NothingBlog/models"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrTagInfoNotFound = errors.New("记录不存在")
	ErrDeleteTagFailed = errors.New("删除Tag失败")
)

// 向数据库中插入数据
func InsertNewTag(tag *models.Tag) error {
	return Db.Create(tag).Error
}

// 查询所有的tag，以及所有的信息
func QueryTagAll() ([]models.Tag, error) {
	var tags []models.Tag
	if err := Db.Select("tag_id, name, desc, article_count, created_at, updated_at").Take(&tags).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tags, nil
}

func QueryTagById(tag *models.Tag) (err error) {
	err = Db.Where("tag_id=?", tag.TagId).Take(tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTagInfoNotFound
	}
	return
}

func QueryTagByName(tag *models.Tag) (err error) {
	err = Db.Where("name=?", tag.Name).Take(tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTagInfoNotFound
	}
	return
}

// 根据ID更新tag信息
func UpdateTagById(tag *models.Tag) (err error) {
	return Db.Model(&models.Tag{}).Updates(tag).Error
}

// // 根据名称更新tag信息
// func UpdateByName(tag *models.Tag) (err error) {

// }

func DeleteTagById(id int64) (err error) {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		zap.L().Warn("删除tag时,打开事务出错", zap.Error(tx.Error))
		return ErrDeleteTagFailed
	}

	// 删除tag_article表中对应的字段
	if tx.Exec("delete from `tag_article` where `tag_id` = ?", id).Error != nil {
		tx.Rollback()
		return err
	}

	// 删除tags表中的对应字段
	if tx.Exec("delete from `tags` where `tag_id` = ?", id).Error != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 删除多个tag，需要指定tag的id
func DeleteTagByIds(ids []int64) (err error) {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		zap.L().Warn("删除tag时,打开事务出错", zap.Error(tx.Error))
		return ErrDeleteTagFailed
	}

	// 删除tag_article表中对应的字段
	if tx.Exec("delete from `tag_article` where `tag_id` = in(?)", ids).Error != nil {
		tx.Rollback()
		return err
	}

	// 删除tags表中的对应字段
	if tx.Exec("delete from `tags` where `tag_id` = ?", ids).Error != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
