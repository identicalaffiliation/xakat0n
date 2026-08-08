package jwt

import "errors"

var (
	ErrPrivateKeyRequired = errors.New("jwt private key is required")
	ErrPublicKeyRequired  = errors.New("jwt public key is required")
	ErrInvalidKeyID       = errors.New("invalid keyID")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidSubject     = errors.New("invalid subject")
)
