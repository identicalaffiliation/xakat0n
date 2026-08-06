package integrations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

func seedItem(t *testing.T, itemID uuid.UUID, stock int) {
	t.Helper()

	_, err := db.Exec(
		context.Background(),
		`INSERT INTO items (id, title, price, is_limited, stock) VALUES ($1, 'test item', 100, true, $2)`,
		itemID,
		stock,
	)
	require.NoError(t, err)
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
		assert.ErrorIs(t, err, domain.ErrItemNotFound)
	})
}
