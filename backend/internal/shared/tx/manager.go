package tx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	pool   *pgxpool.Pool
	logger txLogger
}

type txLogger interface {
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type errorWrapper interface {
	WrapError(ctx context.Context, err error) error
}

func NewManager(pool *pgxpool.Pool, logger txLogger) *Manager {
	return &Manager{
		pool:   pool,
		logger: logger,
	}
}

func (m *Manager) wrapError(ctx context.Context, err error) error {
	wrapper, ok := m.logger.(errorWrapper)
	if !ok {
		return err
	}
	return wrapper.WrapError(ctx, err)
}

func (m *Manager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return m.wrapError(ctx, fmt.Errorf("begin tx: %w", err))
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			m.logger.ErrorContext(ctx,
				"failed to rollback tx",
				"error", err,
			)
		}
	}()

	if err := fn(withTx(ctx, tx)); err != nil {
		return m.wrapError(ctx, fmt.Errorf("do fn: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return m.wrapError(ctx, fmt.Errorf("commit tx: %w", err))
	}

	return nil
}
