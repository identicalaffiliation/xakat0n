package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
)

func TestNewItem(t *testing.T) {
	t.Parallel()

	t.Run("limited item exposes stock and soldOut", func(t *testing.T) {
		id := uuid.New()
		item := &domain.Item{
			ID:          id,
			Title:       "PS5",
			Description: "cons",
			Price:       90000,
			Stock:       1,
			IsLimited:   true,
			Category:    nil,
		}

		actual := NewItem(item, true)

		assert.Equal(t, id, actual.ItemID)
		assert.Equal(t, "PS5", actual.Title)
		require.NotNil(t, actual.Description)
		assert.Equal(t, "cons", *actual.Description)
		assert.Equal(t, int64(90000), actual.Price)
		require.NotNil(t, actual.Stock)
		assert.Equal(t, 1, *actual.Stock)
		assert.True(t, actual.IsLimited)
		assert.Nil(t, actual.Category)
		assert.True(t, actual.SoldOut)
	})

	t.Run("non-limited item has nil stock regardless of db value", func(t *testing.T) {
		item := &domain.Item{
			ID:        uuid.New(),
			Title:     "mug",
			Price:     500,
			Stock:     1,
			IsLimited: false,
		}

		actual := NewItem(item, false)

		assert.False(t, actual.IsLimited)
		assert.Nil(t, actual.Stock)
		assert.False(t, actual.SoldOut)
	})
}

func TestNewItems(t *testing.T) {
	t.Parallel()

	t.Run("empty catalog", func(t *testing.T) {
		actual := NewItems(nil, nil)
		assert.NotNil(t, actual)
		assert.Empty(t, actual)
	})

	t.Run("maps every item preserving order and soldOut", func(t *testing.T) {
		first := uuid.New()
		second := uuid.New()
		items := []*domain.Item{
			{ID: first, Title: "first", Price: 1, IsLimited: true, Stock: 1},
			{ID: second, Title: "second", Price: 2},
		}
		soldOut := map[uuid.UUID]bool{first: true}

		actual := NewItems(items, soldOut)

		require.Len(t, actual, 2)
		assert.Equal(t, first, actual[0].ItemID)
		assert.Equal(t, "first", actual[0].Title)
		assert.True(t, actual[0].SoldOut)
		assert.Equal(t, second, actual[1].ItemID)
		assert.Equal(t, "second", actual[1].Title)
		assert.False(t, actual[1].SoldOut)
	})
}
