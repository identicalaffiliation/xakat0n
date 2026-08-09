package logging

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
)

func WithItemID(ctx context.Context, logger ports.Logger, itemID uuid.UUID) context.Context {
	if logger == nil || itemID == uuid.Nil {
		return ctx
	}
	return logger.WithField(ctx, "itemID", itemID)
}

func WithItemTitle(ctx context.Context, logger ports.Logger, title string) context.Context {
	if logger == nil || title == "" {
		return ctx
	}
	return logger.WithField(ctx, "itemTitle", title)
}

func WithItemIDs(ctx context.Context, logger ports.Logger, itemIDs []uuid.UUID) context.Context {
	if logger == nil {
		return ctx
	}
	return logger.WithField(ctx, "itemIDs", itemIDs)
}
