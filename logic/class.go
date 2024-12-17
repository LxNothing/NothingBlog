package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
	"errors"
	"time"

	"go.uber.org/zap"
)

type LogicClass struct {
}

var (
	ErrDeleteClassByIds = errors.New("该类别中存在文章")
	ErrQueryDatabase    = errors.New("查询数据库出错")
)

var daoClass mysql.DaoClass

// CreateNewClass 创建新的类别，并返回类别的名称和ID
func (lc LogicClass) CreateNewClass(cls *models.Class) (*models.ClassBriefReturn, error) {
	cls.ClassId = snowflake.GetNextId().Int64()
	cls.CreatedAt = time.Now()
	cls.UpdatedAt = cls.CreatedAt

	if err := daoClass.InsertNewClass(cls); err != nil {
		zap.L().Debug("插入新的类别到数据库失败", zap.Error(err))
		return nil, err
	}

	new_class := &models.ClassBriefReturn{
		ClassId:  cls.ClassId,
		Name:     cls.Name,
		AtcCount: 0,
	}
	return new_class, nil
}

func (lc LogicClass) GetClassById(id int64) (*models.ClassEntireReturn, error) {
	cls, err := daoClass.QueryClassesById(id)

	if cls == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	res := cls.BindToEntireClass()
	return res, nil
}

func (lc LogicClass) GetClassByName(name string) (*models.Class, error) {
	cls, err := daoClass.QueryClassesByName(name)

	if cls == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}
	return cls, nil
}

func (lc LogicClass) GetAllClasses() ([]models.ClassBriefReturn, error) {
	classes, err := daoClass.QueryAllClasses()

	if classes == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrArticleQueryFailed
		}
	}

	res := make([]models.ClassBriefReturn, len(classes))
	for idx, v := range classes {
		res[idx] = *v.BindToBriefClass()
	}

	return res, nil
}

// DeleteOneClassById 根据类别ID删除class
func (lc LogicClass) DeleteOneClassById(id int64) error {
	// 查询该类别下是否包含文章 查询出错或者存在文章，不允许删除，返回false
	num, err := articleDb.QueryArticleNumberWithClass(id)
	if err != nil {
		return ErrQueryDatabase
	}

	// 存在文章 - 不允许删除
	if num != 0 {
		return ErrDeleteClassByIds
	}
	if err := daoClass.DeleteClassById(id); err != nil {
		return ErrQueryDatabase
	}
	return nil
}

// DeleteMultiClassById  返回没有被删除的id列表
func (lc LogicClass) DeleteMultiClassById(ids []int64) ([]int64, error) {
	delIds := make([]int64, 0)
	retIds := make([]int64, 0)
	// 先查找不包含文章的class
	for _, id := range ids {
		num, err := articleDb.QueryArticleNumberWithClass(id)
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
	if err := daoClass.DeleteClassByIds(delIds); err != nil {
		return ids, ErrQueryDatabase
	}
	return retIds, nil
}

func (lc LogicClass) UpdateClass(class *models.Class) error {
	return daoClass.UpdateClassById(class)
}
