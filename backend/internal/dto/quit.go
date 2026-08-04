package dto

import "github.com/google/uuid"

type QuitQueueRequest struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
	UserID    uuid.UUID `json:"userId" validate:"required"`
}

func NewQuitQueueRequest(productID, userID uuid.UUID) *QuitQueueRequest {
	return &QuitQueueRequest{
		ProductID: productID,
		UserID:    userID,
	}
}
