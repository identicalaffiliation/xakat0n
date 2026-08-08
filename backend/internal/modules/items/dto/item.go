package dto

import (
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
)

type Item struct {
	ItemID      uuid.UUID `json:"itemId"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Price       int64     `json:"price"`
	Category    *string   `json:"category"`
	IsLimited   bool      `json:"isLimited"`
	Stock       *int      `json:"stock"`
	SoldOut     bool      `json:"soldOut"`
}

func NewItem(d *domain.Item, soldOut bool) Item {
	var stock *int
	if d.IsLimited {
		stock = &d.Stock
	}

	return Item{
		ItemID:      d.ID,
		Title:       d.Title,
		Description: &d.Description,
		Price:       d.Price,
		Category:    d.Category,
		IsLimited:   d.IsLimited,
		Stock:       stock,
		SoldOut:     soldOut,
	}
}

func NewItems(dms []*domain.Item, soldOut map[uuid.UUID]bool) []Item {
	items := make([]Item, 0, len(dms))
	for _, d := range dms {
		items = append(items, NewItem(d, soldOut[d.ID]))
	}

	return items
}
