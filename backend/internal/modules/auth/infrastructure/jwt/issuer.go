package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewClaims(
	username string,
	issuer string,
	subject string,
	audience string,
	issuedAt time.Time,
	expiresAt time.Time,
) *Claims {
	return &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings{audience},
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
}

type Issuer struct {
	privateKey *rsa.PrivateKey
	issuer     string
	audience   string
	keyID      string
	ttl        time.Duration
}

func NewIssuer(
	privateKey *rsa.PrivateKey,
	issuer string,
	audience string,
	keyID string,
	ttl time.Duration,
) *Issuer {
	return &Issuer{
		privateKey: privateKey,
		issuer:     issuer,
		audience:   audience,
		keyID:      keyID,
		ttl:        ttl,
	}
}

func (i *Issuer) Issue(userID uuid.UUID, username string) (string, error) {
	now := time.Now()

	claims := NewClaims(
		username,
		i.issuer,
		userID.String(),
		i.audience,
		now,
		now.Add(i.ttl),
	)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = i.keyID

	return token.SignedString(i.privateKey)
}
