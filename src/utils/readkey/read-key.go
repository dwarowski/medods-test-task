package readkey

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ReadPrivateKey() (*rsa.PrivateKey, error) {
	// Get private key
	privateKeyBytes := os.Getenv("PRIVATE_KEY_PATH")
	if privateKeyBytes == "" {
		return nil, errors.New("DANGER: PRIVATE KEY NOT FOUND")
	}

	// Parse Key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyBytes))
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func ReadPublicKey() (*rsa.PublicKey, error) {
	// Get public key
	publicKeyBytes := os.Getenv("PUBLIC_KEY_PATH")
	if publicKeyBytes == "" {
		return nil, errors.New("DANGER: PUBLIC KEY NOT FOUND")
	}

	// Parse Key
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyBytes))
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
