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

func ReadPublicKey(path string) (any, error) {
	publicKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
