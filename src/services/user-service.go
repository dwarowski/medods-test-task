package services

import (
	"errors"

	dt "github.com/dwarowski/medods-test-task/src/dto"
	"github.com/dwarowski/medods-test-task/src/models"
	"github.com/dwarowski/medods-test-task/src/utils/gentokens"
	"github.com/dwarowski/medods-test-task/src/utils/hashstring"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetByID(db *gorm.DB, id uuid.UUID, userAgent string) (any, error) {
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Generating and saving tokens
	tokens, tokenErr := GenerateAndSaveTokens(db, user.ID, userAgent)
	if tokenErr != nil {
		return nil, tokenErr
	}
	return tokens, nil
}

func CreateUser(db *gorm.DB, dto dt.CreateUserDto, userAgent string) (any, error) {

	// Hash plain password
	HashedPassword, _ := hashstring.Hash(dto.PlainPassword)
	user := models.User{Username: dto.Username, Email: dto.Email, Password: HashedPassword}

	// Save user to db
	createUser := db.Create(&user)
	if createUser.Error != nil {
		return nil, createUser.Error
	}

	// Generating and saving tokens
	tokens, tokenErr := GenerateAndSaveTokens(db, user.ID, userAgent)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return tokens, nil
}

func Login(db *gorm.DB, dto dt.LoginDto, userAgent string) (any, error) {

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

	// Generating and saving tokens
	tokens, tokenErr := GenerateAndSaveTokens(db, user.ID, userAgent)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return tokens, nil
}

func GenerateAndSaveTokens(db *gorm.DB, userID uuid.UUID, userAgent string) (any, error) {

	// Generate access and refresh token
	accessToken, accessTokenId, accErr := gentokens.GenreateAccessToken(userID)
	if accErr != nil {
		return nil, accErr
	}
	refreshToken, tokenId, refErr := gentokens.GenerateRefreshToken(accessTokenId, userID, userAgent)
	if refErr != nil {
		return nil, refErr
	}

	// Hash token Id to check later
	hashedToken, tokenErr := hashstring.Hash(tokenId.String())
	if tokenErr != nil {
		return nil, refErr
	}

	// Save hashed token Id to db
	addToken := db.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", hashedToken)
	if addToken.Error != nil {
		return nil, addToken.Error
	}

	return dt.TokensDto{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
