package domain

import "errors"

var (
	ErrUserAlreadyQueued = errors.New("user already queued")
	ErrInternal          = errors.New("internal server error")
)
