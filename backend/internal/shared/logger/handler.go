package logger

import (
	"context"
	"log/slog"
)

type ContextFields interface {
	FieldsFromContext(ctx context.Context) map[string]any
}

type HandlerMiddleware struct {
	next   slog.Handler
	fields ContextFields
}

func NewHandlerMiddleware(next slog.Handler, fields ContextFields) *HandlerMiddleware {
	return &HandlerMiddleware{
		next:   next,
		fields: fields,
	}
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	for k, v := range h.fields.FieldsFromContext(ctx) {
		rec.Add(k, v)
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{
		next:   h.next.WithAttrs(attrs),
		fields: h.fields,
	}
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{
		next:   h.next.WithGroup(name),
		fields: h.fields,
	}
}
