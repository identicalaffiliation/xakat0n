package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateQueueDomain(t *testing.T) {
	t.Parallel()

	id1, id2 := uuid.New(), uuid.New()
	actual := NewQueue(id1, id2)
	assert.NotNil(t, actual)

	assert.Equal(t, id1, actual.ProductID)
	assert.Equal(t, id2, actual.UserID)
	assert.NotNil(t, actual.ID)
}
