package domain

import (
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	UserID    uuid.UUID
	Status    QueueStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt *time.Time
}

func NewQueue(productID, userID uuid.UUID) *Queue {
	return &Queue{
		ID:        uuid.New(),
		ProductID: productID,
		UserID:    userID,
	}
}
