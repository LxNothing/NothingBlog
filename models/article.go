package models

import "gorm.io/gorm"

// 文章的数据模型

type Article struct {
	gorm.Model
	AuthorId     uint64   `gorm:"type:bigint(20);not null;" json:"author_id"` // 作者id
	TagList      []AtcTag `gorm:"many2many:tag_article" json:"tag_list"`      // 文章标签
	TopFlag      bool     `gorm:"type:bool;default:false;" json:"top_flag"`   // 是否置顶标志 false-没有置顶 true-置顶
	Status       uint8    `gorm:"type:int;not null;default:0;" json:"status"` // 文章状态 0-草稿 1-已提交 2-审核完成
	Title        string   `gorm:"type:varchar(255);not null;" json:"title"`   // 文章标题
	Summary      string   `gorm:"type:varchar(255);" json:"summary"`          // 文章的摘要信息
	Image        string   `gorm:"type:varchar(255);" json:"image"`            // 文章缩略图
	Content      string   `gorm:"type:MediumText;" json:"content"`            // 文章内容
	CommentCount uint32   `gorm:"type:int;default:0;" json:"comment_count"`   // 文章的评论数
	VisitCount   uint32   `gorm:"type:int;default:0;" json:"visit_count"`     // 文章的浏览量
}
