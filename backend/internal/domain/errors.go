package domain

import "errors"

var (
	UserAlreadyQueued = errors.New("user already queued")
	ErrInternal       = errors.New("internal server error")
)
