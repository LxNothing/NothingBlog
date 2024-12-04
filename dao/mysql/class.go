package mysql

import "NothingBlog/models"

// 向数据库中插入数据
func InsertNewClass(cls *models.Class) error {
	return Db.Create(cls).Error
}
