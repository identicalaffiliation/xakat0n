package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCreateQueueRepository(t *testing.T) {
	t.Parallel()
	testPool := &pgxpool.Pool{}
	repo := NewQueueRepository(testPool)
	assert.NotNil(t, repo)
}
