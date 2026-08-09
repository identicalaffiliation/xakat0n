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

type CheckoutUsecase struct {
	advance ports.AdvanceUsecase
	items   ports.ItemsRepository
	queue   ports.QueueRepository
	logger  ports.Logger
	ttl     time.Duration
}

func NewCheckoutUsecase(
	advance ports.AdvanceUsecase,
	items ports.ItemsRepository,
	queue ports.QueueRepository,
	logger ports.Logger,
	ttl time.Duration,
) *CheckoutUsecase {
	return &CheckoutUsecase{
		advance: advance,
		items:   items,
		queue:   queue,
		logger:  logger,
		ttl:     ttl,
	}
}

func (u *CheckoutUsecase) StartCheckout(ctx context.Context, itemID, userID uuid.UUID) (*dto.CheckoutStarted, error) {
	if err := u.advance.AdvanceQueue(ctx, itemID, u.ttl); err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			return nil, domain.ErrItemNotFound
		}

		u.logger.Error("failed to advance queue before checkout", "itemId", itemID, "error", err)
		return nil, domain.ErrInternal
	}

	isLimited, err := u.items.IsLimited(ctx, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			return nil, domain.ErrItemNotFound
		}

		u.logger.Error("failed to check item limited flag", "itemId", itemID, "error", err)
		return nil, domain.ErrInternal
	}

	if !isLimited {
		return dto.NewCheckoutStarted(nil), nil
	}

	queue, ok, err := u.queue.TryStartCheckout(ctx, itemID, userID)
	if err != nil {
		u.logger.Error("failed to start checkout", "itemId", itemID, "userId", userID, "error", err)
		return nil, domain.ErrInternal
	}

	if !ok {
		return nil, domain.ErrNoActiveRight
	}

	return dto.NewCheckoutStarted(dto.NewTicket(queue, time.Now().UTC())), nil
}
