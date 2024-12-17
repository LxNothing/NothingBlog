package models

import (
	"time"

	"gorm.io/gorm"
)

// 用于gorm生成数据表使用
type Class struct {
	gorm.Model
	ClassId  int64  `json:"class_id,string" gorm:"type:bigint(20);unique;not null"` // 类别ID - 由应用层生成
	AtcCount int32  `json:"atc_count" gorm:"type:int;not null;default 0"`           // 该分类下包含的文章数量
	Name     string `json:"name" gorm:"type:varchar(256); unique; not null"`
	Desc     string `json:"desc" gorm:"type:varchar(256)"`
}

type ClassBriefReturn struct {
	ClassId  int64  `json:"class_id,string"` // 类别ID - 由应用层生成
	AtcCount int32  `json:"atc_count"`       // 该分类下包含的文章数量
	Name     string `json:"name"`            // 类别名称
}

// 类别class、的详细信息
type ClassEntireReturn struct {
	ClassBriefReturn
	Desc      string    `json:"desc"`       // 类别的描述信息
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

func (c *Class) BindToBriefClass() *ClassBriefReturn {
	return &ClassBriefReturn{
		ClassId:  c.ClassId,
		AtcCount: c.AtcCount,
		Name:     c.Name,
	}
}

func (c *Class) BindToEntireClass() *ClassEntireReturn {
	return &ClassEntireReturn{
		ClassBriefReturn: *c.BindToBriefClass(),
		Desc:             c.Desc,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.CreatedAt,
	}
}
