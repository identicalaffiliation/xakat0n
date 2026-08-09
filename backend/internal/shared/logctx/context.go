package logctx

import (
	"context"
	"errors"
	"maps"
	"reflect"
)

type contextKey struct{}

var logKey contextKey

type Fields map[string]any

type LogCtx struct{}

func NewLogCtx() *LogCtx {
	return &LogCtx{}
}

func (l *LogCtx) WithField(ctx context.Context, key string, value any) context.Context {
	current, _ := ctx.Value(logKey).(Fields)
	if existing, ok := current[key]; ok && reflect.DeepEqual(existing, value) {
		return ctx
	}

	fields := make(Fields, len(current)+1)
	maps.Copy(fields, current)

	fields[key] = value
	return context.WithValue(ctx, logKey, fields)
}

func (l *LogCtx) FieldsFromContext(ctx context.Context) map[string]any {
	current, _ := ctx.Value(logKey).(Fields)
	return current
}

func (l *LogCtx) WrapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	fields := l.FieldsFromContext(ctx)
	copied := make(Fields, len(fields))

	maps.Copy(copied, fields)

	return &logContextError{
		err:    err,
		fields: copied,
	}
}

func (l *LogCtx) ContextFromError(ctx context.Context, err error) context.Context {
	var wrapped *logContextError
	if !errors.As(err, &wrapped) {
		return ctx
	}

	for k, v := range wrapped.fields {
		ctx = l.WithField(ctx, k, v)
	}

	return ctx
}
