package models

import "time"

// 创建博文的时候数据结构
type BlogsArch struct {
	Id           int64     `json:"id,string" db:"blog_id"`
	AuthorId     int64     `json:"author_id,string" db:"author_id"`
	CommmunityId int64     `json:"community_id,string" db:"community_id" binding:"required"`
	VoteScore    int64     `json:"vote_score,string" db:"vote_score"` // 添加投票分数
	Status       int32     `json:"status" db:"status" `
	Title        string    `json:"title" db:"title"  binding:"required"`
	Content      string    `json:"content" db:"content"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	UpdateTime   time.Time `json:"update_time" db:"update_time"`
}

type ApiBlogDetail struct {
	AuthorName string `json:"username" db:"username"`
	Blog       *BlogsArch
	Community  *CommunityDetail
}
