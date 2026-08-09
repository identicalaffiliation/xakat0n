package domain

import (
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	ID        uuid.UUID
	ItemID uuid.UUID
	UserID    uuid.UUID
	Status    QueueStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt *time.Time
}

func NewQueue(itemID, userID uuid.UUID) *Queue {
	return &Queue{
		ID:        uuid.New(),
		ItemID: itemID,
		UserID:    userID,
	}
}
