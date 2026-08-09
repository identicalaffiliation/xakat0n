package integrations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
)

func ticketStatus(t *testing.T, ticketID uuid.UUID) domain.QueueStatus {
	t.Helper()

	var status domain.QueueStatus
	err := db.QueryRow(context.Background(), `SELECT status FROM queues WHERE id = $1`, ticketID).Scan(&status)
	require.NoError(t, err)

	return status
}

func TestQueueRepository_ExpireStale(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := postgres.NewQueueRepository(db)
	itemID := uuid.New()
	now := time.Now().UTC()
	past := now.Add(-time.Minute)
	future := now.Add(time.Minute)

	expiredOffered := uuid.New()
	seedTicket(t, expiredOffered, itemID, uuid.New(), domain.QueueStatusOffered, now, &past)

	expiredCheckout := uuid.New()
	seedTicket(t, expiredCheckout, itemID, uuid.New(), domain.QueueStatusCheckout, now, &past)

	freshOffered := uuid.New()
	seedTicket(t, freshOffered, itemID, uuid.New(), domain.QueueStatusOffered, now, &future)

	queued := uuid.New()
	seedTicket(t, queued, itemID, uuid.New(), domain.QueueStatusQueued, now, nil)

	require.NoError(t, repo.ExpireStale(ctx, itemID))

	assert.Equal(t, domain.QueueStatusExpired, ticketStatus(t, expiredOffered))
	assert.Equal(t, domain.QueueStatusExpired, ticketStatus(t, expiredCheckout))
	assert.Equal(t, domain.QueueStatusOffered, ticketStatus(t, freshOffered))
	assert.Equal(t, domain.QueueStatusQueued, ticketStatus(t, queued))
}

func TestQueueRepository_CountTaken(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := postgres.NewQueueRepository(db)
	itemID := uuid.New()
	now := time.Now().UTC()

	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusOffered, now, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusCheckout, now, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusPurchased, now, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusQueued, now, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusExpired, now, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusCancelled, now, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusSoldOut, now, nil)

	taken, err := repo.CountTaken(ctx, itemID)
	require.NoError(t, err)
	assert.Equal(t, 3, taken)
}

func TestQueueRepository_CountPurchased(t *testing.T) {
	t.Run("counts only PURCHASED, grouped per item, omits items with none", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		now := time.Now().UTC()

		itemWithTwoPurchases := uuid.New()
		seedTicket(t, uuid.New(), itemWithTwoPurchases, uuid.New(), domain.QueueStatusPurchased, now, nil)
		seedTicket(t, uuid.New(), itemWithTwoPurchases, uuid.New(), domain.QueueStatusPurchased, now, nil)
		seedTicket(t, uuid.New(), itemWithTwoPurchases, uuid.New(), domain.QueueStatusQueued, now, nil)

		itemWithNoPurchases := uuid.New()
		seedTicket(t, uuid.New(), itemWithNoPurchases, uuid.New(), domain.QueueStatusOffered, now, nil)

		itemNotRequested := uuid.New()
		seedTicket(t, uuid.New(), itemNotRequested, uuid.New(), domain.QueueStatusPurchased, now, nil)

		counts, err := repo.CountPurchased(ctx, []uuid.UUID{itemWithTwoPurchases, itemWithNoPurchases})
		require.NoError(t, err)

		assert.Equal(t, 2, counts[itemWithTwoPurchases])
		_, ok := counts[itemWithNoPurchases]
		assert.False(t, ok)
		_, ok = counts[itemNotRequested]
		assert.False(t, ok)
	})

	t.Run("empty input returns empty map without querying", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)

		counts, err := repo.CountPurchased(ctx, nil)
		require.NoError(t, err)
		assert.Empty(t, counts)
	})
}

func TestQueueRepository_PromoteNext(t *testing.T) {
	t.Run("promotes exactly freeSlots earliest by created_at", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		itemID := uuid.New()
		base := time.Now().UTC().Add(-time.Hour)

		ids := make([]uuid.UUID, 5)
		for i := range ids {
			ids[i] = uuid.New()
			seedTicket(t, ids[i], itemID, uuid.New(), domain.QueueStatusQueued, base.Add(time.Duration(i)*time.Millisecond), nil)
		}

		affected, err := repo.PromoteNext(ctx, itemID, 2, 3*time.Second)
		require.NoError(t, err)
		assert.Equal(t, int64(2), affected)

		assert.Equal(t, domain.QueueStatusOffered, ticketStatus(t, ids[0]))
		assert.Equal(t, domain.QueueStatusOffered, ticketStatus(t, ids[1]))
		assert.Equal(t, domain.QueueStatusQueued, ticketStatus(t, ids[2]))
		assert.Equal(t, domain.QueueStatusQueued, ticketStatus(t, ids[3]))
		assert.Equal(t, domain.QueueStatusQueued, ticketStatus(t, ids[4]))
	})

	t.Run("SKIP LOCKED skips a row locked by a concurrent uncommitted tx", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		itemID := uuid.New()
		base := time.Now().UTC().Add(-time.Hour)

		earlier := uuid.New()
		seedTicket(t, earlier, itemID, uuid.New(), domain.QueueStatusQueued, base, nil)
		later := uuid.New()
		seedTicket(t, later, itemID, uuid.New(), domain.QueueStatusQueued, base.Add(time.Second), nil)

		lockingTx, err := db.Begin(ctx)
		require.NoError(t, err)
		defer func() { _ = lockingTx.Rollback(ctx) }()

		var lockedID uuid.UUID
		require.NoError(t, lockingTx.QueryRow(ctx, `SELECT id FROM queues WHERE id = $1 FOR UPDATE`, earlier).Scan(&lockedID))

		affected, err := repo.PromoteNext(ctx, itemID, 1, 3*time.Second)
		require.NoError(t, err)
		assert.Equal(t, int64(1), affected)

		assert.Equal(t, domain.QueueStatusQueued, ticketStatus(t, earlier))
		assert.Equal(t, domain.QueueStatusOffered, ticketStatus(t, later))
	})
}

func TestQueueRepository_MarkSoldOut(t *testing.T) {
	t.Run("purchased reached stock", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		itemID := uuid.New()
		now := time.Now().UTC()

		seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusPurchased, now, nil)
		queuedA := uuid.New()
		seedTicket(t, queuedA, itemID, uuid.New(), domain.QueueStatusQueued, now, nil)
		queuedB := uuid.New()
		seedTicket(t, queuedB, itemID, uuid.New(), domain.QueueStatusQueued, now, nil)

		require.NoError(t, repo.MarkSoldOut(ctx, itemID, 1))

		assert.Equal(t, domain.QueueStatusSoldOut, ticketStatus(t, queuedA))
		assert.Equal(t, domain.QueueStatusSoldOut, ticketStatus(t, queuedB))
	})

	t.Run("purchased below stock leaves queued untouched", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		itemID := uuid.New()
		now := time.Now().UTC()

		seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusPurchased, now, nil)
		queued := uuid.New()
		seedTicket(t, queued, itemID, uuid.New(), domain.QueueStatusQueued, now, nil)

		require.NoError(t, repo.MarkSoldOut(ctx, itemID, 5))

		assert.Equal(t, domain.QueueStatusQueued, ticketStatus(t, queued))
	})
}

func TestQueueRepository_GetLatestTicket(t *testing.T) {
	t.Run("returns the most recent ticket, not the oldest", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		itemID := uuid.New()
		userID := uuid.New()
		base := time.Now().UTC().Add(-time.Hour)

		seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusCancelled, base, nil)
		seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusExpired, base.Add(time.Minute), nil)
		latest := uuid.New()
		seedTicket(t, latest, itemID, userID, domain.QueueStatusQueued, base.Add(2*time.Minute), nil)

		ticket, err := repo.GetLatestTicket(ctx, itemID, userID)
		require.NoError(t, err)
		assert.Equal(t, latest, ticket.ID)
		assert.Equal(t, domain.QueueStatusQueued, ticket.Status)
	})

	t.Run("error - no ticket for pair", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)

		_, err := repo.GetLatestTicket(ctx, uuid.New(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrTicketNotFound)
	})
}

func TestQueueRepository_CountQueuedAhead(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := postgres.NewQueueRepository(db)
	itemID := uuid.New()
	base := time.Now().UTC().Add(-time.Hour)

	first := base
	second := base.Add(time.Minute)
	third := base.Add(2 * time.Minute)

	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusQueued, first, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusQueued, second, nil)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusQueued, third, nil)

	ahead, err := repo.CountQueuedAhead(ctx, itemID, third)
	require.NoError(t, err)
	assert.Equal(t, 2, ahead)

	ahead, err = repo.CountQueuedAhead(ctx, itemID, first)
	require.NoError(t, err)
	assert.Equal(t, 0, ahead)
}

func TestQueueRepository_NextSlotFreeAt(t *testing.T) {
	t.Run("returns the earliest active deadline", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		itemID := uuid.New()
		now := time.Now().UTC()
		soon := now.Add(5 * time.Second)
		later := now.Add(10 * time.Second)

		seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusOffered, now, &later)
		seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusCheckout, now, &soon)

		nextFree, err := repo.NextSlotFreeAt(ctx, itemID)
		require.NoError(t, err)
		require.NotNil(t, nextFree)
		assert.WithinDuration(t, soon, *nextFree, time.Second)
	})

	t.Run("nil when no active rights", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewQueueRepository(db)
		itemID := uuid.New()

		seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusQueued, time.Now().UTC(), nil)

		nextFree, err := repo.NextSlotFreeAt(ctx, itemID)
		require.NoError(t, err)
		assert.Nil(t, nextFree)
	})
}
