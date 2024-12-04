package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
)

func CreateNewTag(tag *models.Tag) error {
	tag.TagId = snowflake.GetNextId().Int64()

	return mysql.InsertNewTag(tag)
}
