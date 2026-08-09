package integrations

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	itemspostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

func newCreateQueueUsecase() *application.CreateQueueUsecase {
	logging := newTestLogger()
	queueRepo := postgres.NewQueueRepository(db, logging)
	itemsRepo := itemspostgres.NewItemsRepository(db, logging)
	advanceUsecase := newAdvanceQueueUsecase(logging)
	txManager := tx.NewManager(db, logging)

	return application.NewCreateQueueUsecase(
		advanceUsecase,
		queueRepo,
		itemsRepo,
		txManager,
		logging,
		3*time.Second,
	)
}

func createTicket(t *testing.T, usecase *application.CreateQueueUsecase, itemID, userID uuid.UUID) (*dto.Ticket, error) {
	t.Helper()

	return usecase.CreateQueue(context.Background(), dto.NewCreateRequest(itemID, userID))
}

func TestCreateQueueUsecase_OfferedImmediately(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	usecase := newCreateQueueUsecase()
	ticket, err := createTicket(t, usecase, itemID, uuid.New())
	require.NoError(t, err)

	assert.Equal(t, domain.QueueStatusOffered, ticket.Status)
	assert.Nil(t, ticket.Position)
	require.NotNil(t, ticket.ExpiresInSeconds)
	assert.LessOrEqual(t, *ticket.ExpiresInSeconds, int64(3))
	assert.Greater(t, *ticket.ExpiresInSeconds, int64(0))
}

func TestCreateQueueUsecase_QueuedWhenSlotHeld(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	now := time.Now().UTC()
	holderExpiresAt := now.Add(30 * time.Second)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusOffered, now, &holderExpiresAt)

	usecase := newCreateQueueUsecase()
	ticket, err := createTicket(t, usecase, itemID, uuid.New())
	require.NoError(t, err)

	assert.Equal(t, domain.QueueStatusQueued, ticket.Status)
	require.NotNil(t, ticket.Position)
	assert.Equal(t, 1, *ticket.Position)
	assert.Nil(t, ticket.ExpiresInSeconds)
}

func TestCreateQueueUsecase_Idempotent(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	userID := uuid.New()
	usecase := newCreateQueueUsecase()

	first, err := createTicket(t, usecase, itemID, userID)
	require.NoError(t, err)
	require.Equal(t, domain.QueueStatusOffered, first.Status)

	second, err := createTicket(t, usecase, itemID, userID)
	require.NoError(t, err)

	assert.Equal(t, first.TicketID, second.TicketID)
	assert.Equal(t, domain.QueueStatusOffered, second.Status)
}

func TestCreateQueueUsecase_NotLimited(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	_, err := db.Exec(
		context.Background(),
		`INSERT INTO items (id, title, price, is_limited, stock) VALUES ($1, 'plain item', 100, false, 1)`,
		itemID,
	)
	require.NoError(t, err)

	usecase := newCreateQueueUsecase()
	_, err = createTicket(t, usecase, itemID, uuid.New())
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrQueueNotApplicable)
}

func TestCreateQueueUsecase_SoldOut(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	now := time.Now().UTC()
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusPurchased, now, nil)

	usecase := newCreateQueueUsecase()
	_, err := createTicket(t, usecase, itemID, uuid.New())
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrItemSoldOut)

	assert.Equal(t, 1, countTicketsByStatus(t, itemID, domain.QueueStatusPurchased))
}

func TestCreateQueueUsecase_ItemNotFound(t *testing.T) {
	truncate(db, t)

	usecase := newCreateQueueUsecase()
	_, err := createTicket(t, usecase, uuid.New(), uuid.New())
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrItemNotFound)
}

// TestCreateQueueUsecase_ReClickDoesNotJumpQueue — ключевая семантика architecture.md:
// у истёкшего права слот сначала уходит голове очереди (advanceQueue отдельной
// транзакцией до создания), и только потом рассматривается новая заявка.
func TestCreateQueueUsecase_ReClickDoesNotJumpQueue(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	now := time.Now().UTC()
	expiredUserID := uuid.New()
	expiredAt := now.Add(-30 * time.Second)
	seedTicket(t, uuid.New(), itemID, expiredUserID, domain.QueueStatusOffered, now.Add(-time.Minute), &expiredAt)

	waiterID := uuid.New()
	seedTicket(t, uuid.New(), itemID, waiterID, domain.QueueStatusQueued, now.Add(-20*time.Second), nil)

	usecase := newCreateQueueUsecase()
	ticket, err := createTicket(t, usecase, itemID, expiredUserID)
	require.NoError(t, err)

	assert.Equal(t, domain.QueueStatusQueued, ticket.Status)
	require.NotNil(t, ticket.Position)
	assert.Equal(t, 1, *ticket.Position)

	waiterTicket, err := newGetMyTicketUsecase().GetMyTicket(context.Background(), itemID, waiterID)
	require.NoError(t, err)
	assert.Equal(t, domain.QueueStatusOffered, waiterTicket.Status)
}

// TestCreateQueueUsecase_ParallelStock1 — NFR1: при 100 параллельных запросах
// на товар со stock=1 ровно один получает OFFERED, остальные — QUEUED.
func TestCreateQueueUsecase_ParallelStock1(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	usecase := newCreateQueueUsecase()

	const users = 100
	statuses := make([]domain.QueueStatus, users)
	errs := make([]error, users)

	var wg sync.WaitGroup
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ticket, err := createTicket(t, usecase, itemID, uuid.New())
			errs[i] = err
			if err == nil {
				statuses[i] = ticket.Status
			}
		}(i)
	}
	wg.Wait()

	offered, queued := 0, 0
	for i := 0; i < users; i++ {
		require.NoError(t, errs[i])
		switch statuses[i] {
		case domain.QueueStatusOffered:
			offered++
		case domain.QueueStatusQueued:
			queued++
		default:
			t.Fatalf("unexpected status %q", statuses[i])
		}
	}

	assert.Equal(t, 1, offered)
	assert.Equal(t, users-1, queued)
}
