package logging

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/ports"
)

func WithItemID(ctx context.Context, logger ports.Logger, itemID uuid.UUID) context.Context {
	if logger == nil {
		return ctx
	}
	return logger.WithField(ctx, "itemID", itemID)
}

func WithUserID(ctx context.Context, logger ports.Logger, userID uuid.UUID) context.Context {
	if logger == nil {
		return ctx
	}
	return logger.WithField(ctx, "userID", userID)
}

func WithTicketID(ctx context.Context, logger ports.Logger, ticketID uuid.UUID) context.Context {
	if logger == nil {
		return ctx
	}
	return logger.WithField(ctx, "ticketID", ticketID)
}
