package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
)

type UserRepository interface {
	GetOrCreate(ctx context.Context, username domain.Username) (uuid.UUID, error)
}
