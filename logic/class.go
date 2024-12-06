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
	ErrDeleteClassByIds = errors.New("该类别中存在文章")
)

func CreateNewClass(cls *models.Class) error {
	cls.ClassId = snowflake.GetNextId().Int64()
	cls.CreatedAt = time.Now()
	cls.UpdatedAt = cls.CreatedAt

	return mysql.InsertNewClass(cls)
}

func GetClassById(id int64) (*models.Class, error) {
	cls, err := mysql.QueryClassesById(id)

	if cls == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	return cls, nil
}

func GetClassByName(name string) (*models.Class, error) {
	cls, err := mysql.QueryClassesByName(name)

	if cls == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	return cls, nil
}

func GetAllClasses() ([]models.Class, error) {
	classes, err := mysql.QueryAllClasses()

	if classes == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	return classes, nil
}

func DeleteClassById(id int64) error {
	ids := make([]int64, 1)
	ids[0] = id

	// 先判断是否存在文章的class是待删class
	num, err := mysql.QueryAticleNumberByClass(ids)
	if err != nil || num != 0 {
		return ErrDeleteClassByIds
	}

	return mysql.DeleteClassByIds(ids)
}

func DeleteMultiClassById(ids []int64) error {
	// 先判断是否存在文章的class是待删class
	num, err := mysql.QueryAticleNumberByClass(ids)
	if err != nil || num != 0 {
		return ErrDeleteClassByIds
	}

	return mysql.DeleteClassByIds(ids)
}
