package models

import (
	"time"

	"gorm.io/gorm"
)

type PrivilegeType uint8
type StatusType uint8

// 文章的数据模型
/*
	对于TagList []Tag，定义了多对多的关系，即文章包含多个tag，tag也可以包含多个文章，因此many2many则是指定这种关系，并
	gorm:"
	many2many:tag_article;   指定生成的中间表为tag_article
	foreignKey:ArticleId;    指定本表中的外键字段ArticleId
	joinForeignKey:ArticleId; 上面的ArticleId关联到中间表的ArticleId字段
	References:TagId; 指定第三个表，即tags表中的关联字段tag_id
	joinReferences:TagId 上面的 tag_id 关联到中间表的tag_id字段
	gorm会自动将 TagId 这种格式转为 tag_id这种格式
*/
type Article struct {
	gorm.Model
	ArticleId    int64         `gorm:"type:bigint(20);unique;not null;" json:"article_id"` // 文章ID - 由应用层雪花算法生成
	AuthorId     int64         `json:"author_id"`                                          // 作者id
	User         User          `gorm:"foreignKey:AuthorId;references:UserId"`
	ClassId      int64         `json:"class_id"`                                                                                                                  // 文章分类ID
	Class        Class         `gorm:"references:ClassId"`                                                                                                        // ClassId同Class中的键名称相同，就不能指定foreignKey字段
	TagList      []Tag         `gorm:"many2many:tag_article;foreignKey:ArticleId;joinForeignKey:ArticleId;References:TagId;joinReferences:TagId" json:"tag_list"` // 文章标签
	TopFlag      bool          `gorm:"type:bool;default:false;" json:"top_flag"`                                                                                  // 是否置顶标志 false-没有置顶 true-置顶
	EnComment    bool          `gorm:"type:bool;default:false;" json:"en_comment"`                                                                                // 是否允许评论 false-不允许 true-允许
	Status       StatusType    `gorm:"type:int;not null;default:0;" json:"status"`                                                                                // 文章状态 0-草稿 1-已提交 2-删除(预删除，实际数据没有删)
	Privilege    PrivilegeType `gorm:"type:int;not bull;default:0" json:"privilege"`                                                                              // 文章权限 0-公开 1-私有
	Title        string        `gorm:"type:varchar(255);not null;unique" json:"title"`                                                                            // 文章标题,不允许重复
	Summary      string        `gorm:"type:varchar(255);" json:"summary"`                                                                                         // 文章的摘要信息
	Image        string        `gorm:"type:varchar(255);" json:"image"`                                                                                           // 文章缩略图，应该存储一个URL
	Content      string        `gorm:"type:MediumText;" json:"content"`                                                                                           // 文章内容
	CommentCount uint32        `gorm:"type:int;default:0;" json:"comment_count"`                                                                                  // 文章的评论数
	VisitCount   uint32        `gorm:"type:int;default:0;" json:"visit_count"`                                                                                    // 文章的浏览量
}

// 返回文章的简略信息
type ArticleBriefReturn struct {
	CreatedAt    time.Time     `json:"created_at"`    // 文章的创建时间
	UpdatedAt    time.Time     `json:"updated_at"`    // 文章的最后修改时间
	ArticleId    int64         `json:"article_id"`    // 文章ID - 由应用层雪花算法生成
	AuthorId     int64         `json:"author_id"`     // 作者id
	AuthorName   string        `json:"author_name"`   // 作者名字
	ClassId      int64         `json:"class_id"`      // 文章分类ID
	ClassName    string        `json:"class_name"`    // 所属类别名
	TagId        []int64       `json:"tag_id"`        // 所属tag的ID列表
	TagName      []string      `json:"tag_name"`      // 所属tag的Name
	TopFlag      bool          `json:"top_flag"`      // 是否置顶标志 false-没有置顶 true-置顶
	EnComment    bool          `json:"en_comment"`    // 是否允许评论 false-不允许 true-允许
	Status       StatusType    `json:"status"`        // 文章状态 0-草稿 1-已提交 2-删除(预删除，实际数据没有删)
	Privilege    PrivilegeType `json:"privilege"`     // 文章权限 0-公开 1-私有
	Title        string        `json:"title"`         // 文章标题,不允许重复
	Summary      string        `json:"summary"`       // 文章的摘要信息
	Image        string        `json:"image"`         // 文章缩略图
	CommentCount uint32        `json:"comment_count"` // 文章的评论数
	VisitCount   uint32        `json:"visit_count"`   // 文章的浏览量
}

// 返回文章的详细数据 - 主要包括文章的具体内容
type ArticleEntireReturn struct {
	ArticleBriefReturn
	Content string `json:"content"` // 文章的具体内容
}

const (
	PrivilegePublic PrivilegeType = iota // 0 - 当前文章是公开的，客户端可见（默认）
	PrivilegePrivte                      // 1- 当前文章私有，客户端不可见
	PrivilegeAll                         // 所有权限
)

const (
	StatusDraft  StatusType = iota // 0 - 草稿
	StatusCommit                   // 1 - 已经提交，外部可以访问（同时Privilege要为公开才可以）
	StatusDelete                   // 2 - 预删除，处于这种状态外部不可见，但是数据仍然在数据库（类似垃圾回收站）
	StatusAll                      // 所有状态
)
