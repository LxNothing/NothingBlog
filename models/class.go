package models

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	ClassId  int64  `json:"class_id,string" gorm:"type:bigint(20);unique;not null"` // 类别ID - 由应用层生成
	AtcCount int32  `json:"atc_count" gorm:"type:int;not null;default 0"`           // 该分类下包含的文章数量
	Name     string `json:"name" gorm:"type:varchar(256)"`
	Desc     string `json:"desc" gorm:"type:varchar(256)"`
}