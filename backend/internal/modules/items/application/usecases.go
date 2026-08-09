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
	repo    ports.ItemsRepository
	soldOut ports.SoldOutChecker
	logger  ports.Logger
}

func NewGetAllItemsUsecase(repo ports.ItemsRepository, soldOut ports.SoldOutChecker, log ports.Logger) *GetAllItemsUsecase {
	return &GetAllItemsUsecase{
		repo:    repo,
		soldOut: soldOut,
		logger:  log,
	}
}

func (u *GetAllItemsUsecase) GetAllItems(ctx context.Context) ([]dto.Item, error) {
	items, err := u.repo.GetAll(ctx)
	if err != nil {
		u.logger.Error(
			"failed to get all items",
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	soldOut, err := soldOutByItemID(ctx, u.soldOut, items)
	if err != nil {
		u.logger.Error(
			"failed to count purchased items",
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	return dto.NewItems(items, soldOut), nil
}

type GetItemUsecase struct {
	repo    ports.ItemsRepository
	soldOut ports.SoldOutChecker
	logger  ports.Logger
}

func NewGetItemUsecase(repo ports.ItemsRepository, soldOut ports.SoldOutChecker, log ports.Logger) *GetItemUsecase {
	return &GetItemUsecase{
		repo:    repo,
		soldOut: soldOut,
		logger:  log,
	}
}

func (u *GetItemUsecase) GetItem(ctx context.Context, itemID uuid.UUID) (*dto.Item, error) {
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

	soldOut, err := soldOutByItemID(ctx, u.soldOut, []*domain.Item{item})
	if err != nil {
		u.logger.Error(
			"failed to count purchased items",
			"itemId", itemID.String(),
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	result := dto.NewItem(item, soldOut[item.ID])
	return &result, nil
}

type GetSimilarItemsUsecase struct {
	repo    ports.ItemsRepository
	soldOut ports.SoldOutChecker
	logger  ports.Logger
}

func NewGetSimilarItemsUsecase(repo ports.ItemsRepository, soldOut ports.SoldOutChecker, log ports.Logger) *GetSimilarItemsUsecase {
	return &GetSimilarItemsUsecase{
		repo:    repo,
		soldOut: soldOut,
		logger:  log,
	}
}

func (u *GetSimilarItemsUsecase) GetSimilarItems(ctx context.Context, itemID uuid.UUID, limit int) ([]dto.Item, error) {
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

	// Не с чем сравнивать — товар без category похожих не имеет.
	if item.Category == nil {
		return []dto.Item{}, nil
	}

	similar, err := u.repo.GetSimilarByCategory(ctx, itemID, *item.Category, limit)
	if err != nil {
		u.logger.Error(
			"failed to get similar items",
			"itemId", itemID.String(),
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	soldOut, err := soldOutByItemID(ctx, u.soldOut, similar)
	if err != nil {
		u.logger.Error(
			"failed to count purchased items",
			"itemId", itemID.String(),
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	return dto.NewItems(similar, soldOut), nil
}

// soldOutByItemID считает soldOut только для лимитированных товаров: у
// нелимитированных стока в контрактном смысле нет, соответственно false
// без обращения к queue-модулю.
func soldOutByItemID(ctx context.Context, checker ports.SoldOutChecker, items []*domain.Item) (map[uuid.UUID]bool, error) {
	limitedIDs := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		if item.IsLimited {
			limitedIDs = append(limitedIDs, item.ID)
		}
	}

	purchased, err := checker.CountPurchased(ctx, limitedIDs)
	if err != nil {
		return nil, err
	}

	soldOut := make(map[uuid.UUID]bool, len(limitedIDs))
	for _, item := range items {
		if item.IsLimited {
			soldOut[item.ID] = purchased[item.ID] >= item.Stock
		}
	}

	return soldOut, nil
}
