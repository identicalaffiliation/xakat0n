package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type UsersRepository struct {
	pool   tx.DBTX
	logger ports.Logger
}

func NewUsersRepository(pool tx.DBTX, loggers ...ports.Logger) *UsersRepository {
	var logger ports.Logger
	if len(loggers) > 0 {
		logger = loggers[0]
	}
	return &UsersRepository{
		pool:   pool,
		logger: logger,
	}
}

// GetOrCreate — атомарный upsert вместо SELECT затем INSERT: под параллельными
// логинами одним и тем же username гонки не будет, ON CONFLICT либо создаёт
// строку, либо возвращает id уже существующей.
func (repo *UsersRepository) GetOrCreate(ctx context.Context, username domain.Username) (uuid.UUID, error) {
	const query = `
		INSERT INTO users (username) VALUES ($1)
		ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
		RETURNING id`

	var userID uuid.UUID
	if err := repo.pool.QueryRow(ctx, query, username.String()).Scan(&userID); err != nil {
		err = fmt.Errorf("get or create user: %w", err)
		if repo.logger != nil {
			err = repo.logger.WrapError(ctx, err)
		}
		return uuid.Nil, err
	}

	return userID, nil
}
