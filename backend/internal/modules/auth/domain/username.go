package domain

import "strings"

type Username string

const (
	usernameMinLength = 3
	usernameMaxLength = 32
)

func NewUsername(value string) (Username, error) {
	value = strings.ToLower(strings.TrimSpace(value))

	if value == "" {
		return "", ErrUsernameRequired
	}

	if len(value) < usernameMinLength || len(value) > usernameMaxLength {
		return "", ErrUsernameLength
	}

	return Username(value), nil
}

func (u Username) String() string {
	return string(u)
}
