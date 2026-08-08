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
		u.logger.Error(
			"failed to finalize checkout result",
			"item_id", itemID,
			"user_id", userID,
			"ticket_id", in.TicketID,
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	if !ok {
		if _, err := u.queue.FindTicket(ctx, itemID, userID, in.TicketID); err != nil {
			if errors.Is(err, domain.ErrTicketNotFound) {
				return nil, domain.ErrTicketNotFound
			}

			u.logger.Error(
				"failed to find ticket for disambiguation",
				"item_id", itemID,
				"user_id", userID,
				"error", err,
			)
			return nil, domain.ErrInternal
		}

		return nil, domain.ErrTooLate
	}

	// Best-effort: ошибка тут не должна превращать уже закоммиченный успешный исход в 500.
	if err := u.advance.AdvanceQueue(ctx, itemID, u.ttl); err != nil {
		u.logger.Error("failed to advance queue after payment callback", "item_id", itemID, "error", err)
	}

	return dto.NewTicket(queue, time.Now().UTC()), nil
}
