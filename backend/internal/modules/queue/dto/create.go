package dto

import "github.com/google/uuid"

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
