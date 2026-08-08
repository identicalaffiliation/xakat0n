package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
)

func TestNewItemResponse(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	createdAt := time.Now().UTC()
	item := &domain.Item{
		ID:          id,
		Title:       "PS5",
		Description: "cons",
		Price:       90000,
		Stock:       1,
		IsLimited:   true,
		Category:    nil,
		CreatedAt:   createdAt,
	}

	actual := NewItemResponse(item)
	require.NotNil(t, actual)

	assert.Equal(t, id, actual.Item.ID)
	assert.Equal(t, "PS5", actual.Item.Title)
	require.NotNil(t, actual.Item.Description)
	assert.Equal(t, "cons", *actual.Item.Description)
	assert.Equal(t, int64(90000), actual.Item.Price)
	assert.Equal(t, 1, actual.Item.Stock)
	assert.True(t, actual.Item.IsLimited)
	assert.Nil(t, actual.Item.Category)
	assert.Equal(t, createdAt, actual.Item.CreatedAt)
}

func TestNewItemsResponse(t *testing.T) {
	t.Parallel()

	t.Run("empty catalog", func(t *testing.T) {
		actual := NewItemsResponse(nil)
		require.NotNil(t, actual)
		assert.NotNil(t, actual.Items)
		assert.Empty(t, actual.Items)
	})

	t.Run("maps every item preserving order", func(t *testing.T) {
		items := []*domain.Item{
			{ID: uuid.New(), Title: "first", Price: 1},
			{ID: uuid.New(), Title: "second", Price: 2},
		}

		actual := NewItemsResponse(items)

		require.Len(t, actual.Items, 2)
		assert.Equal(t, items[0].ID, actual.Items[0].Item.ID)
		assert.Equal(t, "first", actual.Items[0].Item.Title)
		assert.Equal(t, items[1].ID, actual.Items[1].Item.ID)
		assert.Equal(t, "second", actual.Items[1].Item.Title)
	})
}
