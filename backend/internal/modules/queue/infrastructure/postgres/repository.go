package postgres

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type QueueRepository struct {
	pool tx.DBTX
}

func NewQueueRepository(pool tx.DBTX) *QueueRepository {
	return &QueueRepository{
		pool: pool,
	}
}

func (repo *QueueRepository) dbtx(ctx context.Context) tx.DBTX {
	return tx.DBTXFromContext(ctx, repo.pool)
}
