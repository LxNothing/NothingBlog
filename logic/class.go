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
	ErrQueryDatabase    = errors.New("查询数据库出错")
)

func CreateNewClass(cls *models.Class) (int64, error) {
	cls.ClassId = snowflake.GetNextId().Int64()
	cls.CreatedAt = time.Now()
	cls.UpdatedAt = cls.CreatedAt

	return cls.ClassId, mysql.InsertNewClass(cls)
}

func GetClassById(id int64) (*models.ResponseClassDetail, error) {
	cls, err := mysql.QueryClassesById(id)

	if cls == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}

	res := &models.ResponseClassDetail{
		ResponseClassBrief: models.ResponseClassBrief{
			ClassId:  cls.ClassId,
			AtcCount: cls.AtcCount,
			Name:     cls.Name,
		},
		Desc:      cls.Desc,
		CreatedAt: cls.CreatedAt,
		UpdatedAt: cls.UpdatedAt,
	}
	return res, nil
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

func GetAllClasses() ([]models.ResponseClassBrief, error) {
	classes, err := mysql.QueryAllClasses()

	if classes == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}

	res := make([]models.ResponseClassBrief, len(classes))
	for idx, v := range classes {
		res[idx] = models.ResponseClassBrief{
			ClassId:  v.ClassId,
			Name:     v.Name,
			AtcCount: v.AtcCount,
		}
	}

	return res, nil
}

// DeleteOneClassById 根据类别ID删除class
func DeleteOneClassById(id int64) error {
	// 查询该类别下是否包含文章 查询出错或者存在文章，不允许删除，返回false
	num, err := mysql.QueryArticleNumberWithClass(id)
	if err != nil {
		return ErrQueryDatabase
	}

	// 存在文章 - 不允许删除
	if num != 0 {
		return ErrDeleteClassByIds
	}
	if err := mysql.DeleteClassById(id); err != nil {
		return ErrQueryDatabase
	}
	return nil
}

// DeleteMultiClassById  返回没有被删除的id列表
func DeleteMultiClassById(ids []int64) ([]int64, error) {
	delIds := make([]int64, 0)
	retIds := make([]int64, 0)
	// 先查找不包含文章的class
	for _, id := range ids {
		num, err := mysql.QueryArticleNumberWithClass(id)
		if err != nil {
			zap.L().Debug("查询数据库出错", zap.Error(err))
			continue
		}

		if num != 0 {
			retIds = append(retIds, id)
		} else {
			delIds = append(delIds, id)
		}
	}
	if len(delIds) == 0 {
		return ids, ErrDeleteClassByIds
	}
	if err := mysql.DeleteClassByIds(delIds); err != nil {
		return ids, ErrQueryDatabase
	}
	return retIds, nil
}

func UpdateClass(class *models.Class) error {
	return mysql.UpdateClassById(class)
}
