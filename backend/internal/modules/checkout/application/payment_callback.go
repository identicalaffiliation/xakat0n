package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/ports"
)

type PaymentCallbackUsecase struct {
	advance ports.AdvanceUsecase
	queue   ports.QueueRepository
	logger  ports.Logger
	ttl     time.Duration
}

func NewPaymentCallbackUsecase(
	advance ports.AdvanceUsecase,
	queue ports.QueueRepository,
	logger ports.Logger,
	ttl time.Duration,
) *PaymentCallbackUsecase {
	return &PaymentCallbackUsecase{
		advance: advance,
		queue:   queue,
		logger:  logger,
		ttl:     ttl,
	}
}

func (u *PaymentCallbackUsecase) HandleCallback(
	ctx context.Context,
	itemID, userID uuid.UUID,
	in *dto.PaymentCallbackRequest,
) (*dto.Ticket, error) {
	queue, ok, err := u.queue.FinalizeCheckoutResult(ctx, itemID, userID, in.TicketID, in.Paid())
	if err != nil {
		ctx = u.logger.ContextFromError(ctx, err)
		u.logger.ErrorContext(ctx,
			"failed to finalize checkout result",
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	if !ok {
		if _, err := u.queue.FindTicket(ctx, itemID, userID, in.TicketID); err != nil {
			if errors.Is(err, domain.ErrTicketNotFound) {
				return nil, domain.ErrTicketNotFound
			}

			ctx = u.logger.ContextFromError(ctx, err)
			u.logger.ErrorContext(ctx,
				"failed to find ticket for disambiguation",
				"error", err,
			)
			return nil, domain.ErrInternal
		}

		return nil, domain.ErrTooLate
	}

	// Best-effort: ошибка тут не должна превращать уже закоммиченный успешный исход в 500.
	if err := u.advance.AdvanceQueue(ctx, itemID, u.ttl); err != nil {
		ctx = u.logger.ContextFromError(ctx, err)
		u.logger.ErrorContext(ctx, "failed to advance queue after payment callback", "error", err)
	}

	return dto.NewTicket(queue, time.Now().UTC()), nil
}
