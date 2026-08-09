package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
)

func TestCreateLogger(t *testing.T) {
	t.Parallel()

	logger, err := NewLogger(
		&config.APIConfig{
			LoggerConfig: config.LoggerConfig{
				Level:  LevelDebug,
				Format: JsonFormat,
			},
		},
	)
	require.NoError(t, err)
	assert.NotNil(t, logger)
}
