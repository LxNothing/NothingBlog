package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	return mysql.QueryCommunityList()
}

func GetCommunityInfoById(id int64) (c *models.CommunityDetail, err error) {
	return mysql.QueryCommunityById(id)
}

// 查询社区id是否存在
func QueryCommunityExistedById(id int64) bool {
	return mysql.QueryCommunityExistedById(id)
}
