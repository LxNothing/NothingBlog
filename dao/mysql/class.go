package mysql

import (
	"NothingBlog/models"
	"errors"

	"gorm.io/gorm"
)

// 向数据库中插入数据
func InsertNewClass(cls *models.Class) error {
	return Db.Create(cls).Error
}

func QueryAllClasses() ([]models.Class, error) {
	var class []models.Class
	if err := Db.Find(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return class, nil
}

func QueryClassesByName(name string) (*models.Class, error) {
	class := new(models.Class)

	if err := Db.Where("name=?", name).Find(class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return class, nil
}
