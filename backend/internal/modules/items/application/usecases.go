package application

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
)

type GetAllItemsUsecase struct {
	repo   ports.ItemsRepository
	logger ports.Logger
}

func NewGetAllItemsUsecase(repo ports.ItemsRepository, log ports.Logger) *GetAllItemsUsecase {
	return &GetAllItemsUsecase{
		repo:   repo,
		logger: log,
	}
}

func (u *GetAllItemsUsecase) GetAllItems(ctx context.Context) (*dto.ItemsResponse, error) {
	items, err := u.repo.GetAll(ctx)
	if err != nil {
		u.logger.Error(
			"failed to get all items",
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	return dto.NewItemsResponse(items), nil
}

type GetItemUsecase struct {
	repo   ports.ItemsRepository
	logger ports.Logger
}

func NewGetItemUsecase(repo ports.ItemsRepository, log ports.Logger) *GetItemUsecase {
	return &GetItemUsecase{
		repo:   repo,
		logger: log,
	}
}

func (u *GetItemUsecase) GetItem(ctx context.Context, itemID uuid.UUID) (*dto.ItemResponse, error) {
	item, err := u.repo.GetItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			return nil, err
		}

		u.logger.Error(
			"failed to get item by id",
			"itemId", itemID.String(),
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	return dto.NewItemResponse(item), nil
}
