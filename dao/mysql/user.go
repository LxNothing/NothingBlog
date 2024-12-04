package mysql

import (
	"NothingBlog/models"
	"errors"

	"gorm.io/gorm"
)

// 查询所有的用户信息
func QueryAllUser() ([]models.User, error) {
	var users []models.User
	if err := Db.Select("user_id, user_name, email, user_icon, desc").Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return users, nil
}

// 根据名字查询用户是否已经存在 不存在返回nil 存在或者查询出错返回 error
// 注意gorm的版本 本代码是基于V2版本的gorm，不能用V1版本，否则这个代码会有问题
func QueryUserByName(user *models.User) (err error) {
	// 注 Find不会产生 ErrRecordNotFound
	err = Db.Where("user_name=?", user.UserName).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return ErrUserExisted
}

// 设置用户的密码
func SetUserPassword(userName string, newPwd string) error {
	return Db.Model(&models.User{}).Where("`user_name` = ?", userName).Update("password", newPwd).Error
}

// 根据用户的名称和邮箱查找用户
// func QueryUserByNameAndEmail(user *models.User) (err error) {
// 	// 注 Find不会产生 ErrRecordNotFound
// 	err = Db.Where("user_name=?", user.UserName).Take(user).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil
// 	}
// 	return ErrUserExisted
// }

// 插入用户数据
func InsertUser(u *models.User) (err error) {
	return Db.Create(u).Error
}

// 根据用户ID查询用户名称
func QueryUsernameById(id int64) (name string, err error) {
	// sqlStr := "select username from user where user_id = ?"

	// err = db.Get(&name, sqlStr, id)

	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		err = ErrorUserNotExist
	// 	}
	// 	return
	// }
	return "", nil
}
