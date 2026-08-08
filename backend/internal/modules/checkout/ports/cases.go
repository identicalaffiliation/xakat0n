package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/dto"
)

type AdvanceUsecase interface {
	AdvanceQueue(ctx context.Context, itemID uuid.UUID, ttl time.Duration) error
}

type CheckoutUsecase interface {
	StartCheckout(ctx context.Context, itemID, userID uuid.UUID) (*dto.CheckoutStarted, error)
}

type PaymentCallbackUsecase interface {
	HandleCallback(ctx context.Context, itemID, userID uuid.UUID, in *dto.PaymentCallbackRequest) (*dto.Ticket, error)
}
