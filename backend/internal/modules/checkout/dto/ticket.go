package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
)

// Ticket — тот же JSON-контракт Ticket из api-contract.yaml, что и у queue/dto.Ticket.
type Ticket struct {
	TicketID              uuid.UUID          `json:"ticketId"`
	ItemID                uuid.UUID          `json:"itemId"`
	Status                domain.QueueStatus `json:"status"`
	Position              *int               `json:"position"`
	NextSlotFreeInSeconds *int64             `json:"nextSlotFreeInSeconds"`
	ExpiresInSeconds      *int64             `json:"expiresInSeconds"`
	ExpiresAt             *time.Time         `json:"expiresAt"`
	CreatedAt             time.Time          `json:"createdAt"`
	ServerTime            time.Time          `json:"serverTime"`
}

func NewTicket(queue *domain.Queue, now time.Time) *Ticket {
	t := &Ticket{
		TicketID:   queue.ID,
		ItemID:     queue.ItemID,
		Status:     queue.Status,
		ExpiresAt:  queue.ExpiresAt,
		CreatedAt:  queue.CreatedAt,
		ServerTime: now,
	}

	// expiresInSeconds заполняется только для CHECKOUT — PURCHASED уже не отсчитывает окно.
	if queue.Status == domain.QueueStatusCheckout && queue.ExpiresAt != nil {
		t.ExpiresInSeconds = clampSeconds(queue.ExpiresAt.Sub(now))
	}

	return t
}

// clampSeconds защищает отсчёт от расхождения часов клиента и сервера (NFR8).
func clampSeconds(d time.Duration) *int64 {
	seconds := int64(d.Seconds())
	if seconds < 0 {
		seconds = 0
	}

	return &seconds
}
