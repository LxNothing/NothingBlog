package mysql

import (
	"NothingBlog/models"
)

func QueryCommunityList() (communityList []*models.Community, err error) {
	// sqlStr := "select community_id, community_name from community"

	// if err = db.Select(&communityList, sqlStr); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		err = ErrInvalidCommunityId
	// 		zap.L().Warn("查找到社区列表数据为空")
	// 	}
	// }
	return nil, nil
}

// 根据社区id查询社区的具体信息
func QueryCommunityById(id int64) (comm *models.CommunityDetail, err error) {
	// sqlStr := `select
	//  		community_id, community_name, introduction, create_time, update_time
	// 		from community
	// 		where community_id = ?`

	// comm = new(models.CommunityDetail)

	// if err = db.Get(comm, sqlStr, id); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		err = ErrInvalidCommunityId
	// 		zap.L().Warn("查找到社区详细数据为空")
	// 	}
	// }
	return nil, nil
}

func QueryCommunityExistedById(id int64) bool {
	// sqlStr := `select 1 from community where community_id = ?`
	// var tmp int64
	// if err := db.Get(&tmp, sqlStr, id); err != nil {
	// 	return false
	// }
	return true
}
