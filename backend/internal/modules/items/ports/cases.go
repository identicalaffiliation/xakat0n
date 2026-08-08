package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/dto"
)

type GetAllItemsUsecase interface {
	GetAllItems(ctx context.Context) (*dto.ItemsResponse, error)
}

type GetItemUsecase interface {
	GetItem(ctx context.Context, id uuid.UUID) (*dto.ItemResponse, error)
}
