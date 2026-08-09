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

type GetMyTicketUsecase struct {
	advance ports.AdvanceUsecase
	repo    ports.QueueRepository
	logger  ports.Logger
	ttl     time.Duration
}

func NewGetMyTicketUsecase(
	advance ports.AdvanceUsecase,
	repo ports.QueueRepository,
	logger ports.Logger,
	ttl time.Duration,
) *GetMyTicketUsecase {
	return &GetMyTicketUsecase{
		advance: advance,
		repo:    repo,
		logger:  logger,
		ttl:     ttl,
	}
}

func (u *GetMyTicketUsecase) GetMyTicket(ctx context.Context, itemID, userID uuid.UUID) (*dto.Ticket, error) {
	if err := u.advance.AdvanceQueue(ctx, itemID, u.ttl); err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			// Контракт GET /queue/me отдельного 404 под "товар не найден" не
			// определяет — только TicketNotFound. Раз товара нет, заявки по
			// нему тем более нет.
			return nil, domain.ErrTicketNotFound
		}

		u.logger.Error(
			"failed to advance queue before reading ticket",
			"itemId", itemID,
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	ticket, err := u.repo.GetLatestTicket(ctx, itemID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrTicketNotFound) {
			return nil, domain.ErrTicketNotFound
		}

		u.logger.Error(
			"failed to get latest ticket",
			"itemId", itemID,
			"userId", userID,
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	now := time.Now().UTC()
	result := dto.NewTicket(ticket, now)

	if ticket.Status == domain.QueueStatusQueued {
		ahead, err := u.repo.CountQueuedAhead(ctx, itemID, ticket.CreatedAt)
		if err != nil {
			u.logger.Error("failed to count queued ahead", "error", err)
			return nil, domain.ErrInternal
		}

		nextFree, err := u.repo.NextSlotFreeAt(ctx, itemID)
		if err != nil {
			u.logger.Error("failed to get next slot free at", "error", err)
			return nil, domain.ErrInternal
		}

		result.SetQueuedFields(ahead, nextFree, now)
	}

	return result, nil
}
