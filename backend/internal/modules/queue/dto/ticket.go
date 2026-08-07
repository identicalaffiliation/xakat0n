package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

// Ticket — плоский DTO по схеме Ticket из api-contract.yaml (camelCase),
// намеренно не переиспользует Queue/CreateQueueResponse (вложенность в "queue").
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
