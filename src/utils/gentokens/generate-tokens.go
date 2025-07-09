package gentokens

import (
	"time"

	"github.com/dwarowski/medods-test-task/src/utils/readkey"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenreateAccessToken(userId uuid.UUID) (string, uuid.UUID, error) {
	tokenId := uuid.New()
	payload := jwt.MapClaims{
		"id":   tokenId,
		"guid": userId,
		"exp":  time.Now().Add(time.Minute * 4).Unix(),
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

func GenerateRefreshToken(accessTokenId uuid.UUID, userId uuid.UUID) (string, uuid.UUID, error) {
	tokenId := uuid.New()
	payload := jwt.MapClaims{
		"id":       tokenId,
		"atid":     accessTokenId,
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
