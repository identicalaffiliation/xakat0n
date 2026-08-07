package integrations

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	checkout "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout"
	checkoutapplication "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/application"
	checkoutdomain "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
	itemspostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
)

func newCheckoutUsecase(ttl time.Duration) *checkoutapplication.CheckoutUsecase {
	queueRepo := postgres.NewQueueRepository(db)
	itemsRepo := itemspostgres.NewItemsRepository(db)
	advance := checkout.NewAdvanceAdapter(newAdvanceQueueUsecase())

	return checkoutapplication.NewCheckoutUsecase(advance, itemsRepo, queueRepo, slog.Default(), ttl)
}

func TestCheckoutUsecase_Success(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	userID := uuid.New()
	expiresAt := time.Now().UTC().Add(30 * time.Second)
	seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusOffered, time.Now().UTC(), &expiresAt)

	usecase := newCheckoutUsecase(3 * time.Second)
	result, err := usecase.StartCheckout(context.Background(), itemID, userID)
	require.NoError(t, err)

	assert.True(t, result.QueueApplied)
	require.NotNil(t, result.Ticket)
	assert.Equal(t, checkoutdomain.QueueStatusCheckout, result.Ticket.Status)
	// Окно не перезапускается: expiresAt тот же, что был выставлен при выдаче права.
	assert.WithinDuration(t, expiresAt, *result.Ticket.ExpiresAt, time.Millisecond)
}

func TestCheckoutUsecase_NonLimitedItemSkipsQueue(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	_, err := db.Exec(
		context.Background(),
		`INSERT INTO items (id, title, price, is_limited, stock) VALUES ($1, 'not limited', 100, false, 1)`,
		itemID,
	)
	require.NoError(t, err)

	usecase := newCheckoutUsecase(3 * time.Second)
	result, err := usecase.StartCheckout(context.Background(), itemID, uuid.New())
	require.NoError(t, err)

	assert.False(t, result.QueueApplied)
	assert.Nil(t, result.Ticket)
}

func TestCheckoutUsecase_NoActiveRight(t *testing.T) {
	itemID := uuid.New()
	usecase := newCheckoutUsecase(3 * time.Second)

	t.Run("no ticket at all", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		_, err := usecase.StartCheckout(context.Background(), itemID, uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrNoActiveRight)
	})

	t.Run("still queued, never offered", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		// Слот занят другим держателем, иначе advanceQueue (вызывается первым шагом
		// StartCheckout) сам продвинул бы этот QUEUED-тикет в OFFERED, и checkout
		// в этом же вызове прошёл бы — тест проверял бы не тот сценарий.
		holderExpiresAt := time.Now().UTC().Add(30 * time.Second)
		seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusOffered, time.Now().UTC(), &holderExpiresAt)

		userID := uuid.New()
		seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusQueued, time.Now().UTC().Add(time.Second), nil)

		_, err := usecase.StartCheckout(context.Background(), itemID, userID)
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrNoActiveRight)
	})

	t.Run("offered but already expired", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		userID := uuid.New()
		expiredAt := time.Now().UTC().Add(-time.Second)
		seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusOffered, time.Now().UTC().Add(-time.Minute), &expiredAt)

		_, err := usecase.StartCheckout(context.Background(), itemID, userID)
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrNoActiveRight)
	})

	t.Run("belongs to another user", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		expiresAt := time.Now().UTC().Add(30 * time.Second)
		seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusOffered, time.Now().UTC(), &expiresAt)

		_, err := usecase.StartCheckout(context.Background(), itemID, uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrNoActiveRight)
	})
}

func TestCheckoutUsecase_ItemNotFound(t *testing.T) {
	truncate(db, t)

	usecase := newCheckoutUsecase(3 * time.Second)
	_, err := usecase.StartCheckout(context.Background(), uuid.New(), uuid.New())
	require.Error(t, err)
	assert.ErrorIs(t, err, checkoutdomain.ErrItemNotFound)
}
