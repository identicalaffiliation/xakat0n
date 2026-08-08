package domain

import "errors"

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrUsernameLength   = errors.New("username must contain from 3 to 32 characters")
	ErrUsernameFormat   = errors.New("username has invalid format")
)
