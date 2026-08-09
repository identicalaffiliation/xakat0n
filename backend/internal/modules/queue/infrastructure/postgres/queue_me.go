package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

func (repo *QueueRepository) GetLatestTicket(ctx context.Context, productID, userID uuid.UUID) (*domain.Queue, error) {
	const query = `
		SELECT id, product_id, user_id, status, created_at, updated_at, expires_at
		FROM queues
		WHERE product_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT 1`

	var q domain.Queue
	err := repo.dbtx(ctx).QueryRow(ctx, query, productID, userID).Scan(
		&q.ID,
		&q.ProductID,
		&q.UserID,
		&q.Status,
		&q.CreatedAt,
		&q.UpdatedAt,
		&q.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.wrapError(ctx, domain.ErrTicketNotFound)
		}

		return nil, repo.wrapError(ctx, fmt.Errorf("get latest ticket: %w", err))
	}

	return &q, nil
}

func (repo *QueueRepository) CountQueuedAhead(ctx context.Context, productID uuid.UUID, createdAt time.Time) (int, error) {
	const query = `
		SELECT COUNT(*) FROM queues
		WHERE product_id = $1 AND status = 'QUEUED'::queue_status AND created_at < $2`

	var count int
	if err := repo.dbtx(ctx).QueryRow(ctx, query, productID, createdAt).Scan(&count); err != nil {
		return 0, repo.wrapError(ctx, fmt.Errorf("count queued ahead: %w", err))
	}

	return count, nil
}

func (repo *QueueRepository) NextSlotFreeAt(ctx context.Context, productID uuid.UUID) (*time.Time, error) {
	const query = `
		SELECT MIN(expires_at) FROM queues
		WHERE product_id = $1 AND status IN ('OFFERED', 'CHECKOUT')`

	var nextFree *time.Time
	if err := repo.dbtx(ctx).QueryRow(ctx, query, productID).Scan(&nextFree); err != nil {
		return nil, repo.wrapError(ctx, fmt.Errorf("next slot free at: %w", err))
	}

	return nextFree, nil
}
