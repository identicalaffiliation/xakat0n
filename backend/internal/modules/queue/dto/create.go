package dto

import "github.com/google/uuid"

type CreateQueueRequest struct {
	ItemID uuid.UUID `json:"itemId" validate:"required"`
	UserID uuid.UUID `json:"userId" validate:"required"`
}

func NewCreateRequest(itemID, userID uuid.UUID) *CreateQueueRequest {
	return &CreateQueueRequest{
		ItemID: itemID,
		UserID: userID,
	}
}
