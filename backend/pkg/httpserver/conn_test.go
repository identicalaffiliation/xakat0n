package httpserver

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/config"
)

func TestCreateConn(t *testing.T) {
	t.Parallel()
	require.NotNil(t, SetupServer(&config.ServerConfig{}, nil, nil))
}
