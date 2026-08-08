package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
)

func TestNewTicket(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	t.Run("checkout fills expiresInSeconds", func(t *testing.T) {
		t.Parallel()

		expiresAt := now.Add(42 * time.Second)
		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID:    uuid.New(),
			Status:    domain.QueueStatusCheckout,
			ExpiresAt: &expiresAt,
			CreatedAt: now,
		}

		ticket := NewTicket(queue, now)
		require.NotNil(t, ticket.ExpiresInSeconds)
		assert.Equal(t, int64(42), *ticket.ExpiresInSeconds)
	})

	t.Run("purchased leaves expiresInSeconds nil even with expiresAt set", func(t *testing.T) {
		t.Parallel()

		expiresAt := now.Add(42 * time.Second)
		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID:    uuid.New(),
			Status:    domain.QueueStatusPurchased,
			ExpiresAt: &expiresAt,
			CreatedAt: now,
		}

		ticket := NewTicket(queue, now)
		assert.Nil(t, ticket.ExpiresInSeconds)
	})

	t.Run("clamps already-expired window to zero, not negative", func(t *testing.T) {
		t.Parallel()

		expiresAt := now.Add(-5 * time.Second)
		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID:    uuid.New(),
			Status:    domain.QueueStatusCheckout,
			ExpiresAt: &expiresAt,
			CreatedAt: now,
		}

		ticket := NewTicket(queue, now)
		require.NotNil(t, ticket.ExpiresInSeconds)
		assert.Equal(t, int64(0), *ticket.ExpiresInSeconds)
	})

	t.Run("position and nextSlotFreeInSeconds are always nil", func(t *testing.T) {
		t.Parallel()

		queue := &domain.Queue{
			ID:        uuid.New(),
			ItemID:    uuid.New(),
			Status:    domain.QueueStatusCheckout,
			CreatedAt: now,
		}

		ticket := NewTicket(queue, now)
		assert.Nil(t, ticket.Position)
		assert.Nil(t, ticket.NextSlotFreeInSeconds)
	})
}
