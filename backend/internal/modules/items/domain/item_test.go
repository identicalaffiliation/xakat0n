package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewItem(t *testing.T) {
	require.NotNil(t, NewItem("a", "b", 1, true))
}
