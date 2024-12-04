package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrTageQueryFailed = errors.New("查询Tag数据库出错")
)

func CreateNewTag(tag *models.Tag) error {
	tag.TagId = snowflake.GetNextId().Int64()

	return mysql.InsertNewTag(tag)
}

func GetAllTags() ([]models.Tag, error) {
	tags, err := mysql.QueryAllTags()

	if tags == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrTageQueryFailed
		}

	}
	return tags, err
}

func GetTagByName(name string) (*models.Tag, error) {
	tag := new(models.Tag)
	tag.Name = name
	err := mysql.QueryTagByName(tag)

	if err != nil {
		zap.L().Warn("查询数据库出错", zap.Error(err))
		return nil, ErrTageQueryFailed
	}

	return tag, nil
}

func GetTagByClassId(clsId int64) ([]models.Tag, error) {
	tags, err := mysql.QuerytTagByClassId(clsId)
	if tags == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrTageQueryFailed
		}
		return nil, nil
	}

	var res = make([]models.Tag, 0, len(tags))

	for _, v := range tags {
		res = append(res, v)
	}

	return res, nil
}
