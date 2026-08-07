package logger

import (
	"errors"
	"log/slog"
	"os"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
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

func NewLogger(cfg *config.APIConfig) (*slog.Logger, error) {
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

	h, ok := handlers[cfg.LoggerConfig.Format]
	if !ok {
		return nil, ErrInvalidLoggerFormat
	}

	h = NewHandlerMiddleware(h)

	slogger := slog.New(h)
	return slogger, nil
}
