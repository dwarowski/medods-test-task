package services

import (
	models "github.com/dwarowski/medods-test-task/src/model"
	"gorm.io/gorm"
)

func GetByID(db *gorm.DB, id int) (any, error) {
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
