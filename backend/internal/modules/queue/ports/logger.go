package ports

import "context"

type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)

	WithField(
		ctx context.Context,
		key string,
		value any,
	) context.Context

	WrapError(
		ctx context.Context,
		err error,
	) error

	ContextFromError(
		ctx context.Context,
		err error,
	) context.Context
}
