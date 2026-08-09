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
	advance   ports.AdvanceUsecase
	repo      ports.QueueRepository
	items     ports.ItemsRepository
	txManager ports.TxManager
	logger    ports.Logger
	ttl       time.Duration
}

func NewCreateQueueUsecase(
	advance ports.AdvanceUsecase,
	repo ports.QueueRepository,
	items ports.ItemsRepository,
	manager ports.TxManager,
	logger ports.Logger,
	ttl time.Duration,
) *CreateQueueUsecase {
	return &CreateQueueUsecase{
		advance:   advance,
		repo:      repo,
		items:     items,
		txManager: manager,
		logger:    logger,
		ttl:       ttl,
	}
}

func (u *CreateQueueUsecase) CreateQueue(
	ctx context.Context,
	in *dto.CreateQueueRequest,
) (*dto.Ticket, error) {
	if err := u.advance.AdvanceQueue(ctx, in.ItemID, u.ttl); err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			return nil, domain.ErrItemNotFound
		}

		u.logger.Error(
			"failed to advance queue before create",
			"itemId", in.ItemID,
			"userId", in.UserID,
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	created, err := u.createQueueTicket(ctx, in)
	if err != nil {
		if !errors.Is(err, domain.ErrQueueNotApplicable) &&
			!errors.Is(err, domain.ErrItemSoldOut) &&
			!errors.Is(err, domain.ErrItemNotFound) {
			u.logger.Error(
				"failed to put user in queue",
				"itemId", in.ItemID,
				"userId", in.UserID,
				"error", err,
			)
		}

		return nil, err
	}

	return u.buildTicket(ctx, created)
}

func (u *CreateQueueUsecase) createQueueTicket(
	ctx context.Context,
	in *dto.CreateQueueRequest,
) (domain.Queue, error) {
	var created domain.Queue
	err := u.txManager.WithTx(ctx, func(ctx context.Context) error {
		item, err := u.items.LockStock(ctx, in.ItemID)
		if err != nil {
			return err
		}

		if !item.IsLimited {
			return domain.ErrQueueNotApplicable
		}

		soldOut, err := u.repo.IsSoldOut(ctx, in.ItemID, item.Stock)
		if err != nil {
			return err
		}
		if soldOut {
			return domain.ErrItemSoldOut
		}

		existing, err := u.repo.GetActiveTicket(ctx, in.ItemID, in.UserID)
		if err != nil {
			return err
		}
		if existing != nil {
			created = *existing
			return nil
		}

		queue, err := u.repo.CreateQueue(
			ctx,
			domain.NewQueue(in.ItemID, in.UserID),
		)
		if err != nil {
			if errors.Is(err, domain.ErrUserAlreadyQueued) {
				existing, getErr := u.repo.GetActiveTicket(
					ctx,
					in.ItemID,
					in.UserID,
				)
				if getErr != nil {
					return getErr
				}
				if existing != nil {
					created = *existing
					return nil
				}

				return domain.ErrItemSoldOut
			}

			return err
		}

		taken, err := u.repo.CountTaken(ctx, in.ItemID)
		if err != nil {
			return err
		}

		if taken < item.Stock {
			promoted, expiresAt, err := u.repo.TryPromoteUser(
				ctx,
				queue.ID,
				in.ItemID,
				u.ttl,
			)
			if err != nil {
				return err
			}

			if promoted {
				queue.Status = domain.QueueStatusOffered
				queue.ExpiresAt = expiresAt
			}
		}

		created = *queue
		return nil
	})

	return created, err
}

func (u *CreateQueueUsecase) buildTicket(ctx context.Context, created domain.Queue) (*dto.Ticket, error) {
	now := time.Now().UTC()
	ticket := dto.NewTicket(&created, now)

	if created.Status != domain.QueueStatusQueued {
		return ticket, nil
	}

	ahead, err := u.repo.CountQueuedAhead(ctx, created.ItemID, created.CreatedAt)
	if err != nil {
		u.logger.Error("failed to count queued ahead", "error", err)
		return nil, domain.ErrInternal
	}

	nextFree, err := u.repo.NextSlotFreeAt(ctx, created.ItemID)
	if err != nil {
		u.logger.Error("failed to get next slot free at", "error", err)
		return nil, domain.ErrInternal
	}

	ticket.SetQueuedFields(ahead, nextFree, now)
	return ticket, nil
}
