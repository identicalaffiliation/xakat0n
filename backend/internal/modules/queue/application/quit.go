package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
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
			u.logger.Warn(
				"queue not found",
				"error", err,
				"userId", userID,
				"itemId", itemID,
			)
			return nil, err
		}
		u.logger.Error(
			"failed to Quit from queue",
			"error", err,
			"userId", userID,
			"itemId", itemID,
		)
		return nil, domain.ErrInternal
	}

	err = u.advance.AdvanceQueue(ctx, itemID, u.ttl)
	if err != nil {
		u.logger.Error(
			"error advance queue",
			"itemId", itemID,
			"userId", userID,
		)
	}

	return dto.NewTicket(queue, time.Now().UTC()), nil
}
