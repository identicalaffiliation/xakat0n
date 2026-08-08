package domain

import (
	"time"

	"github.com/google/uuid"
)

// Queue — узкая проекция тикета для checkout-модуля.
type Queue struct {
	ID        uuid.UUID
	ItemID    uuid.UUID
	UserID    uuid.UUID
	Status    QueueStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt *time.Time
}
