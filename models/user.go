package models

import (
	"github.com/jinzhu/gorm"
)

// 描述用户的数据模型
type User struct {
	gorm.Model
	UserId   int64  `json:"user_id" gorm:"type:bigint(20);unique;not null"`   // 用户ID，由应用层雪花算法计算得出
	UserName string `json:"username" gorm:"type:varchar(30);unique;not null"` // 用户名
	Email    string `json:"email" gorm:"type:varchar(30)"`                    // 邮箱
	Password string `json:"password" gorm:"type:varchar(100);not null"`       // 密码
	UserIcon string `json:"user_icon" gorm:"type:varchar(255)"`               // 头像
	Desc     string `json:"desc" gorm:"type:varchar(255)"`                    // 个人描述-
}

type UserBrief struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"username"`
}

func (u *User) BindToBriefUser() UserBrief {
	return UserBrief{
		UserId:   u.UserId,
		UserName: u.UserName,
	}
}
