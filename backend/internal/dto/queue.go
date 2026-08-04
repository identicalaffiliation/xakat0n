package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
)

type Queue struct {
	ID        uuid.UUID          `json:"id"`
	ProductID uuid.UUID          `json:"productId"`
	UserID    uuid.UUID          `json:"userId"`
	Status    domain.QueueStatus `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	ExpiresAt *time.Time         `json:"expiresAt,omitempty"`
}

type QueueResponse struct {
	Queue Queue `json:"queue"`
}

func NewCreateResponse(queue *domain.Queue) *QueueResponse {
	return &QueueResponse{
		Queue: Queue{
			ID:        queue.ID,
			ProductID: queue.ProductID,
			UserID:    queue.UserID,
			Status:    queue.Status,
			CreatedAt: queue.CreatedAt,
			UpdatedAt: queue.UpdatedAt,
			ExpiresAt: queue.ExpiresAt,
		},
	}
}
