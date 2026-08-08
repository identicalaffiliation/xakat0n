package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
)

// QueueRepository — реализация живёт в internal/modules/queue.
type QueueRepository interface {
	TryStartCheckout(ctx context.Context, itemID, userID uuid.UUID) (*domain.Queue, bool, error)
	FinalizeCheckoutResult(ctx context.Context, itemID, userID, ticketID uuid.UUID, paid bool) (*domain.Queue, bool, error)
	FindTicket(ctx context.Context, itemID, userID, ticketID uuid.UUID) (*domain.Queue, error)
}

// ItemsRepository — реализация живёт в internal/modules/items.
type ItemsRepository interface {
	IsLimited(ctx context.Context, itemID uuid.UUID) (bool, error)
}
