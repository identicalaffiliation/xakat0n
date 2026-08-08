package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
)

type fakeItemsRepo struct {
	items     []*domain.Item
	item      *domain.Item
	getErr    error
	getAllErr error
}

func (f *fakeItemsRepo) CreateItem(ctx context.Context, item *domain.Item) error {
	return nil
}

func (f *fakeItemsRepo) GetAll(ctx context.Context) ([]*domain.Item, error) {
	return f.items, f.getAllErr
}

func (f *fakeItemsRepo) GetItemByID(ctx context.Context, itemID uuid.UUID) (*domain.Item, error) {
	return f.item, f.getErr
}

type fakeLogger struct{}

func (l *fakeLogger) Debug(msg string, args ...any) {}
func (l *fakeLogger) Error(msg string, args ...any) {}

type fakeSoldOutChecker struct {
	purchased map[uuid.UUID]int
	err       error
}

func (f *fakeSoldOutChecker) CountPurchased(ctx context.Context, itemIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if f.err != nil {
		return nil, f.err
	}

	counts := make(map[uuid.UUID]int, len(itemIDs))
	for _, id := range itemIDs {
		counts[id] = f.purchased[id]
	}

	return counts, nil
}

func TestGetAllItemsUsecase(t *testing.T) {
	t.Parallel()

	t.Run("success returns mapped items with soldOut only for limited ones", func(t *testing.T) {
		limited := uuid.New()
		unlimited := uuid.New()
		items := []*domain.Item{
			{ID: limited, Title: "a", Price: 1, IsLimited: true, Stock: 1},
			{ID: unlimited, Title: "b", Price: 2},
		}
		checker := &fakeSoldOutChecker{purchased: map[uuid.UUID]int{limited: 1}}
		usecase := NewGetAllItemsUsecase(&fakeItemsRepo{items: items}, checker, &fakeLogger{})

		response, err := usecase.GetAllItems(context.Background())
		require.NoError(t, err)
		require.Len(t, response, 2)
		assert.Equal(t, limited, response[0].ItemID)
		assert.True(t, response[0].SoldOut)
		assert.Equal(t, unlimited, response[1].ItemID)
		assert.False(t, response[1].SoldOut)
	})

	t.Run("repository error is wrapped as internal", func(t *testing.T) {
		usecase := NewGetAllItemsUsecase(&fakeItemsRepo{getAllErr: errors.New("db down")}, &fakeSoldOutChecker{}, &fakeLogger{})

		_, err := usecase.GetAllItems(context.Background())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternal)
	})

	t.Run("sold out checker error is wrapped as internal", func(t *testing.T) {
		items := []*domain.Item{{ID: uuid.New(), Title: "a", Price: 1, IsLimited: true, Stock: 1}}
		checker := &fakeSoldOutChecker{err: errors.New("queue db down")}
		usecase := NewGetAllItemsUsecase(&fakeItemsRepo{items: items}, checker, &fakeLogger{})

		_, err := usecase.GetAllItems(context.Background())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternal)
	})

	t.Run("empty catalog returns empty slice, not nil", func(t *testing.T) {
		usecase := NewGetAllItemsUsecase(&fakeItemsRepo{items: nil}, &fakeSoldOutChecker{}, &fakeLogger{})

		response, err := usecase.GetAllItems(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response)
	})

	t.Run("multiple limited items resolve soldOut independently in one batch", func(t *testing.T) {
		soldOutID := uuid.New()
		inStockID := uuid.New()
		items := []*domain.Item{
			{ID: soldOutID, Title: "a", Price: 1, IsLimited: true, Stock: 2},
			{ID: inStockID, Title: "b", Price: 2, IsLimited: true, Stock: 5},
		}
		checker := &fakeSoldOutChecker{purchased: map[uuid.UUID]int{soldOutID: 2, inStockID: 1}}
		usecase := NewGetAllItemsUsecase(&fakeItemsRepo{items: items}, checker, &fakeLogger{})

		response, err := usecase.GetAllItems(context.Background())
		require.NoError(t, err)
		require.Len(t, response, 2)
		assert.True(t, response[0].SoldOut)
		assert.False(t, response[1].SoldOut)
	})
}

func TestGetItemUsecase(t *testing.T) {
	t.Parallel()

	t.Run("success returns mapped item", func(t *testing.T) {
		item := &domain.Item{ID: uuid.New(), Title: "a", Price: 1}
		usecase := NewGetItemUsecase(&fakeItemsRepo{item: item}, &fakeSoldOutChecker{}, &fakeLogger{})

		response, err := usecase.GetItem(context.Background(), item.ID)
		require.NoError(t, err)
		assert.Equal(t, item.ID, response.ItemID)
		assert.Equal(t, "a", response.Title)
	})

	t.Run("limited item sold out when purchased count reaches stock", func(t *testing.T) {
		item := &domain.Item{ID: uuid.New(), Title: "a", Price: 1, IsLimited: true, Stock: 1}
		checker := &fakeSoldOutChecker{purchased: map[uuid.UUID]int{item.ID: 1}}
		usecase := NewGetItemUsecase(&fakeItemsRepo{item: item}, checker, &fakeLogger{})

		response, err := usecase.GetItem(context.Background(), item.ID)
		require.NoError(t, err)
		assert.True(t, response.SoldOut)
	})

	t.Run("limited item not sold out while purchased count below stock", func(t *testing.T) {
		item := &domain.Item{ID: uuid.New(), Title: "a", Price: 1, IsLimited: true, Stock: 3}
		checker := &fakeSoldOutChecker{purchased: map[uuid.UUID]int{item.ID: 1}}
		usecase := NewGetItemUsecase(&fakeItemsRepo{item: item}, checker, &fakeLogger{})

		response, err := usecase.GetItem(context.Background(), item.ID)
		require.NoError(t, err)
		assert.False(t, response.SoldOut)
		require.NotNil(t, response.Stock)
		assert.Equal(t, 3, *response.Stock)
	})

	t.Run("not found is passed through", func(t *testing.T) {
		usecase := NewGetItemUsecase(&fakeItemsRepo{getErr: domain.ErrItemNotFound}, &fakeSoldOutChecker{}, &fakeLogger{})

		_, err := usecase.GetItem(context.Background(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrItemNotFound)
	})

	t.Run("repository error is wrapped as internal", func(t *testing.T) {
		usecase := NewGetItemUsecase(&fakeItemsRepo{getErr: errors.New("db down")}, &fakeSoldOutChecker{}, &fakeLogger{})

		_, err := usecase.GetItem(context.Background(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternal)
	})
}
