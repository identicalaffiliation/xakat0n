package domain

import "errors"

var (
	ErrUserAlreadyQueued = errors.New("user already queued")
	ErrQueueNotFound     = errors.New("queue not found")
	ErrInternal          = errors.New("internal server error")
)
