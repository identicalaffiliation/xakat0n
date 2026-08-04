package dto

import (
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
)

type CreateQueueRequest struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
	UserID    uuid.UUID `json:"userId" validate:"required"`
}

func NewCreateRequest(productID, userID uuid.UUID) *CreateQueueRequest {
	return &CreateQueueRequest{
		ProductID: productID,
		UserID:    userID,
	}
}

type CreateQueueResponse struct {
	Queue Queue `json:"queue"`
}

func NewCreateResponse(queue *domain.Queue) *CreateQueueResponse {
	return &CreateQueueResponse{
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
