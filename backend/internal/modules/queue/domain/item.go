package domain

import "github.com/google/uuid"

// Item — узкое представление товара, нужное только queue-модулю для advanceQueue.
type Item struct {
	ID        uuid.UUID
	Stock     int
	IsLimited bool
}
