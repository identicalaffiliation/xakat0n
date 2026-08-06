package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

// Ticket — плоский DTO по схеме Ticket из api-contract.yaml (snake_case),
// намеренно не переиспользует Queue/CreateQueueResponse (camelCase, вложенность).
type Ticket struct {
	TicketID              uuid.UUID          `json:"ticket_id"`
	ItemID                uuid.UUID          `json:"item_id"`
	Status                domain.QueueStatus `json:"status"`
	Position              *int               `json:"position"`
	NextSlotFreeInSeconds *int64             `json:"next_slot_free_in_seconds"`
	ExpiresInSeconds      *int64             `json:"expires_in_seconds"`
	ExpiresAt             *time.Time         `json:"expires_at"`
	CreatedAt             time.Time          `json:"created_at"`
	ServerTime            time.Time          `json:"server_time"`
}

func NewTicket(queue *domain.Queue, now time.Time) *Ticket {
	t := &Ticket{
		TicketID:   queue.ID,
		ItemID:     queue.ProductID,
		Status:     queue.Status,
		ExpiresAt:  queue.ExpiresAt,
		CreatedAt:  queue.CreatedAt,
		ServerTime: now,
	}

	isActiveWindow := queue.Status == domain.QueueStatusOffered || queue.Status == domain.QueueStatusCheckout
	if isActiveWindow && queue.ExpiresAt != nil {
		t.ExpiresInSeconds = clampSeconds(queue.ExpiresAt.Sub(now))
	}

	return t
}

// SetQueuedFields заполняет поля, которые заполняются только при status = QUEUED.
func (t *Ticket) SetQueuedFields(ahead int, nextSlotFreeAt *time.Time, now time.Time) {
	position := ahead + 1
	t.Position = &position

	if nextSlotFreeAt != nil {
		t.NextSlotFreeInSeconds = clampSeconds(nextSlotFreeAt.Sub(now))
	}
}

// clampSeconds защищает отсчёт от расхождения часов клиента и сервера (NFR8):
// отрицательный остаток (уже истекло, но advanceQueue ещё не пробежался) не должен уйти в минус.
func clampSeconds(d time.Duration) *int64 {
	seconds := int64(d.Seconds())
	if seconds < 0 {
		seconds = 0
	}

	return &seconds
}
