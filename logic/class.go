package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
	"time"

	"go.uber.org/zap"
)

func CreateNewClass(cls *models.Class) error {
	cls.ClassId = snowflake.GetNextId().Int64()
	cls.CreatedAt = time.Now()
	cls.UpdatedAt = cls.CreatedAt

	return mysql.InsertNewClass(cls)
}

func GetClassByName(name string) (*models.Class, error) {
	cls, err := mysql.QueryClassesByName(name)

	if cls == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	return cls, err
}

func GetAllClasses() ([]models.Class, error) {
	classes, err := mysql.QueryAllClasses()

	if classes == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	return classes, err
}
