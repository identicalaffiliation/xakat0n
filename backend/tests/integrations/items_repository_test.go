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
	category := "test category"
	expected.Category = &category
	expected.Stock = 7
	imagePath := "test.jpg"
	expected.ImagePath = &imagePath

	require.NoError(t, repo.CreateItem(ctx, expected))

	actual, err := repo.GetAll(ctx)
	require.NoError(t, err)

	assert.NotNil(t, actual)
	assert.Equal(t, expected.Title, actual[0].Title)
	assert.Equal(t, expected.Description, actual[0].Description)
	assert.Equal(t, expected.Price, actual[0].Price)
	assert.Equal(t, expected.Category, actual[0].Category)
	assert.Equal(t, expected.Stock, actual[0].Stock)
	assert.Equal(t, expected.ImagePath, actual[0].ImagePath)
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

func TestItemsRepository_GetSimilarByCategory(t *testing.T) {
	t.Run("returns items in the same category, excluding itself", func(t *testing.T) {
		truncate(db, t)
		ctx := context.Background()
		repo := postgres.NewItemsRepository(db)

		category := "Недвижимость на Луне"
		other := "Хозтовары"

		self := domain.NewItem("Море Спокойствия", "d", 100, true)
		self.Category = &category
		require.NoError(t, repo.CreateItem(ctx, self))

		sameCategory := domain.NewItem("Море Дождей", "d", 100, true)
		sameCategory.Category = &category
		require.NoError(t, repo.CreateItem(ctx, sameCategory))

		otherCategory := domain.NewItem("Пакет пакетов", "d", 100, true)
		otherCategory.Category = &other
		require.NoError(t, repo.CreateItem(ctx, otherCategory))

		all, err := repo.GetAll(ctx)
		require.NoError(t, err)
		selfID := findByTitle(t, all, "Море Спокойствия").ID

		actual, err := repo.GetSimilarByCategory(ctx, selfID, category, 6)
		require.NoError(t, err)

		require.Len(t, actual, 1)
		assert.Equal(t, "Море Дождей", actual[0].Title)
	})

	t.Run("respects limit", func(t *testing.T) {
		truncate(db, t)
		ctx := context.Background()
		repo := postgres.NewItemsRepository(db)

		category := "Недвижимость на Луне"
		for _, title := range []string{"a", "b", "c", "d"} {
			item := domain.NewItem(title, "d", 100, true)
			item.Category = &category
			require.NoError(t, repo.CreateItem(ctx, item))
		}

		all, err := repo.GetAll(ctx)
		require.NoError(t, err)

		actual, err := repo.GetSimilarByCategory(ctx, all[0].ID, category, 2)
		require.NoError(t, err)
		assert.Len(t, actual, 2)
	})

	t.Run("no matches returns empty slice", func(t *testing.T) {
		truncate(db, t)
		ctx := context.Background()
		repo := postgres.NewItemsRepository(db)

		category := "Недвижимость на Луне"
		item := domain.NewItem("a", "d", 100, true)
		item.Category = &category
		require.NoError(t, repo.CreateItem(ctx, item))

		all, err := repo.GetAll(ctx)
		require.NoError(t, err)

		actual, err := repo.GetSimilarByCategory(ctx, all[0].ID, category, 6)
		require.NoError(t, err)
		assert.Empty(t, actual)
	})
}

func findByTitle(t *testing.T, items []*domain.Item, title string) *domain.Item {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return item
		}
	}

	t.Fatalf("item with title %q not found", title)
	return nil
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
