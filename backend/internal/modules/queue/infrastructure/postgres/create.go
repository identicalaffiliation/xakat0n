package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	pgerrors "github.com/identicalaffiliation/xakat0n/backend/internal/shared/postgres"
)

func (repo *QueueRepository) CreateQueue(ctx context.Context, queue *domain.Queue) (*domain.Queue, error) {
	const query string = `INSERT INTO
	queues (id, item_id, user_id)
		VALUES ($1, $2, $3) RETURNING
		id, item_id, user_id, status,
		created_at, updated_at, expires_at`

	var created domain.Queue
	err := repo.dbtx(ctx).QueryRow(
		ctx,
		query,
		queue.ID,
		queue.ItemID,
		queue.UserID,
	).Scan(
		&created.ID,
		&created.ItemID,
		&created.UserID,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.ExpiresAt,
	)
	if err != nil {
		if pgerrors.IsUniqueViolation(err) {
			return nil, domain.ErrUserAlreadyQueued
		}

		return nil, fmt.Errorf("queue user: %w", err)
	}

	return &created, nil
}

func (repo *QueueRepository) GetActiveTicket(
	ctx context.Context,
	itemID, userID uuid.UUID,
) (*domain.Queue, error) {
	const query = `
		SELECT id, item_id, user_id, status, created_at, updated_at, expires_at
		FROM queues
		WHERE item_id = $1 AND user_id = $2
		AND status IN ('QUEUED', 'OFFERED', 'CHECKOUT')
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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get active ticket: %w", err)
	}

	return &q, nil
}

func (repo *QueueRepository) TryPromoteUser(
	ctx context.Context,
	queueID, itemID uuid.UUID,
	ttl time.Duration,
) (bool, *time.Time, error) {
	const query string = `
		UPDATE queues
		SET status = 'OFFERED',
		updated_at=now(),
		expires_at=now() + $3::interval
		WHERE id = $1
		AND status = 'QUEUED'::queue_status
		AND NOT EXISTS (
			SELECT 1 FROM queues
			WHERE item_id = $2
			AND status IN ('OFFERED', 'CHECKOUT')
		)
		AND NOT EXISTS(
			SELECT 1 FROM queues
			WHERE item_id = $2
			AND status = 'SOLD_OUT'::queue_status
		)
		AND NOT EXISTS(
			SELECT 1 FROM queues AS q
			WHERE q.item_id = $2
			AND q.status = 'QUEUED'::queue_status
			AND q.created_at < queues.created_at
		)

		RETURNING expires_at`

	var expiresAt time.Time
	ttlStr := fmt.Sprintf(
		"%d seconds",
		int(ttl.Seconds()),
	)

	err := repo.dbtx(ctx).QueryRow(
		ctx,
		query,
		queueID,
		itemID,
		ttlStr,
	).Scan(&expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, nil
		}

		if pgerrors.IsUniqueViolation(err) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf("promote user: %w", err)
	}

	return true, &expiresAt, nil
}
