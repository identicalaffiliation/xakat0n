package tx

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestDBTXFromContext(t *testing.T) {
	t.Parallel()

	fallback := &pgxpool.Pool{}
	assert.Same(t, DBTX(fallback), DBTXFromContext(context.Background(), fallback))

	txDB := &pgxpool.Pool{}
	ctx := withTx(context.Background(), txDB)
	assert.Same(t, DBTX(txDB), DBTXFromContext(ctx, fallback))
}
