package models

import "time"

// Community的简略信息
type Community struct {
	Id   int64  `json:"id,string" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

// Community的详细信息
type CommunityDetail struct {
	Id         int64     `json:"id,string" db:"community_id"`
	Name       string    `json:"name" db:"community_name"`
	Desc       string    `json:"desc,omitempty" db:"introduction"` // omitempty表示在序列化时该字段为空就不展示了
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}
