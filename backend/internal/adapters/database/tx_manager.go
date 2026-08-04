package database

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/xakat0n/backend/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	pool   *pgxpool.Pool
	logger ports.Logger
}

func NewTxManager(pool *pgxpool.Pool, slogger ports.Logger) *TxManager {
	return &TxManager{
		pool:   pool,
		logger: slogger,
	}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			m.logger.Error(
				"failed to rollback tx",
				"error", err,
			)
		}
	}()

	if err := fn(ctx); err != nil {
		return fmt.Errorf("do fn: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	
	return nil
}
