package domain

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID
	Title       string
	Description string
	Price       int64
	Category    *string
	IsLimited   bool
	Stock       int
	CreatedAt   time.Time
}

func NewItem(title, desscription string, price int64) *Item {
	return &Item{
		Title:       title,
		Description: desscription,
		Price:       price,
	}
}
