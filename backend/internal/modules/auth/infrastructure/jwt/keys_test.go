package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadPrivateKey(t *testing.T) {
	t.Parallel()

	expected, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	der, err := x509.MarshalPKCS8PrivateKey(expected)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "private.pem")
	data := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	require.NoError(t, os.WriteFile(path, data, 0o600))

	actual, err := LoadPrivateKey(path)
	require.NoError(t, err)
	assert.Equal(t, expected.N, actual.N)
	assert.Equal(t, expected.E, actual.E)
}

func TestLoadPrivateKeyErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		_, err := LoadPrivateKey(filepath.Join(t.TempDir(), "missing.pem"))
		require.Error(t, err)
		assert.ErrorContains(t, err, "read private key")
	})

	t.Run("invalid PEM", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "private.pem")
		require.NoError(t, os.WriteFile(path, []byte("not a private key"), 0o600))

		_, err := LoadPrivateKey(path)
		require.Error(t, err)
		assert.ErrorContains(t, err, "parse private key")
	})
}
