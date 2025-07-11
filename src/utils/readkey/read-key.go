package readkey

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ReadPrivateKey() (*rsa.PrivateKey, error) {
	// Get private key
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		return nil, errors.New("DANGER: PRIVATE KEY NOT FOUND")
	}
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	// Parse Key
	privateKey, parseErr := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if parseErr != nil {
		return nil, parseErr
	}
	return privateKey, nil
}

func ReadPublicKey() (*rsa.PublicKey, error) {
	// Get public key
	publicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		return nil, errors.New("DANGER: PUBLIC KEY NOT FOUND")
	}
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	// Parse Key
	publicKey, parseErr := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if parseErr != nil {
		return nil, parseErr
	}
	return publicKey, nil
}
