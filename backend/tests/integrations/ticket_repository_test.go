package integrations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domain2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	postgres2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
)

func TestQueueRepository_CreateQueue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truncate(db, t)

		itemsRepo := postgres2.NewItemsRepository(db)
		i := domain2.NewItem("a", "a", 1)
		require.NoError(t, itemsRepo.CreateItem(context.Background(), i))
		items, err := itemsRepo.GetAll(context.Background())
		require.NoError(t, err)
		item := items[0]
		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		expected := domain.NewQueue(item.ID, uuid.New())

		actual, err := repo.CreateQueue(ctx, expected)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.UserID, actual.UserID)
		assert.Equal(t, expected.ProductID, actual.ProductID)
		assert.Equal(t, domain.QueueStatusQueued, actual.Status)
		assert.Nil(t, actual.ExpiresAt)
		assert.WithinDuration(t, time.Now().UTC(), actual.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), actual.UpdatedAt, time.Second)
	})

	t.Run("error - user already queued", func(t *testing.T) {
		truncate(db, t)

		itemsRepo := postgres2.NewItemsRepository(db)
		i := domain2.NewItem("a", "a", 1)
		require.NoError(t, itemsRepo.CreateItem(context.Background(), i))
		items, err := itemsRepo.GetAll(context.Background())
		require.NoError(t, err)
		item := items[0]
		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		expected := domain.NewQueue(item.ID, uuid.New())

		_, err = repo.CreateQueue(ctx, expected)
		require.NoError(t, err)

		_, err = repo.CreateQueue(ctx, expected)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyQueued)
	})
}

func TestQueueRepository_TryPromoteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		truncate(db, t)

		itemsRepo := postgres2.NewItemsRepository(db)
		i := domain2.NewItem("a", "a", 1)
		require.NoError(t, itemsRepo.CreateItem(context.Background(), i))
		items, err := itemsRepo.GetAll(context.Background())
		require.NoError(t, err)
		item := items[0]
		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		expected := domain.NewQueue(item.ID, uuid.New())

		_, err = repo.CreateQueue(ctx, expected)
		require.NoError(t, err)

		promoted, ttl, err := repo.TryPromoteUser(
			ctx,
			expected.ID,
			expected.ProductID,
			time.Second*2,
		)
		require.NoError(t, err)
		assert.NotNil(t, ttl)
		expectedExpiresAt := time.Now().Add(2 * time.Second)
		assert.WithinDuration(
			t,
			expectedExpiresAt,
			*ttl,
			1*time.Second,
		)
		assert.Equal(t, true, promoted)
	})

	t.Run("catch race condition", func(t *testing.T) {
		truncate(db, t)

		itemsRepo := postgres2.NewItemsRepository(db)
		i := domain2.NewItem("a", "a", 1)
		require.NoError(t, itemsRepo.CreateItem(context.Background(), i))
		items, err := itemsRepo.GetAll(context.Background())
		require.NoError(t, err)
		item := items[0]
		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		expected := domain.NewQueue(item.ID, uuid.New())
		_, err = repo.CreateQueue(ctx, expected)
		require.NoError(t, err)

		_, _, err = repo.TryPromoteUser(
			ctx,
			expected.ID,
			expected.ProductID,
			time.Second*2,
		)
		require.NoError(t, err)
		promoted, _, err := repo.TryPromoteUser(
			ctx,
			expected.ID,
			expected.ProductID,
			time.Second*2,
		)
		require.NoError(t, err)
		assert.Equal(t, false, promoted)
	})
}
