package integrations

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	postgres2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

func countQueues(ctx context.Context, t *testing.T, id uuid.UUID) int {
	t.Helper()

	var count int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM queues WHERE id = $1`, id).Scan(&count)
	require.NoError(t, err)

	return count
}

func TestTxManager_WithTx(t *testing.T) {
	t.Run("writes made through a repository inside fn are actually committed", func(t *testing.T) {
		truncate(db, t)

		itemsRepo := postgres2.NewItemsRepository(db)
		i := domain2.NewItem("a", "a", 1, true)
		require.NoError(t, itemsRepo.CreateItem(context.Background(), i))
		items, err := itemsRepo.GetAll(context.Background())
		require.NoError(t, err)
		item := items[0]
		ctx := context.Background()

		repo := postgres.NewQueueRepository(db)
		txManager := tx.NewManager(db, slog.Default())
		queue := domain.NewQueue(item.ID, uuid.New())

		err = txManager.WithTx(ctx, func(ctx context.Context) error {
			_, err := repo.CreateQueue(ctx, queue)
			return err
		})
		require.NoError(t, err)

		assert.Equal(t, 1, countQueues(ctx, t, queue.ID))
	})

	t.Run("error inside fn rolls back writes made through a repository", func(t *testing.T) {
		truncate(db, t)

		itemsRepo := postgres2.NewItemsRepository(db)
		i := domain2.NewItem("a", "a", 1, true)
		require.NoError(t, itemsRepo.CreateItem(context.Background(), i))
		items, err := itemsRepo.GetAll(context.Background())
		require.NoError(t, err)
		item := items[0]
		ctx := context.Background()

		repo := postgres.NewQueueRepository(db)
		txManager := tx.NewManager(db, slog.Default())
		queue := domain.NewQueue(item.ID, uuid.New())
		boom := errors.New("boom")

		// Пока repo.pool был жёстко привязан к *pgxpool.Pool, CreateQueue выполнялся
		// как отдельный autocommit-запрос в обход открытой WithTx-транзакции — строка
		// оставалась в БД, даже когда fn возвращала ошибку и WithTx откатывался.
		err = txManager.WithTx(ctx, func(ctx context.Context) error {
			if _, err := repo.CreateQueue(ctx, queue); err != nil {
				return err
			}
			return boom
		})
		require.ErrorIs(t, err, boom)

		assert.Equal(t, 0, countQueues(ctx, t, queue.ID))
	})
}
