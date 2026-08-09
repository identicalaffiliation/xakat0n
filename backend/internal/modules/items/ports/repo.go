package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type ItemsRepository interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetAll(ctx context.Context) ([]*domain.Item, error)
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*domain.Item, error)
	GetSimilarByCategory(ctx context.Context, itemID uuid.UUID, category string, limit int) ([]*domain.Item, error)
}

// SoldOutChecker — реализуется queue-модулем (queues — его таблица), items
// сам её не запрашивает и только сравнивает результат со Stock.
type SoldOutChecker interface {
	CountPurchased(ctx context.Context, itemIDs []uuid.UUID) (map[uuid.UUID]int, error)
}
