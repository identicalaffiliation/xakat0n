package domain

import "errors"

var (
	ErrItemNotFound   = errors.New("item not found")
	ErrTicketNotFound = errors.New("ticket not found")
	ErrNoActiveRight  = errors.New("no active right")
	ErrTooLate        = errors.New("checkout result too late")
	ErrInternal       = errors.New("internal server error")
)
