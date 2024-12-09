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

func GetAllTags() ([]models.ResponseTagBrief, error) {
	tags, err := mysql.QueryAllTags()

	if tags == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrTageQueryFailed
		}

	}
	res := make([]models.ResponseTagBrief, len(tags))
	for idx, v := range tags {
		res[idx] = models.ResponseTagBrief{
			TagId:    v.TagId,
			Name:     v.Name,
			AtcCount: int32(v.ArticleCount),
		}
	}

	return res, err
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

func GetTagById(id int64) (*models.Tag, error) {
	tag := new(models.Tag)
	tag.TagId = id
	err := mysql.QueryTagById(tag)

	if err != nil {
		zap.L().Warn("查询数据库出错", zap.Error(err))
		return nil, ErrTageQueryFailed
	}

	return tag, nil
}

// 根据给定的类别名称 获取该类别下所有文章的不重复tag
func GetTagByClassId(clsId int64) ([]models.ResponseTagBrief, error) {
	tags, err := mysql.QuerytTagByClassId(clsId)
	if tags == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrTageQueryFailed
		}
		return nil, nil
	}

	var res = make([]models.ResponseTagBrief, 0, len(tags))

	for _, v := range tags {
		res = append(res, models.ResponseTagBrief{
			TagId:    v.TagId,
			AtcCount: int32(v.ArticleCount),
			Name:     v.Name,
		})
	}

	return res, nil
}

// 删除单个tag
func DeleteTagById(id int64) error {
	ids := make([]int64, 1)
	ids[0] = id
	return mysql.DeleteTagByIds(ids)
}

// 删除多个tag
func DeleteMultiTagById(ids []int64) error {
	return mysql.DeleteTagByIds(ids)
}

func UpdateTag(tag *models.Tag) error {
	return mysql.UpdateTagById(tag)
}
