package services

import (
	"fmt"

	"github.com/dwarowski/medods-test-task/src/dto"
	"github.com/dwarowski/medods-test-task/src/models"
	"golang.org/x/crypto/bcrypt"
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

func CreateUser(db *gorm.DB, dto dto.CreateUserDto) (any, error) {
	HashedPassword, _ := HashPassword(dto.PlainPassword)
	fmt.Println(HashedPassword)
	user := models.User{Username: dto.Username, Email: dto.Email, Password: HashedPassword}

	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
