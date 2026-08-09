package logger

import "context"

func (l *Logging) WithField(ctx context.Context, key string, value any) context.Context {
	return l.context.WithField(ctx, key, value)
}

func (l *Logging) WrapError(ctx context.Context, err error) error {
	return l.context.WrapError(ctx, err)
}

func (l *Logging) ContextFromError(ctx context.Context, err error) context.Context {
	return l.context.ContextFromError(ctx, err)
}
