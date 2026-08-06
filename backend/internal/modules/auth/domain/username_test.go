package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUsername(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		expected  Username
		expectErr error
	}{
		{
			name:     "valid username",
			value:    "aleksei",
			expected: Username("aleksei"),
		},
		{
			name:     "normalizes case and whitespace",
			value:    "  Aleksei  ",
			expected: Username("aleksei"),
		},
		{
			name:     "minimum length",
			value:    "abc",
			expected: Username("abc"),
		},
		{
			name:     "maximum length",
			value:    strings.Repeat("a", usernameMaxLength),
			expected: Username(strings.Repeat("a", usernameMaxLength)),
		},
		{
			name:      "empty username",
			value:     "",
			expectErr: ErrUsernameRequired,
		},
		{
			name:      "whitespace-only username",
			value:     "   ",
			expectErr: ErrUsernameRequired,
		},
		{
			name:      "too short",
			value:     "ab",
			expectErr: ErrUsernameLength,
		},
		{
			name:      "too long",
			value:     strings.Repeat("a", usernameMaxLength+1),
			expectErr: ErrUsernameLength,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := NewUsername(test.value)
			if test.expectErr != nil {
				require.ErrorIs(t, err, test.expectErr)
				assert.Empty(t, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestUsernameString(t *testing.T) {
	t.Parallel()

	username := Username("aleksei")

	assert.Equal(t, "aleksei", username.String())
}
