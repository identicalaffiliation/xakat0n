package logctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFieldDoesNotReplaceContextForSameValue(t *testing.T) {
	logContext := NewLogCtx()
	ctx := logContext.WithField(context.Background(), "userID", "user-1")

	same := logContext.WithField(ctx, "userID", "user-1")

	assert.Equal(t, ctx, same)
	assert.Equal(t, map[string]any{"userID": "user-1"}, logContext.FieldsFromContext(same))
}

func TestWithFieldReplacesChangedValue(t *testing.T) {
	logContext := NewLogCtx()
	ctx := logContext.WithField(context.Background(), "queueID", "queue-1")

	updated := logContext.WithField(ctx, "queueID", "queue-2")

	assert.NotEqual(t, ctx, updated)
	assert.Equal(t, "queue-2", logContext.FieldsFromContext(updated)["queueID"])
}
