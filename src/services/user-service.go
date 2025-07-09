package services

import (
	"errors"
	"time"

	dt "github.com/dwarowski/medods-test-task/src/dto"
	"github.com/dwarowski/medods-test-task/src/models"
	"github.com/dwarowski/medods-test-task/src/utils/hashstring"
	"github.com/dwarowski/medods-test-task/src/utils/readkey"
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

	// Hash plain password
	HashedPassword, _ := hashstring.Hash(dto.PlainPassword)
	user := models.User{Username: dto.Username, Email: dto.Email, Password: HashedPassword}

	// Save user to db
	createUser := db.Create(&user)
	if createUser.Error != nil {
		return nil, createUser.Error
	}

	// Generate access and refresh token
	accessToken, accErr := GenreateAccessToken(user.ID)
	if accErr != nil {
		return nil, accErr
	}
	refreshToken, tokenId, refErr := GenerateRefreshToken(user.ID)
	if refErr != nil {
		return nil, refErr
	}

	// Hash token Id to check later
	hashedToken, tokenErr := hashstring.Hash(tokenId.String())
	if tokenErr != nil {
		return nil, refErr
	}

	// Save hashed token Id to db
	addToken := db.Model(&user).Update("refresh_token", hashedToken)
	if addToken.Error != nil {
		return nil, addToken.Error
	}

	return dt.TokensDto{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func Login(db *gorm.DB, dto dt.LoginDto) (any, error) {

	// Check if user with this email exsist
	var user models.User
	result := db.Table("users").Where("email = ?", dto.Email).First(&user)
	if result.Error != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if password is correct
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.PlainPassword))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access and refresh token
	accessToken, accErr := GenreateAccessToken(user.ID)
	if accErr != nil {
		return nil, accErr
	}
	refreshToken, tokenId, refErr := GenerateRefreshToken(user.ID)
	if refErr != nil {
		return nil, refErr
	}

	// Hash token Id to check later
	hashedToken, tokenErr := hashstring.Hash(tokenId.String())
	if tokenErr != nil {
		return nil, refErr
	}

	// Save hashed token Id to db
	addToken := db.Model(&user).Update("refresh_token", hashedToken)
	if addToken.Error != nil {
		return nil, addToken.Error
	}

	return dt.TokensDto{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func GenreateAccessToken(userId uuid.UUID) (string, error) {
	payload := jwt.MapClaims{
		"id":   uuid.New(),
		"guid": userId,
		"exp":  time.Now().Add(time.Minute * 4).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, payload)

	secret, err := readkey.ReadRSAKey("keys/private.pem")
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

	secret, err := readkey.ReadRSAKey("keys/private.pem")
	if err != nil {
		return "", uuid.Nil, err
	}

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", uuid.Nil, err
	}

	return signedToken, tokenId, nil
}
