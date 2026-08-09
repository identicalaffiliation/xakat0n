package logging

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/ports"
)

func WithUserID(ctx context.Context, logger ports.Logger, userID uuid.UUID) context.Context {
	if logger == nil {
		return ctx
	}
	return logger.WithField(ctx, "userID", userID)
}

func WithUsername(ctx context.Context, logger ports.Logger, username string) context.Context {
	if logger == nil {
		return ctx
	}
	return logger.WithField(ctx, "username", username)
}
