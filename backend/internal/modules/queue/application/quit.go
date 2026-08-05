package application

import (
	"context"
	"errors"
	"time"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
)

type QuitQueueUsecase struct {
	repo      ports.QueueRepository
	txManager ports.TxManager
	logger    ports.Logger
	ttl       time.Duration
}

func NewQuitQueueUsecase(
	repo ports.QueueRepository,
	manager ports.TxManager,
	logger ports.Logger,
	ttl time.Duration,
) *QuitQueueUsecase {
	return &QuitQueueUsecase{
		repo:      repo,
		txManager: manager,
		logger:    logger,
		ttl:       ttl,
	}
}

func (u *QuitQueueUsecase) QuitQueue(ctx context.Context, in *dto.QuitQueueRequest) (*dto.QueueResponse, error) {
	var queue *domain.Queue
	err := u.txManager.WithTx(ctx, func(ctx context.Context) error {
		result, err := u.repo.QuitQueue(ctx, in.ProductID, in.UserID, u.ttl)
		if err != nil {
			return err
		}

		queue = result
		return nil
	})
	if err != nil {
		if errors.Is(err, domain.ErrQueueNotFound) {
			u.logger.Warn(
				"queue not found",
				"error", err,
				"userID", in.UserID,
				"product_id", in.ProductID,
			)
			return nil, err
		}
		u.logger.Error(
			"failed to Quit from queue",
			"error", err,
			"userID", in.UserID,
			"product_id", in.ProductID,
		)
		return nil, err
	}
	return dto.NewQueueResponse(queue), nil
}
