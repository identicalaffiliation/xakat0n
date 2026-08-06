package domain

import "errors"

var (
	ErrItemNotFound = errors.New("item not found")
	ErrInternal     = errors.New("internal server error")
)
