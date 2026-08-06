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
}
