package services

import (
	"errors"

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
	user := models.User{Username: dto.Username, Email: dto.Email, Password: HashedPassword}

	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func Login(db *gorm.DB, dto dto.LoginDto) (any, error) {
	var user models.User
	result := db.Table("users").Where("email = ?", dto.Email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.PlainPassword))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return true, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
