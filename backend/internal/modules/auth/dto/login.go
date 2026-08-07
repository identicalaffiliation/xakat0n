package dto

import "github.com/google/uuid"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
}

func NewLoginRequest(username string) *LoginRequest {
	return &LoginRequest{Username: username}
}

type LoginResponse struct {
	UserID   uuid.UUID `json:"userId" validate:"required"`
	Username string    `json:"username" validate:"required"`
	Token    string    `json:"token" validate:"required"`
}

func NewLoginResponse(userID uuid.UUID, username, token string) *LoginResponse {
	return &LoginResponse{
		UserID:   userID,
		Username: username,
		Token:    token,
	}
}
