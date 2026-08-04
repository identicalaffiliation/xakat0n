package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCreateRequest(t *testing.T) {
	t.Parallel()

	id1, id2 := uuid.New(), uuid.New()
	actual := NewCreateRequest(id1, id2)
	require.NotNil(t, actual)

	assert.Equal(t, actual.UserID, id2)
	assert.Equal(t, actual.ProductID, id1)
}

func TestCreateCreateResponse(t *testing.T) {
	t.Parallel()

	queue := domain.Queue{
		ID:        uuid.New(),
		ProductID: uuid.New(),
		UserID:    uuid.New(),
		Status:    domain.QueueStatusCheckout,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	actual := NewCreateResponse(&queue)
	require.NotNil(t, actual)

	assert.Equal(t, queue.ID, actual.Queue.ID)
	assert.Equal(t, queue.ProductID, actual.Queue.ProductID)
	assert.Equal(t, queue.UserID, actual.Queue.UserID)
	assert.Equal(t, queue.Status, actual.Queue.Status)
	assert.Equal(t, queue.CreatedAt, actual.Queue.CreatedAt)
	assert.Equal(t, queue.UpdatedAt, actual.Queue.UpdatedAt)
}
