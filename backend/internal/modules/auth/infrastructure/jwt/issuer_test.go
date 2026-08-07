package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssuerIssue(t *testing.T) {
	t.Parallel()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	const (
		issuerName = "mock-auth"
		audience   = "xakat0n-api"
		keyID      = "mock-key-1"
		username   = "aleksei"
		ttl        = time.Hour
	)

	userID := uuid.New()
	issuedAt := time.Now()
	issuer := NewIssuer(privateKey, issuerName, audience, keyID, ttl)

	rawToken, err := issuer.Issue(userID, username)
	require.NoError(t, err)
	require.NotEmpty(t, rawToken)

	claims := new(Claims)
	token, err := golangjwt.ParseWithClaims(
		rawToken,
		claims,
		func(token *golangjwt.Token) (any, error) {
			assert.Equal(t, golangjwt.SigningMethodRS256, token.Method)
			assert.Equal(t, keyID, token.Header["kid"])

			return &privateKey.PublicKey, nil
		},
		golangjwt.WithIssuer(issuerName),
		golangjwt.WithAudience(audience),
		golangjwt.WithValidMethods([]string{golangjwt.SigningMethodRS256.Alg()}),
	)
	require.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, userID.String(), claims.Subject)
	require.NotNil(t, claims.IssuedAt)
	require.NotNil(t, claims.ExpiresAt)
	assert.WithinDuration(t, issuedAt, claims.IssuedAt.Time, time.Second)
	assert.WithinDuration(t, issuedAt.Add(ttl), claims.ExpiresAt.Time, time.Second)
}

func TestIssuerIssueRejectsInvalidPrivateKey(t *testing.T) {
	t.Parallel()

	issuer := NewIssuer(nil, "mock-auth", "xakat0n-api", "mock-key-1", time.Hour)

	_, err := issuer.Issue(uuid.New(), "aleksei")
	require.Error(t, err)
	assert.ErrorContains(t, err, "private key is required")
}
