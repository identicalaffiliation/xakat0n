package httpserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
)

func TestNew(t *testing.T) {
	t.Parallel()
	require.NotNil(t, New(&config.ServerConfig{}, http.NewServeMux()))
}

func TestNewRouter(t *testing.T) {
	t.Parallel()
	require.NotNil(t, NewRouter())
}
