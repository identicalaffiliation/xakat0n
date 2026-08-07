package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type QueueRepository interface {
	CreateQueue(ctx context.Context, queue *domain.Queue) (*domain.Queue, error)
	TryPromoteUser(ctx context.Context, queueID, productID uuid.UUID, ttl time.Duration) (bool, *time.Time, error)

	ExpireStale(ctx context.Context, productID uuid.UUID) error
	CountTaken(ctx context.Context, productID uuid.UUID) (int, error)
	PromoteNext(ctx context.Context, productID uuid.UUID, freeSlots int, ttl time.Duration) (int64, error)
	MarkSoldOut(ctx context.Context, productID uuid.UUID, stock int) error

	IsSoldOut(ctx context.Context, productID uuid.UUID, stock int) (bool, error)
	GetActiveTicket(ctx context.Context, productID, userID uuid.UUID) (*domain.Queue, error)
	GetLatestTicket(ctx context.Context, productID, userID uuid.UUID) (*domain.Queue, error)
	CountQueuedAhead(ctx context.Context, productID uuid.UUID, createdAt time.Time) (int, error)
	NextSlotFreeAt(ctx context.Context, productID uuid.UUID) (*time.Time, error)
}

// ItemsRepository — consumer-defined интерфейс: queue-модуль не читает items
// напрямую, а зависит только от того, что реально нужно advanceQueue. Реализация
// живёт в internal/modules/items.
type ItemsRepository interface {
	LockStock(ctx context.Context, itemID uuid.UUID) (*domain.Item, error)
}
