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

func QueryClassesById(id int64) (*models.Class, error) {
	class := new(models.Class)

	if err := Db.Where("class_id = ?", id).Find(class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return class, nil
}

func QueryClassesByName(name string) (*models.Class, error) {
	class := new(models.Class)

	if err := Db.Where("name = ?", name).Find(class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return class, nil
}

// 注意：文章表中的class_id关联到这class表，所以不能直接删除本表中的值
// 业务实现是：如果有文章属于这个类别，那么这个类别不允许删除，只能更改
// 因此这里没有更新文章表中的class_id
func DeleteClassByIds(ids []int64) error {
	return Db.Where("class_id in (ids)").Unscoped().Delete(&models.Class{}).Error
}

func UpdateClassById(class *models.Class) error {
	return Db.Select("updated_at", "name", "desc").Where("class_id = ?", class.ClassId).Updates(class).Error
}
