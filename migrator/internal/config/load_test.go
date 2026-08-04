package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tempDir := os.TempDir()
	path := filepath.Join(tempDir, "config.yml")

	body := `
logger:
  format: text
  level: debug
`

	require.NoError(t, os.WriteFile(path, []byte(body), 0o644))

	res, err := LoadConfig(path)
	require.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.LoggerConfig.Format, "text")
	assert.Equal(t, res.LoggerConfig.Level, "debug")
}
