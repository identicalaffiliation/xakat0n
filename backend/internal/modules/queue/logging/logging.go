package logging

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
)

func WithQueueID(
	ctx context.Context,
	logger ports.Logger,
	queueID uuid.UUID,
) context.Context {
	return logger.WithField(
		ctx,
		"queueID",
		queueID,
	)
}

func WithItemID(
	ctx context.Context,
	logger ports.Logger,
	itemID uuid.UUID,
) context.Context {
	return logger.WithField(
		ctx,
		"itemID",
		itemID,
	)
}

func WithUserID(
	ctx context.Context,
	logger ports.Logger,
	userID uuid.UUID,
) context.Context {
	return logger.WithField(
		ctx,
		"userID",
		userID,
	)
}
