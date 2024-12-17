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

type LogicTag struct {
}

var daoTag mysql.DaoTag

func (lt LogicTag) CreateNewTag(tag *models.Tag) error {
	tag.TagId = snowflake.GetNextId().Int64()

	return daoTag.InsertNewTag(tag)
}

func (lt LogicTag) GetAllTags() ([]models.TagBriefReturn, error) {
	tags, err := daoTag.QueryAllTags()

	if tags == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrTageQueryFailed
		}

	}
	res := make([]models.TagBriefReturn, len(tags))
	for idx, v := range tags {
		res[idx] = *v.BindToBriefTag()
	}

	return res, err
}

func (lt LogicTag) GetTagByName(name string) (*models.Tag, error) {
	tag := new(models.Tag)
	tag.Name = name
	err := daoTag.QueryTagByName(tag)

	if err != nil {
		zap.L().Warn("查询数据库出错", zap.Error(err))
		return nil, ErrTageQueryFailed
	}

	return tag, nil
}

func (lt LogicTag) GetTagById(id int64) (*models.TagEntireReturn, error) {
	tag := new(models.Tag)
	tag.TagId = id
	err := daoTag.QueryTagById(tag)

	if err != nil {
		zap.L().Warn("查询数据库出错", zap.Error(err))
		return nil, ErrTageQueryFailed
	}

	res := tag.BindToEntireTag()
	return res, nil
}

// 根据给定的类别名称 获取该类别下所有文章的不重复tag
func (lt LogicTag) GetTagByClassId(clsId int64) ([]models.TagBriefReturn, error) {
	tags, err := daoTag.QuerytTagByClassId(clsId)
	if tags == nil {
		if err != nil {
			zap.L().Warn("查询数据库出错", zap.Error(err))
			return nil, ErrTageQueryFailed
		}
		return nil, nil
	}

	var res = make([]models.TagBriefReturn, 0, len(tags))

	for _, v := range tags {
		res = append(res, *v.BindToBriefTag())
	}

	return res, nil
}

// 删除单个tag
func (lt LogicTag) DeleteTagById(id int64) error {
	ids := make([]int64, 1)
	ids[0] = id
	return daoTag.DeleteTagByIds(ids)
}

// 删除多个tag
func (lt LogicTag) DeleteMultiTagById(ids []int64) error {
	return daoTag.DeleteTagByIds(ids)
}

func (lt LogicTag) UpdateTag(tag *models.Tag) error {
	return daoTag.UpdateTagById(tag)
}
