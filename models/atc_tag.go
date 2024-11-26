package models

import "gorm.io/gorm"

// 描述文章标签的数据模型
type AtcTag struct {
	gorm.Model
	Name         string    `gorm:"type:varchar(100);not null;" json:"name"`           // 标签名称
	Desc         string    `gorm:"type:varchar(100);" json:"desc"`                    // 标签的描述信息
	ArticleList  []Article `gorm:"many2many:tag_article" json:"article_list"`         // 文章列表
	ArticleCount uint32    `gorm:"type:int;not null;default 0;" json:"article_count"` // 该标签下拥有的文章数量
}
