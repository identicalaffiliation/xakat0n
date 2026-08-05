package tx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/logger"
)

type Manager struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewManager(pool *pgxpool.Pool, l logger.Logger) *Manager {
	return &Manager{
		pool:   pool,
		logger: l,
	}
}

func (m *Manager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			m.logger.Error(
				"failed to rollback tx",
				"error", err,
			)
		}
	}()

	if err := fn(withTx(ctx, tx)); err != nil {
		return fmt.Errorf("do fn: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
