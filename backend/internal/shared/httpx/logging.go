package httpx

import (
	"context"

	"github.com/google/uuid"
)

type LogContext interface {
	WithField(ctx context.Context, key string, value any) context.Context
}

func withLogUserID(
	ctx context.Context,
	logctx LogContext,
	userID uuid.UUID,
) context.Context {
	return logctx.WithField(ctx, "userID", userID)
}
