package integrations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	postgres2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
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

func TestQueueRepository_QuitQueue(t *testing.T) {
	const ttl = 2 * time.Second

	t.Run("success - cancel queued user", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		expected := domain.NewQueue(uuid.New(), uuid.New())

		_, err := repo.CreateQueue(ctx, expected)
		require.NoError(t, err)

		actual, err := repo.QuitQueue(ctx, expected.ProductID, expected.UserID, ttl)
		require.NoError(t, err)
		require.NotNil(t, actual)

		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.ProductID, actual.ProductID)
		assert.Equal(t, expected.UserID, actual.UserID)
		assert.Equal(t, domain.QueueStatusCancelled, actual.Status)

		var status domain.QueueStatus
		err = db.QueryRow(ctx, `SELECT status FROM queues WHERE id = $1`, expected.ID).
			Scan(&status)
		require.NoError(t, err)
		assert.Equal(t, domain.QueueStatusCancelled, status)
	})

	t.Run("success - promote next user after offered user quits", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		productID := uuid.New()
		first := domain.NewQueue(productID, uuid.New())
		second := domain.NewQueue(productID, uuid.New())

		_, err := repo.CreateQueue(ctx, first)
		require.NoError(t, err)
		_, err = repo.CreateQueue(ctx, second)
		require.NoError(t, err)

		promoted, _, err := repo.TryPromoteUser(ctx, first.ID, productID, ttl)
		require.NoError(t, err)
		require.True(t, promoted)

		cancelled, err := repo.QuitQueue(ctx, productID, first.UserID, ttl)
		require.NoError(t, err)
		require.NotNil(t, cancelled)
		assert.Equal(t, domain.QueueStatusCancelled, cancelled.Status)

		var (
			status    domain.QueueStatus
			expiresAt *time.Time
		)
		err = db.QueryRow(ctx, `SELECT status, expires_at FROM queues WHERE id = $1`, second.ID).
			Scan(&status, &expiresAt)
		require.NoError(t, err)
		assert.Equal(t, domain.QueueStatusOffered, status)
		require.NotNil(t, expiresAt)
		assert.WithinDuration(t, time.Now().Add(ttl), *expiresAt, time.Second)
	})

	t.Run("error - queue not found", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)

		actual, err := repo.QuitQueue(ctx, uuid.New(), uuid.New(), ttl)

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrQueueNotFound)
		assert.Nil(t, actual)
	})
}
