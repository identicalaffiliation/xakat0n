package application

import (
	"context"
	"errors"
	"time"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
)

type CreateQueueUsecase struct {
	repo      ports.QueueRepository
	txManager ports.TxManager
	logger    ports.Logger
	ttl       time.Duration
}

func NewCreateQueueUsecase(
	repo ports.QueueRepository,
	manager ports.TxManager,
	logger ports.Logger,
	ttl time.Duration,
) *CreateQueueUsecase {
	return &CreateQueueUsecase{
		repo:      repo,
		txManager: manager,
		logger:    logger,
		ttl:       ttl,
	}
}

func (u *CreateQueueUsecase) CreateQueue(
	ctx context.Context,
	in *dto.CreateQueueRequest,
) (*dto.CreateQueueResponse, error) {
	var created domain.Queue
	err := u.txManager.WithTx(ctx, func(ctx context.Context) error {
		queue, err := u.repo.CreateQueue(
			ctx,
			domain.NewQueue(in.ProductID, in.UserID),
		)
		if err != nil {
			u.logger.Error(
				"failed to create queue",
				"error", err,
			)
			return err
		}

		promoted, expiredAt, err := u.repo.TryPromoteUser(
			ctx,
			queue.ID,
			queue.ProductID,
			u.ttl,
		)
		if err != nil {
			return err
		}

		if promoted {
			queue.ExpiresAt = expiredAt
			queue.Status = domain.QueueStatusOffered
		}

		created = *queue
		return nil
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyQueued) {
			return nil, domain.ErrUserAlreadyQueued
		}

		u.logger.Error(
			"failed to put user in queue",
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	return dto.NewCreateResponse(&created), nil
}
