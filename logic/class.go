package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
)

func CreateNewClass(cls *models.Class) error {
	cls.ClassId = snowflake.GetNextId().Int64()

	return mysql.InsertNewClass(cls)
}
