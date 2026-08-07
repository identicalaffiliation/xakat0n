package ports

import (
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
)

type TokenIssuer interface {
	Issue(userID uuid.UUID, username domain.Username) (string, error)
}
