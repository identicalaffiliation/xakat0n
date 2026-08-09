package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/logctx"
)

type customLoggerStub struct {
	level string
	msg   string
}

func (l *customLoggerStub) DebugContext(_ context.Context, msg string, _ ...any) {
	l.level, l.msg = "debug", msg
}
func (l *customLoggerStub) InfoContext(_ context.Context, msg string, _ ...any) {
	l.level, l.msg = "info", msg
}
func (l *customLoggerStub) WarnContext(_ context.Context, msg string, _ ...any) {
	l.level, l.msg = "warn", msg
}
func (l *customLoggerStub) ErrorContext(_ context.Context, msg string, _ ...any) {
	l.level, l.msg = "error", msg
}

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

func TestLoggingAcceptsCustomLogger(t *testing.T) {
	base := &customLoggerStub{}
	logging := NewLogging(base, logctx.NewLogCtx())
	ctx := logging.WithField(context.Background(), "itemID", "item-1")

	logging.InfoContext(ctx, "custom logger works")

	assert.Equal(t, "info", base.level)
	assert.Equal(t, "custom logger works", base.msg)
}
