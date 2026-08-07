package logger

import "maps"

import "context"

type contextKey struct{}

var logKey contextKey

type Fields map[string]any

func WithField(ctx context.Context, key string, value any) context.Context {
	current, _ := ctx.Value(logKey).(Fields)

	fields := make(Fields, len(current)+1)
	maps.Copy(fields, current)

	fields[key] = value
	return context.WithValue(ctx, logKey, fields)
}

func FieldsFromContext(ctx context.Context) Fields {
	current, _ := ctx.Value(logKey).(Fields)
	return current
}
