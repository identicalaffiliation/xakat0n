package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/dto"
)

type GetAllItemsUsecase interface {
	GetAllItems(ctx context.Context) ([]dto.Item, error)
}

type GetItemUsecase interface {
	GetItem(ctx context.Context, id uuid.UUID) (*dto.Item, error)
}

type GetSimilarItemsUsecase interface {
	GetSimilarItems(ctx context.Context, itemID uuid.UUID, limit int) ([]dto.Item, error)
}
