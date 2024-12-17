package models

import (
	"time"

	"gorm.io/gorm"
)

// 描述文章标签的数据模型
type Tag struct {
	gorm.Model
	TagId        int64     `json:"tag_id" gorm:"type:bigint(20);not null;unique"`     // tag id 由应用层生成
	Name         string    `json:"name" gorm:"type:varchar(100);not null;unique"`     // 标签名称
	Desc         string    `json:"desc" gorm:"type:varchar(100);"`                    // 标签的描述信息
	ArticleList  []Article `json:"article_list" gorm:"many2many:tag_article"`         // 反向引用 文章列表
	ArticleCount uint32    `json:"article_count" gorm:"type:int;not null;default 0;"` // 该标签下拥有的文章数量
}

// 返回的 tag 简略信息
type TagBriefReturn struct {
	TagId        int64  `json:"tag_id,string"` // tag id 由应用层生成
	Name         string `json:"name"`          // 标签名称
	ArticleCount uint32 `json:"article_count"` // 该标签下拥有的文章数量
}

func (t *Tag) BindToBriefTag() *TagBriefReturn {
	return &TagBriefReturn{
		TagId:        t.TagId,
		Name:         t.Name,
		ArticleCount: t.ArticleCount,
	}
}

type TagEntireReturn struct {
	TagBriefReturn
	Desc      string    `json:"desc"`      // 标签的描述信息
	CreatedAt time.Time `json:"create_at"` // 该tag的创建时间
	UpdatedAt time.Time `json:"update_at"` // 该tag的修改时间
}

func (t *Tag) BindToEntireTag() *TagEntireReturn {
	return &TagEntireReturn{
		TagBriefReturn: *t.BindToBriefTag(),
		Desc:           t.Desc,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
