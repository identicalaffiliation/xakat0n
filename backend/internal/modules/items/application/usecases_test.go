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

func TestGetAllItemsUsecase(t *testing.T) {
	t.Parallel()

	t.Run("success returns mapped items", func(t *testing.T) {
		items := []*domain.Item{
			{ID: uuid.New(), Title: "a", Price: 1},
			{ID: uuid.New(), Title: "b", Price: 2},
		}
		usecase := NewGetAllItemsUsecase(&fakeItemsRepo{items: items}, &fakeLogger{})

		response, err := usecase.GetAllItems(context.Background())
		require.NoError(t, err)
		require.Len(t, response.Items, 2)
		assert.Equal(t, items[0].ID, response.Items[0].Item.ID)
		assert.Equal(t, items[1].ID, response.Items[1].Item.ID)
	})

	t.Run("repository error is wrapped as internal", func(t *testing.T) {
		usecase := NewGetAllItemsUsecase(&fakeItemsRepo{getAllErr: errors.New("db down")}, &fakeLogger{})

		_, err := usecase.GetAllItems(context.Background())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternal)
	})
}

func TestGetItemUsecase(t *testing.T) {
	t.Parallel()

	t.Run("success returns mapped item", func(t *testing.T) {
		item := &domain.Item{ID: uuid.New(), Title: "a", Price: 1}
		usecase := NewGetItemUsecase(&fakeItemsRepo{item: item}, &fakeLogger{})

		response, err := usecase.GetItem(context.Background(), item.ID)
		require.NoError(t, err)
		assert.Equal(t, item.ID, response.Item.ID)
		assert.Equal(t, "a", response.Item.Title)
	})

	t.Run("not found is passed through", func(t *testing.T) {
		usecase := NewGetItemUsecase(&fakeItemsRepo{getErr: domain.ErrItemNotFound}, &fakeLogger{})

		_, err := usecase.GetItem(context.Background(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrItemNotFound)
	})

	t.Run("repository error is wrapped as internal", func(t *testing.T) {
		usecase := NewGetItemUsecase(&fakeItemsRepo{getErr: errors.New("db down")}, &fakeLogger{})

		_, err := usecase.GetItem(context.Background(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternal)
	})
}
