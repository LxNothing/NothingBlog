package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidStatus    = errors.New("状态无效")
	ErrInvalidPrivilege = errors.New("权限无效")
)

// 文章的数据模型 - gorm 迁移数据库时使用
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
	ArticleId    int64  `gorm:"type:bigint(20);unique;not null;" json:"article_id"` // 文章ID - 由应用层雪花算法生成
	AuthorId     int64  `json:"author_id"`                                          // 作者id
	User         User   `gorm:"foreignKey:AuthorId;references:UserId"`
	ClassId      int64  `json:"class_id"`                                                                                                                  // 文章分类ID
	Class        Class  `gorm:"references:ClassId"`                                                                                                        // ClassId同Class中的键名称相同，就不能指定foreignKey字段
	TagList      []Tag  `gorm:"many2many:tag_article;foreignKey:ArticleId;joinForeignKey:ArticleId;References:TagId;joinReferences:TagId" json:"tag_list"` // 文章标签
	TopFlag      bool   `gorm:"type:bool;default:false;" json:"top_flag"`                                                                                  // 是否置顶标志 false-没有置顶 true-置顶
	EnComment    bool   `gorm:"type:bool;default:false;" json:"en_comment"`                                                                                // 是否允许评论 false-不允许 true-允许
	Status       uint8  `gorm:"type:int;not null;default:1;" json:"status"`                                                                                // 文章状态 1-草稿 2-已提交 3-删除(预删除，实际数据没有删)
	Privilege    uint8  `gorm:"type:int;not bull;default:1" json:"privilege"`                                                                              // 文章权限 1-公开 2-私有
	Title        string `gorm:"type:varchar(255);not null;unique" json:"title"`                                                                            // 文章标题,不允许重复
	Summary      string `gorm:"type:varchar(255);" json:"summary"`                                                                                         // 文章的摘要信息
	Image        string `gorm:"type:varchar(255);" json:"image"`                                                                                           // 文章缩略图，应该存储一个URL
	Content      string `gorm:"type:MediumText;" json:"content"`                                                                                           // 文章内容
	CommentCount uint32 `gorm:"type:int;default:0;" json:"comment_count"`                                                                                  // 文章的评论数
	VisitCount   uint32 `gorm:"type:int;default:0;" json:"visit_count"`                                                                                    // 文章的浏览量
}

// 返回文章的简略信息
type ArticleBriefReturn struct {
	CreatedAt    time.Time        `json:"created_at"`    // 文章的创建时间
	UpdatedAt    time.Time        `json:"updated_at"`    // 文章的最后修改时间
	ArticleId    int64            `json:"article_id"`    // 文章ID - 由应用层雪花算法生成
	Author       UserBrief        `json:"author"`        // 只需要返回简略的作者即可
	Class        ClassBriefReturn `json:"class"`         // 简略的类别信息
	Tag          []TagBriefReturn `json:"tag"`           // 简略的tag信息
	TopFlag      bool             `json:"top_flag"`      // 是否置顶标志 false-没有置顶 true-置顶
	EnComment    bool             `json:"en_comment"`    // 是否允许评论 false-不允许 true-允许
	Status       uint8            `json:"status"`        // 文章状态 1-草稿 2-已提交 3-删除(预删除，实际数据没有删)
	Privilege    uint8            `json:"privilege"`     // 文章权限 1-公开 2-私有
	Title        string           `json:"title"`         // 文章标题,不允许重复
	Summary      string           `json:"summary"`       // 文章的摘要信息
	Image        string           `json:"image"`         // 文章缩略图
	CommentCount uint32           `json:"comment_count"` // 文章的评论数
	VisitCount   uint32           `json:"visit_count"`   // 文章的浏览量
}

// 返回文章的详细数据 - 主要包括文章的具体内容
type ArticleEntireReturn struct {
	ArticleBriefReturn
	Content string `json:"content"` // 文章的具体内容
}

func (a *Article) BindToBriefArticle() *ArticleBriefReturn {
	res := ArticleBriefReturn{
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
		ArticleId:    a.ArticleId,
		Tag:          make([]TagBriefReturn, 0),
		TopFlag:      a.TopFlag,
		EnComment:    a.EnComment,
		Status:       a.Status,
		Privilege:    a.Privilege,
		Title:        a.Title,
		Summary:      a.Summary,
		Image:        a.Image,
		CommentCount: a.CommentCount,
		VisitCount:   a.VisitCount,
	}

	res.Author = a.User.BindToBriefUser()
	res.Class = *a.Class.BindToBriefClass()
	for _, v := range a.TagList {
		res.Tag = append(res.Tag, *v.BindToBriefTag())
	}
	return &res
}

func (a *Article) BindToEntireArticle() *ArticleEntireReturn {
	res := new(ArticleEntireReturn)
	res.ArticleBriefReturn = *(a.BindToBriefArticle())
	res.Content = a.Content
	return res
}

// 权限的常量定义
const (
	AllPrivilege = "全部"
	Public       = "公开"
	Private      = "私有"
)

// 文章状态的常量定义
const (
	AllStatus = "全部"
	Drift     = "草稿"
	Commit    = "已发布"
	Recycle   = "回收站"
)

func PrivilegeStringToNumber(priv string) (uint8, error) {
	var strToUintStatus = map[string]uint8{
		"全部": 0,
		"公开": 1,
		"私有": 2,
	}
	num, ok := strToUintStatus[priv]
	if !ok {
		return 0, ErrInvalidStatus
	}
	return num, nil
}

func PrivilegeNumberToString(priv uint8) (string, error) {
	var uintToStrStatus = map[uint8]string{
		0: "全部",
		1: "草稿",
	}
	res, ok := uintToStrStatus[priv]
	if !ok {
		return "", ErrInvalidStatus
	}
	return res, nil
}

func StatusStringToNumber(stus string) (uint8, error) {
	var strToUintStatus = map[string]uint8{
		AllStatus: 0,
		Drift:     1,
		Commit:    2,
		Recycle:   3,
	}
	num, ok := strToUintStatus[stus]
	if !ok {
		return 0, ErrInvalidStatus
	}
	return num, nil
}

func StatusNumberToString(stus uint8) (string, error) {
	var uintToStrStatus = map[uint8]string{
		0: AllStatus,
		1: Drift,
		2: Commit,
		3: Recycle,
	}
	res, ok := uintToStrStatus[stus]
	if !ok {
		return "", ErrInvalidStatus
	}
	return res, nil
}
