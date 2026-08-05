package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// The path is supplied by trusted application configuration, not by an HTTP request.
	data, err := os.ReadFile(path) //nolint:gosec // G304: configured private-key path is intentional.
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	privateKey, err := golangjwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	return privateKey, nil
}
