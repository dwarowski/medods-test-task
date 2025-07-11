package gentokens

import (
	"time"

	"github.com/dwarowski/medods-test-task/src/utils/readkey"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	GUID uuid.UUID `json:"guid"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	ATID      uuid.UUID `json:"atid"` // Access Token ID
	GUID      uuid.UUID `json:"guid"`
	UserAgent string    `json:"userAgent"`
	IpAdress  string    `json:"ipadress"`
	jwt.RegisteredClaims
}

func GenreateAccessToken(userId uuid.UUID) (string, uuid.UUID, error) {
	tokenId := uuid.New()
	payload := AccessTokenClaims{
		GUID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 4)), // Set expiration
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, payload)

	secret, err := readkey.ReadPrivateKey()
	if err != nil {
		return "", uuid.Nil, err
	}

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", uuid.Nil, err
	}

	return signedToken, tokenId, nil
}

func GenerateRefreshToken(accessTokenId uuid.UUID, userId uuid.UUID, userAgent string, ipAdress string) (string, uuid.UUID, error) {
	tokenId := uuid.New()
	payload := RefreshTokenClaims{
		ATID:      accessTokenId,
		GUID:      userId,
		UserAgent: userAgent, // add hash to protect data
		IpAdress:  ipAdress,  // add hash to protect data
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Set expiration
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, payload)

	secret, err := readkey.ReadPrivateKey()
	if err != nil {
		return "", uuid.Nil, err
	}

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", uuid.Nil, err
	}

	return signedToken, tokenId, nil
}
