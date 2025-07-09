package readkey

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

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
