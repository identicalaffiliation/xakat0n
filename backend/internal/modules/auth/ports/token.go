package ports

import "github.com/google/uuid"

type TokenIssuer interface {
	Issue(userID uuid.UUID, username string) (string, error)
}
