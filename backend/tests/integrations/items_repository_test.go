package integrations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	queuedomain "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

func TestItemsRepository_CreateItem(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := postgres.NewItemsRepository(db)
	require.NoError(t, repo.CreateItem(ctx, domain.NewItem(
		"test title",
		"test desc",
		1000,
		true,
	)))
}

func TestItemsRepository_GetAll(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := postgres.NewItemsRepository(db)
	expected := domain.NewItem(
		"test title",
		"test desc",
		1000,
		true,
	)

	require.NoError(t, repo.CreateItem(ctx, expected))

	actual, err := repo.GetAll(ctx)
	require.NoError(t, err)

	assert.NotNil(t, actual)
	assert.Equal(t, expected.Title, actual[0].Title)
	assert.Equal(t, expected.Description, actual[0].Description)
	assert.Equal(t, expected.Price, actual[0].Price)
}

func TestItemsRepository_GetItemByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truncate(db, t)
		ctx := context.Background()
		repo := postgres.NewItemsRepository(db)

		require.NoError(t, repo.CreateItem(ctx, domain.NewItem("a", "b", 1, false)))
		actuals, err := repo.GetAll(ctx)
		require.NoError(t, err)

		id := actuals[0].ID
		actual, err := repo.GetItemByID(ctx, id)
		require.NoError(t, err)

		assert.NotNil(t, actual)
		assert.Equal(t, actual.ID, id)
		assert.Equal(t, actual.Title, actuals[0].Title)
		assert.Equal(t, actual.Description, actuals[0].Description)
		assert.Equal(t, actual.Price, actuals[0].Price)
		assert.Equal(t, 1, actuals[0].Stock)
		assert.Equal(t, false, actuals[0].IsLimited)
	})

	t.Run("error - item not found", func(t *testing.T) {
		truncate(db, t)
		ctx := context.Background()
		repo := postgres.NewItemsRepository(db)

		require.NoError(t, repo.CreateItem(ctx, domain.NewItem("a", "b", 1, false)))

		actual, err := repo.GetItemByID(ctx, uuid.New())
		require.Error(t, err)

		assert.Nil(t, actual)
		assert.ErrorIs(t, err, domain.ErrItemNotFound)
	})
}

func TestItemsRepository_LockStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truncate(db, t)

		itemID := uuid.New()
		seedItem(t, itemID, 5)

		repo := postgres.NewItemsRepository(db)
		item, err := repo.LockStock(context.Background(), itemID)
		require.NoError(t, err)

		assert.Equal(t, itemID, item.ID)
		assert.Equal(t, 5, item.Stock)
		assert.True(t, item.IsLimited)
	})

	t.Run("error - item not found", func(t *testing.T) {
		truncate(db, t)

		repo := postgres.NewItemsRepository(db)
		_, err := repo.LockStock(context.Background(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, queuedomain.ErrItemNotFound)
	})
}
