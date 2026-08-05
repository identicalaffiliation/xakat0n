package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsUniqueViolation(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, IsUniqueViolation(&pgconn.PgError{Code: UniqueViolationCode}))
	assert.Equal(t, false, IsUniqueViolation(errors.New("an error")))
}
