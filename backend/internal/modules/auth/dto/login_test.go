package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoginRequest(t *testing.T) {
	t.Parallel()

	actual := NewLoginRequest("alice")
	require.NotNil(t, actual)
	assert.Equal(t, "alice", actual.Username)
}

func TestNewLoginResponse(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actual := NewLoginResponse(userID, "alice", "jwt-token")
	require.NotNil(t, actual)

	assert.Equal(t, userID, actual.UserID)
	assert.Equal(t, "alice", actual.Username)
	assert.Equal(t, "jwt-token", actual.Token)
}
