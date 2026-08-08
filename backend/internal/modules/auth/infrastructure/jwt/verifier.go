package jwt

import (
	"crypto/rsa"
	"fmt"

	golangjwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Verifier struct {
	publicKey *rsa.PublicKey
	issuer    string
	audience  string
	keyID     string
}

func NewVerifier(
	publicKey *rsa.PublicKey,
	issuer string,
	audience string,
	keyID string,
) *Verifier {
	return &Verifier{
		publicKey: publicKey,
		issuer:    issuer,
		audience:  audience,
		keyID:     keyID,
	}
}

func (v *Verifier) Verify(tokenString string) (uuid.UUID, error) {
	if v.publicKey == nil {
		return uuid.Nil, ErrPrivateKeyRequired
	}

	claims := new(Claims)

	token, err := golangjwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *golangjwt.Token) (any, error) {
			keyID, ok := t.Header["kid"].(string)
			if !ok || keyID != v.keyID {
				return nil, ErrInvalidKeyID
			}

			return v.publicKey, nil
		},
		golangjwt.WithValidMethods([]string{
			golangjwt.SigningMethodRS256.Alg(),
		}),
		golangjwt.WithIssuer(v.issuer),
		golangjwt.WithAudience(v.audience),
		golangjwt.WithExpirationRequired(),
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidSubject, err)
	}
	return userID, nil
}
