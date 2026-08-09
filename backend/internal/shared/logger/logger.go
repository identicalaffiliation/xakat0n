package logger

import (
	"errors"
	"log/slog"
	"os"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/logctx"
)

const (
	LevelDebug = "debug"
	LevelError = "error"
	JsonFormat = "json"
	TextFormat = "text"
)

var (
	ErrInvalidLoggerLevel  = errors.New("invalid logger level")
	ErrInvalidLoggerFormat = errors.New("invalid logger format")

	stdout = os.Stdout
)

type Logging struct {
	*slog.Logger
	context *logctx.LogCtx
}

func NewLogger(cfg *config.APIConfig) (*Logging, error) {
	levels := map[string]slog.Level{
		LevelDebug: slog.LevelDebug,
		LevelError: slog.LevelError,
	}

	level, ok := levels[cfg.LoggerConfig.Level]
	if !ok {
		return nil, ErrInvalidLoggerLevel
	}

	handlers := map[string]slog.Handler{
		JsonFormat: slog.NewJSONHandler(stdout, &slog.HandlerOptions{
			Level: level,
		}),

		TextFormat: slog.NewTextHandler(stdout, &slog.HandlerOptions{
			Level: level,
		}),
	}

	handler, ok := handlers[cfg.LoggerConfig.Format]
	if !ok {
		return nil, ErrInvalidLoggerFormat
	}
	logContext := logctx.NewLogCtx()
	handler = NewHandlerMiddleware(handler, logContext)

	return &Logging{
		Logger:  slog.New(handler),
		context: logContext,
	}, nil
}
