package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
)

type ItemsResponse struct {
	Items []ItemResponse `json:"items"`
}

type ItemResponse struct {
	Item Item `json:"item"`
}

type Item struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Price       int64     `json:"price"`
	Stock       int       `json:"stock"`
	IsLimited   bool      `json:"isLimited"`
	Category    *string   `json:"category"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewItemResponse(d *domain.Item) *ItemResponse {
	return &ItemResponse{
		Item: Item{
			ID:          d.ID,
			Title:       d.Title,
			Description: &d.Description,
			Price:       d.Price,
			Stock:       d.Stock,
			Category:    d.Category,
			IsLimited:   d.IsLimited,
			CreatedAt:   d.CreatedAt,
		},
	}
}

func NewItemsResponse(dms []*domain.Item) *ItemsResponse {
	var items []ItemResponse
	for _, d := range dms {
		items = append(items, ItemResponse{
			Item: Item{
				ID:          d.ID,
				Title:       d.Title,
				Description: &d.Description,
				Price:       d.Price,
				Stock:       d.Stock,
				Category:    d.Category,
				IsLimited:   d.IsLimited,
				CreatedAt:   d.CreatedAt,
			},
		})
	}

	return &ItemsResponse{
		Items: items,
	}
}
