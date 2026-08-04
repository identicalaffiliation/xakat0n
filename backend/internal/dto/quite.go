package dto

import "github.com/google/uuid"

type QuiteQueueRequest struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}

func NewQuiteQueueRequest(productID uuid.UUID) *QuiteQueueRequest {
	return &QuiteQueueRequest{ProductID: productID}
}
