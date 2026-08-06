package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type ItemsRepository struct {
	pool tx.DBTX
}

func NewItemsRepository(pool tx.DBTX) *ItemsRepository {
	return &ItemsRepository{
		pool: pool,
	}
}

func (repo *ItemsRepository) dbtx(ctx context.Context) tx.DBTX {
	return tx.DBTXFromContext(ctx, repo.pool)
}

func (repo *ItemsRepository) LockStock(ctx context.Context, itemID uuid.UUID) (*domain.Item, error) {
	const query = `SELECT id, stock, is_limited FROM items WHERE id = $1 FOR UPDATE`

	var item domain.Item
	err := repo.dbtx(ctx).QueryRow(ctx, query, itemID).Scan(&item.ID, &item.Stock, &item.IsLimited)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrItemNotFound
		}

		return nil, fmt.Errorf("lock stock: %w", err)
	}

	return &item, nil
}
