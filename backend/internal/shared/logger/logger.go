package logger

import (
	"context"
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

type ContextLogger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type Logging struct {
	logger  ContextLogger
	context *logctx.LogCtx
}

func NewLogging(logger ContextLogger, logContext *logctx.LogCtx) *Logging {
	return &Logging{
		logger:  logger,
		context: logContext,
	}
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

	return NewLogging(slog.New(handler), logContext), nil
}

func (l *Logging) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *Logging) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *Logging) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *Logging) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
