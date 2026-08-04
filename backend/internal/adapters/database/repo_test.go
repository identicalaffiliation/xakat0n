package database

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCreateQueueRepository(t *testing.T) {
	t.Parallel()
	testPool := &pgxpool.Pool{}
	repo := NewQueueRepository(testPool)
	assert.NotNil(t, repo)
}

func TestCreateTxManager(t *testing.T) {
	t.Parallel()
	testPool := &pgxpool.Pool{}
	manager := NewTxManager(testPool, slog.Default())
	assert.NotNil(t, manager)
}

func TestCheckUnique(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, checkUniqueViolation(&pgconn.PgError{Code: UniqueViolationCode}))
	assert.Equal(t, false, checkUniqueViolation(errors.New("an error")))
}
