package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type UsersRepository struct {
	pool tx.DBTX
}

func NewUsersRepository(pool tx.DBTX) *UsersRepository {
	return &UsersRepository{
		pool: pool,
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
		return uuid.Nil, fmt.Errorf("get or create user: %w", err)
	}

	return userID, nil
}
