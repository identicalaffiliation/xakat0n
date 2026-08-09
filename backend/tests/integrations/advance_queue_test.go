package integrations

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	itemspostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

func newAdvanceQueueUsecase() *application.AdvanceQueueUsecase {
	queueRepo := postgres.NewQueueRepository(db)
	itemsRepo := itemspostgres.NewItemsRepository(db)
	txManager := tx.NewManager(db, slog.Default())

	return application.NewAdvanceQueueUsecase(itemsRepo, queueRepo, txManager, slog.Default())
}

func countTicketsByStatus(t *testing.T, itemID uuid.UUID, status domain.QueueStatus) int {
	t.Helper()

	var count int
	err := db.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM queues WHERE item_id = $1 AND status = $2::queue_status`,
		itemID,
		status,
	).Scan(&count)
	require.NoError(t, err)

	return count
}

// seedQueuedBatch сеет n QUEUED-тикетов на один товар с гарантированно возрастающим
// created_at (репозиторий CreateQueue не даёт задать created_at явно, а конкурентный
// INSERT ... DEFAULT now() внутри одной операции не гарантирует детерминированный порядок).
func seedQueuedBatch(t *testing.T, itemID uuid.UUID, n int) []uuid.UUID {
	t.Helper()

	base := time.Now().UTC().Add(-time.Hour)
	ids := make([]uuid.UUID, n)
	for i := range ids {
		ids[i] = uuid.New()
		seedTicket(t, ids[i], itemID, uuid.New(), domain.QueueStatusQueued, base.Add(time.Duration(i)*time.Millisecond), nil)
	}

	return ids
}

func TestAdvanceQueueUsecase_NFR1_Atomicity(t *testing.T) {
	t.Run("stock=1: exactly 1 OFFERED out of 100 concurrent calls", func(t *testing.T) {
		truncate(db, t)

		itemID := uuid.New()
		seedItem(t, itemID, 1)
		seedQueuedBatch(t, itemID, 100)

		usecase := newAdvanceQueueUsecase()

		var wg sync.WaitGroup
		errs := make(chan error, 100)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := usecase.AdvanceQueue(context.Background(), itemID, 3*time.Second)
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)

		for err := range errs {
			require.NoError(t, err)
		}

		assert.Equal(t, 1, countTicketsByStatus(t, itemID, domain.QueueStatusOffered))
		assert.Equal(t, 99, countTicketsByStatus(t, itemID, domain.QueueStatusQueued))
	})

	t.Run("stock=5: exactly 5 OFFERED out of 100 concurrent calls, earliest by created_at", func(t *testing.T) {
		truncate(db, t)

		itemID := uuid.New()
		seedItem(t, itemID, 5)
		ids := seedQueuedBatch(t, itemID, 100)

		usecase := newAdvanceQueueUsecase()

		var wg sync.WaitGroup
		errs := make(chan error, 100)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := usecase.AdvanceQueue(context.Background(), itemID, 3*time.Second)
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)

		for err := range errs {
			require.NoError(t, err)
		}

		assert.Equal(t, 5, countTicketsByStatus(t, itemID, domain.QueueStatusOffered))
		assert.Equal(t, 95, countTicketsByStatus(t, itemID, domain.QueueStatusQueued))

		for i, id := range ids {
			wantOffered := i < 5
			status := ticketStatus(t, id)
			if wantOffered {
				assert.Equal(t, domain.QueueStatusOffered, status, "ticket %d (earliest) should be OFFERED", i)
			} else {
				assert.Equal(t, domain.QueueStatusQueued, status, "ticket %d should remain QUEUED", i)
			}
		}
	})
}

func TestAdvanceQueueUsecase_ExpireAndPromoteInSingleCall(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	now := time.Now().UTC()
	past := now.Add(-time.Second)

	expiredOffered := uuid.New()
	seedTicket(t, expiredOffered, itemID, uuid.New(), domain.QueueStatusOffered, now.Add(-time.Minute), &past)

	next := uuid.New()
	seedTicket(t, next, itemID, uuid.New(), domain.QueueStatusQueued, now, nil)

	usecase := newAdvanceQueueUsecase()
	err := usecase.AdvanceQueue(context.Background(), itemID, 3*time.Second)
	require.NoError(t, err)

	assert.Equal(t, domain.QueueStatusExpired, ticketStatus(t, expiredOffered))
	assert.Equal(t, domain.QueueStatusOffered, ticketStatus(t, next))
}

func TestAdvanceQueueUsecase_SoldOut(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusPurchased, time.Now().UTC(), nil)
	ids := seedQueuedBatch(t, itemID, 3)

	usecase := newAdvanceQueueUsecase()
	err := usecase.AdvanceQueue(context.Background(), itemID, 3*time.Second)
	require.NoError(t, err)

	for _, id := range ids {
		assert.Equal(t, domain.QueueStatusSoldOut, ticketStatus(t, id))
	}
}

func TestAdvanceQueueUsecase_Idempotent(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)
	seedQueuedBatch(t, itemID, 3)

	usecase := newAdvanceQueueUsecase()
	ctx := context.Background()

	err := usecase.AdvanceQueue(ctx, itemID, 3*time.Second)
	require.NoError(t, err)

	offeredAfterFirst := countTicketsByStatus(t, itemID, domain.QueueStatusOffered)
	queuedAfterFirst := countTicketsByStatus(t, itemID, domain.QueueStatusQueued)

	err = usecase.AdvanceQueue(ctx, itemID, 3*time.Second)
	require.NoError(t, err)

	assert.Equal(t, offeredAfterFirst, countTicketsByStatus(t, itemID, domain.QueueStatusOffered))
	assert.Equal(t, queuedAfterFirst, countTicketsByStatus(t, itemID, domain.QueueStatusQueued))
}

func TestAdvanceQueueUsecase_ItemNotFound(t *testing.T) {
	truncate(db, t)

	usecase := newAdvanceQueueUsecase()
	err := usecase.AdvanceQueue(context.Background(), uuid.New(), 3*time.Second)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrItemNotFound)
}
