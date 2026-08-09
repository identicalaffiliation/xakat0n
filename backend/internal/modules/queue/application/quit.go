package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/logging"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
)

type QuitQueueUsecase struct {
	queue     ports.QueueRepository
	advance   ports.AdvanceUsecase
	txManager ports.TxManager
	logger    ports.Logger
	ttl       time.Duration
}

func NewQuitQueueUsecase(
	queue ports.QueueRepository,
	advance ports.AdvanceUsecase,
	manager ports.TxManager,
	logger ports.Logger,
	ttl time.Duration,
) *QuitQueueUsecase {
	return &QuitQueueUsecase{
		queue:     queue,
		advance:   advance,
		txManager: manager,
		logger:    logger,
		ttl:       ttl,
	}
}

func (u *QuitQueueUsecase) QuitQueue(ctx context.Context, itemID, userID uuid.UUID) (*dto.Ticket, error) {
	queue, err := u.queue.QuitQueue(ctx, itemID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrQueueNotFound) {
			return nil, u.logger.WrapError(ctx, err)
		}
		ctx = u.logger.ContextFromError(ctx, err)
		u.logger.ErrorContext(
			ctx,
			"failed to Quit from queue",
			"error", err,
		)
		return nil, u.logger.WrapError(ctx, domain.ErrInternal)
	}
	ctx = logging.WithQueueID(ctx, u.logger, queue.ID)

	err = u.advance.AdvanceQueue(ctx, itemID, u.ttl)
	if err != nil {
		u.logger.ErrorContext(
			u.logger.ContextFromError(ctx, err),
			"failed to advance queue after quit",
			"error", err,
		)
	}

	return dto.NewTicket(queue, time.Now().UTC()), nil
}
