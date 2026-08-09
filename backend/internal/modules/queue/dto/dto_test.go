package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

func TestCreateCreateRequest(t *testing.T) {
	t.Parallel()

	id1, id2 := uuid.New(), uuid.New()
	actual := NewCreateRequest(id1, id2)
	require.NotNil(t, actual)

	assert.Equal(t, actual.UserID, id2)
	assert.Equal(t, actual.ItemID, id1)
}

func TestNewTicket(t *testing.T) {
	t.Parallel()

	queueID, itemID := uuid.New(), uuid.New()
	now := time.Now().UTC()

	t.Run("offered sets active window fields", func(t *testing.T) {
		expiresAt := now.Add(90 * time.Second)
		queue := &domain.Queue{
			ID:        queueID,
			ItemID: itemID,
			Status:    domain.QueueStatusOffered,
			ExpiresAt: &expiresAt,
			CreatedAt: now.Add(-time.Minute),
		}

		ticket := NewTicket(queue, now)

		assert.Equal(t, queueID, ticket.TicketID)
		assert.Equal(t, itemID, ticket.ItemID)
		assert.Equal(t, domain.QueueStatusOffered, ticket.Status)
		assert.Equal(t, now, ticket.ServerTime)
		assert.Equal(t, queue.CreatedAt, ticket.CreatedAt)
		assert.Equal(t, &expiresAt, ticket.ExpiresAt)
		require.NotNil(t, ticket.ExpiresInSeconds)
		assert.InDelta(t, int64(90), *ticket.ExpiresInSeconds, 1)
		assert.Nil(t, ticket.Position)
		assert.Nil(t, ticket.NextSlotFreeInSeconds)
	})

	t.Run("checkout also has active window", func(t *testing.T) {
		expiresAt := now.Add(10 * time.Second)
		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID: itemID,
			Status:    domain.QueueStatusCheckout,
			ExpiresAt: &expiresAt,
		}

		ticket := NewTicket(queue, now)

		require.NotNil(t, ticket.ExpiresInSeconds)
		assert.InDelta(t, int64(10), *ticket.ExpiresInSeconds, 1)
	})

	t.Run("queued has no window", func(t *testing.T) {
		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID: itemID,
			Status:    domain.QueueStatusQueued,
		}

		ticket := NewTicket(queue, now)

		assert.Nil(t, ticket.ExpiresInSeconds)
		assert.Nil(t, ticket.ExpiresAt)
	})

	t.Run("already expired is clamped to zero", func(t *testing.T) {
		expiresAt := now.Add(-5 * time.Second)
		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID: itemID,
			Status:    domain.QueueStatusOffered,
			ExpiresAt: &expiresAt,
		}

		ticket := NewTicket(queue, now)

		require.NotNil(t, ticket.ExpiresInSeconds)
		assert.Equal(t, int64(0), *ticket.ExpiresInSeconds)
	})

	t.Run("terminal status without expires at", func(t *testing.T) {
		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID: itemID,
			Status:    domain.QueueStatusPurchased,
		}

		ticket := NewTicket(queue, now)

		assert.Nil(t, ticket.ExpiresInSeconds)
		assert.Nil(t, ticket.Position)
	})
}

func TestSetQueuedFields(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	t.Run("sets position and next slot estimate", func(t *testing.T) {
		nextFree := now.Add(42 * time.Second)
		ticket := NewTicket(&domain.Queue{
			ID:        uuid.New(),
			ItemID: uuid.New(),
			Status:    domain.QueueStatusQueued,
			CreatedAt: now,
		}, now)

		ticket.SetQueuedFields(2, &nextFree, now)

		require.NotNil(t, ticket.Position)
		assert.Equal(t, 3, *ticket.Position)
		require.NotNil(t, ticket.NextSlotFreeInSeconds)
		assert.InDelta(t, int64(42), *ticket.NextSlotFreeInSeconds, 1)
	})

	t.Run("first in line is position one", func(t *testing.T) {
		ticket := NewTicket(&domain.Queue{
			ID:        uuid.New(),
			ItemID: uuid.New(),
			Status:    domain.QueueStatusQueued,
			CreatedAt: now,
		}, now)

		ticket.SetQueuedFields(0, nil, now)

		require.NotNil(t, ticket.Position)
		assert.Equal(t, 1, *ticket.Position)
		assert.Nil(t, ticket.NextSlotFreeInSeconds)
	})
}

func TestClampSeconds(t *testing.T) {
	t.Parallel()

	t.Run("negative duration clamps to zero", func(t *testing.T) {
		result := clampSeconds(-10 * time.Second)
		require.NotNil(t, result)
		assert.Equal(t, int64(0), *result)
	})

	t.Run("zero duration stays zero", func(t *testing.T) {
		result := clampSeconds(0)
		require.NotNil(t, result)
		assert.Equal(t, int64(0), *result)
	})

	t.Run("positive duration is kept", func(t *testing.T) {
		result := clampSeconds(90 * time.Second)
		require.NotNil(t, result)
		assert.Equal(t, int64(90), *result)
	})

	t.Run("sub-second duration truncates", func(t *testing.T) {
		result := clampSeconds(500 * time.Millisecond)
		require.NotNil(t, result)
		assert.Equal(t, int64(0), *result)
	})
}
