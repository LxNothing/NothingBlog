package table

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
)

// sql表的自动迁移，实际生产环境一般不用 - 表的迁移都是很慎重的操作，需要谨慎
func DbTableInit() error {
	return mysql.Db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Class{},
		&models.Article{}, &models.Comment{})
}
