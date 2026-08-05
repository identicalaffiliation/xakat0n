package tx

import (
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	t.Parallel()
	testPool := &pgxpool.Pool{}
	manager := NewManager(testPool, slog.Default())
	assert.NotNil(t, manager)
}
