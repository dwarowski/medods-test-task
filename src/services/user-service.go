package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	dt "github.com/dwarowski/medods-test-task/src/dto"
	"github.com/dwarowski/medods-test-task/src/models"
	"github.com/dwarowski/medods-test-task/src/utils/gentokens"
	"github.com/dwarowski/medods-test-task/src/utils/hashstring"
	"github.com/dwarowski/medods-test-task/src/utils/readkey"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Get tokens for user by its id
func GetByID(db *gorm.DB, id uuid.UUID, userAgent string, ipAdress string) (*dt.TokensDto, error) {

	// Trying to find user
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Generating and saving tokens
	tokens, tokenErr := GenerateAndSaveTokens(db, user.ID, userAgent, ipAdress)
	if tokenErr != nil {
		return nil, tokenErr
	}
	return tokens, nil
}

// Create user from email username and password
func CreateUser(db *gorm.DB, dto dt.CreateUserDto, userAgent string, ipAdress string) (*dt.TokensDto, error) {

	// Hash plain password
	HashedPassword, _ := hashstring.Hash(dto.PlainPassword)
	user := models.User{Username: dto.Username, Email: dto.Email, Password: HashedPassword}

	// Save user to db
	createUser := db.Create(&user)
	if createUser.Error != nil {
		return nil, createUser.Error
	}

	// Generating and saving tokens
	tokens, tokenErr := GenerateAndSaveTokens(db, user.ID, userAgent, ipAdress)
	if tokenErr != nil {
		return nil, tokenErr
	}
	return tokens, nil
}

// Log in via email and password
func Login(db *gorm.DB, dto dt.LoginDto, userAgent string, ipAdress string) (*dt.TokensDto, error) {

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
	tokens, tokenErr := GenerateAndSaveTokens(db, user.ID, userAgent, ipAdress)
	if tokenErr != nil {
		return nil, tokenErr
	}
	return tokens, nil
}

// Tokens generation and save to db
func GenerateAndSaveTokens(db *gorm.DB, userID uuid.UUID, userAgent string, ipAdress string) (*dt.TokensDto, error) {

	// Generate access and refresh token
	accessToken, accessTokenId, accErr := gentokens.GenreateAccessToken(userID)
	if accErr != nil {
		return nil, accErr
	}
	refreshToken, tokenId, refErr := gentokens.GenerateRefreshToken(accessTokenId, userID, userAgent, ipAdress)
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

	return &dt.TokensDto{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Refresh tokens
func RefreshToken(db *gorm.DB, dto dt.TokensDto, userAgent string, ipAdress string) (*dt.TokensDto, error) {

	// Parse refresh token and save payload with defined structure
	refreshToken := &gentokens.RefreshTokenClaims{}
	_, refTknErr := jwt.ParseWithClaims(dto.RefreshToken, refreshToken, func(t *jwt.Token) (any, error) { return readkey.ReadPublicKey() })
	if refTknErr != nil {
		return nil, refTknErr
	}

	// Parse access token and save payload with defined structure
	accessToken := &gentokens.AccessTokenClaims{}
	_, accTknErr := jwt.ParseWithClaims(dto.AccessToken, accessToken, func(t *jwt.Token) (any, error) { return readkey.ReadPublicKey() })
	if accTknErr != nil {
		fmt.Printf("\x1b[43mWarning: %v\x1b[0m\n ", accTknErr)
	}

	// Check if tokens are pair
	if accessToken.ID != refreshToken.ATID.String() {
		return nil, errors.New("token id mismatch")
	}

	// Get user from db
	userModel := models.User{ID: refreshToken.GUID}
	findRes := db.First(&userModel)
	if findRes.Error != nil {
		return nil, findRes.Error
	}

	// compare refresh token in db
	isTokenValid := bcrypt.CompareHashAndPassword([]byte(userModel.RefreshToken), []byte(refreshToken.ID))
	if isTokenValid != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Deny refresh from different devices
	if userAgent != refreshToken.UserAgent {
		removeToken := db.Model(&models.User{}).Where("id = ?", refreshToken.GUID).Update("refresh_token", "")
		if removeToken.Error != nil {
			return nil, removeToken.Error
		}
		return nil, errors.New("can't refresh from different device")
	}

	// Warn if ip changed
	if ipAdress != refreshToken.IpAdress {
		payload := dt.SendToWebhookDto{
			UID:       refreshToken.GUID,
			IpAddress: ipAdress,
			UserAgent: userAgent,
			Timestamp: time.Now(),
		}
		jsonData, _ := json.Marshal(payload)
		webhookURL := "URL"
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error sending webhook: %v", err)
		} else {
			defer resp.Body.Close()
			fmt.Printf("Webhook sent successfully. Status code: %d", resp.StatusCode)
		}

	}

	// Refresh tokens
	tokens, tokenErr := GenerateAndSaveTokens(db, refreshToken.GUID, userAgent, ipAdress)
	if tokenErr != nil {
		return nil, tokenErr
	}
	return tokens, nil
}

// Get User id
func GetUUID(db *gorm.DB, accessToken string) (*dt.GetUUIDDto, error) {

	// Parse access token and save payload with defined structure
	payload := &gentokens.AccessTokenClaims{}
	_, parseErr := jwt.ParseWithClaims(accessToken, payload, func(t *jwt.Token) (any, error) { return readkey.ReadPublicKey() })
	if parseErr != nil {
		return nil, parseErr
	}

	// Find user and return id
	var user models.User
	result := db.First(&user, payload.GUID)
	if result.Error != nil {
		return nil, result.Error
	}
	userId := &dt.GetUUIDDto{Uuid: user.ID}
	return userId, nil
}
