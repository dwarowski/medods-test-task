package services

import (
	"errors"
	"os"
	"time"

	dt "github.com/dwarowski/medods-test-task/src/dto"
	"github.com/dwarowski/medods-test-task/src/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func CreateUser(db *gorm.DB, dto dt.CreateUserDto) (any, error) {
	HashedPassword, _ := HashPassword(dto.PlainPassword)
	user := models.User{Username: dto.Username, Email: dto.Email, Password: HashedPassword}

	createUser := db.Create(&user)
	if createUser.Error != nil {
		return nil, createUser.Error
	}

	accessToken, accErr := GenreateAccessToken(user.ID)
	if accErr != nil {
		return nil, accErr
	}
	refreshToken, tokenId, refErr := GenerateRefreshToken(user.ID)
	if refErr != nil {
		return nil, refErr
	}

	hashedToken, tokenErr := HashPassword(tokenId.String())
	if tokenErr != nil {
		return nil, refErr
	}

	addToken := db.Model(&user).Update("refresh_token", hashedToken)
	if addToken.Error != nil {
		return nil, addToken.Error
	}

	tokens := dt.TokensDto{AccessToken: accessToken, RefreshToken: refreshToken}
	return tokens, nil
}

func Login(db *gorm.DB, dto dt.LoginDto) (any, error) {
	var user models.User
	result := db.Table("users").Where("email = ?", dto.Email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.PlainPassword))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, accErr := GenreateAccessToken(user.ID)
	if accErr != nil {
		return nil, accErr
	}
	refreshToken, tokenId, refErr := GenerateRefreshToken(user.ID)
	if refErr != nil {
		return nil, refErr
	}

	hashedToken, tokenErr := HashPassword(tokenId.String())
	if tokenErr != nil {
		return nil, refErr
	}

	addToken := db.Model(&user).Update("refresh_token", hashedToken)
	if addToken.Error != nil {
		return nil, addToken.Error
	}

	tokens := dt.TokensDto{AccessToken: accessToken, RefreshToken: refreshToken}
	return tokens, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ReadRSAKey(path string) (any, error) {
	privateKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func GenreateAccessToken(userId uuid.UUID) (string, error) {
	payload := jwt.MapClaims{
		"id":   uuid.New(),
		"guid": userId,
		"exp":  time.Now().Add(time.Minute * 4).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, payload)

	secret, err := ReadRSAKey("keys/private.pem")
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateRefreshToken(userId uuid.UUID) (string, uuid.UUID, error) {
	tokenId := uuid.New()
	payload := jwt.MapClaims{
		"id":       tokenId,
		"guid":     userId,
		"deviceId": uuid.New(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, payload)

	secret, err := ReadRSAKey("keys/private.pem")
	if err != nil {
		return "", uuid.Nil, err
	}

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", uuid.Nil, err
	}

	return signedToken, tokenId, nil
}
