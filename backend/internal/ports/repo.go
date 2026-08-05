package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type QueueRepository interface {
	CreateQueue(ctx context.Context, queue *domain.Queue) (*domain.Queue, error)
	TryPromoteUser(ctx context.Context, queueID, productID uuid.UUID, ttl time.Duration) (bool, *time.Time, error)
	QuitQueue(ctx context.Context, productID, userID uuid.UUID, ttl time.Duration) (*domain.Queue, error)
}
