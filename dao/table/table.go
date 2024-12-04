package table

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
)

func DbTableInit() error {
	return mysql.Db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Class{}, &models.Article{})
}
