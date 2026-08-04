package httpserver

import (
	"testing"

	"github.com/identicalaffiliation/xakat0n/backend/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCreateConn(t *testing.T) {
	t.Parallel()
	require.NotNil(t, SetupServer(&config.ServerConfig{}, nil))
}
