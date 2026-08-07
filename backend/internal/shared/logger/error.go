package logger

import (
	"maps"
	"context"
	"errors"
)

type errorWithLogContext struct {
	err    error
	fields Fields
}

func (e *errorWithLogContext) Error() string {
	return e.err.Error()
}

func (e *errorWithLogContext) Unwrap() error {
	return e.err
}

func WrapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	fields := FieldsFromContext(ctx)
	copied := make(Fields, len(fields))

	maps.Copy(copied, fields)

	return &errorWithLogContext{
		err:    err,
		fields: copied,
	}
}

func ContextFromError(ctx context.Context, err error) context.Context {
	var wrapped *errorWithLogContext
	if !errors.As(err, &wrapped) {
		return ctx
	}

	for k, v := range wrapped.fields {
		ctx = WithField(ctx, k, v)
	}

	return ctx
}
