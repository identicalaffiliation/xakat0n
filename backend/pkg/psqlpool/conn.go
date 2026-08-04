package psqlpool

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/xakat0n/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, func(), error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DbUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("parse postgres config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MaxConnLifetime = cfg.MaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("ping postgres pool: %w", err)
	}

	cleanup := func() {
		pool.Close()
	}

	return pool, cleanup, nil
}
