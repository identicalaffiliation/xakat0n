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

func (repo *QueueRepository) GetLatestTicket(ctx context.Context, itemID, userID uuid.UUID) (*domain.Queue, error) {
	const query = `
		SELECT id, item_id, user_id, status, created_at, updated_at, expires_at
		FROM queues
		WHERE item_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT 1`

	var q domain.Queue
	err := repo.dbtx(ctx).QueryRow(ctx, query, itemID, userID).Scan(
		&q.ID,
		&q.ItemID,
		&q.UserID,
		&q.Status,
		&q.CreatedAt,
		&q.UpdatedAt,
		&q.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTicketNotFound
		}

		return nil, fmt.Errorf("get latest ticket: %w", err)
	}

	return &q, nil
}

func (repo *QueueRepository) CountQueuedAhead(ctx context.Context, itemID uuid.UUID, createdAt time.Time) (int, error) {
	const query = `
		SELECT COUNT(*) FROM queues
		WHERE item_id = $1 AND status = 'QUEUED'::queue_status AND created_at < $2`

	var count int
	if err := repo.dbtx(ctx).QueryRow(ctx, query, itemID, createdAt).Scan(&count); err != nil {
		return 0, fmt.Errorf("count queued ahead: %w", err)
	}

	return count, nil
}

func (repo *QueueRepository) NextSlotFreeAt(ctx context.Context, itemID uuid.UUID) (*time.Time, error) {
	const query = `
		SELECT MIN(expires_at) FROM queues
		WHERE item_id = $1 AND status IN ('OFFERED', 'CHECKOUT')`

	var nextFree *time.Time
	if err := repo.dbtx(ctx).QueryRow(ctx, query, itemID).Scan(&nextFree); err != nil {
		return nil, fmt.Errorf("next slot free at: %w", err)
	}

	return nextFree, nil
}
