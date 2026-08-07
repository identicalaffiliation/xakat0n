package domain

import "errors"

var (
	ErrUserAlreadyQueued = errors.New("user already queued")
	ErrInternal          = errors.New("internal server error")
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrItemNotFound      = errors.New("item not found")
)
