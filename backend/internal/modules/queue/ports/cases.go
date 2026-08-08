package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
)

type CreateUsecase interface {
	CreateQueue(ctx context.Context, in *dto.CreateQueueRequest) (*dto.Ticket, error)
}

// AdvanceUsecase — переиспользуемый контракт продвижения очереди товара.
// ttl передаётся явным параметром вызова, а не полем конструктора: так его
// сможет переиспользовать любой другой вызывающий (например, POST /queue)
// со своим собственным TTL, не наследуя чужую конфигурацию.
type AdvanceUsecase interface {
	AdvanceQueue(ctx context.Context, itemID uuid.UUID, ttl time.Duration) error
}

type GetMeUsecase interface {
	GetMyTicket(ctx context.Context, itemID, userID uuid.UUID) (*dto.Ticket, error)
}
