package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
)

type DBTX interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type QueueRepository struct {
	pool DBTX
}

func NewQueueRepository(pool DBTX) *QueueRepository {
	return &QueueRepository{
		pool: pool,
	}
}

func (repo *QueueRepository) CreateQueue(ctx context.Context, queue *domain.Queue) (*domain.Queue, error) {
	const query string = `INSERT INTO 
    queues (id, product_id, user_id)
		VALUES ($1, $2, $3) RETURNING 
		id, product_id, user_id, status, 
		created_at, updated_at, expires_at`

	var created domain.Queue
	err := repo.pool.QueryRow(
		ctx,
		query,
		queue.ID,
		queue.ProductID,
		queue.UserID,
	).Scan(
		&created.ID,
		&created.ProductID,
		&created.UserID,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.ExpiresAt,
	)
	if err != nil {
		if checkUniqueViolation(err) {
			return nil, domain.ErrUserAlreadyQueued
		}

		return nil, fmt.Errorf("queue user: %w", err)
	}

	return &created, nil
}

func (repo *QueueRepository) TryPromoteUser(
	ctx context.Context,
	queueID, productID uuid.UUID,
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
			WHERE product_id = $2
			AND status IN ('OFFERED', 'CHECKOUT')
		)
		AND NOT EXISTS(
			SELECT 1 FROM queues
			WHERE product_id = $2
			AND status = 'SOLD_OUT'::queue_status
		)
		AND NOT EXISTS(
			SELECT 1 FROM queues AS q
			WHERE q.product_id = $2
			AND q.status = 'QUEUED'::queue_status
			AND q.created_at < queues.created_at
		)
		
		RETURNING expires_at`

	var expiresAt time.Time
	ttlStr := fmt.Sprintf(
		"%d seconds",
		int(ttl.Seconds()),
	)

	err := repo.pool.QueryRow(
		ctx,
		query,
		queueID,
		productID,
		ttlStr,
	).Scan(&expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, nil
		}

		if checkUniqueViolation(err) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf("promote user: %w", err)
	}

	return true, &expiresAt, nil
}

func (repo *QueueRepository) QuitQueue(ctx context.Context, productID, userID uuid.UUID) error {
	const query string = `
		UPDATE queues
		SET status = 'CANCELLED'
		WHERE product_id = $1 AND user_id = $2
	`
	result, err := repo.pool.Exec(ctx, query, productID, userID)
	if err != nil {
		return fmt.Errorf("quit queue: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrQueueNotFound
	}
	return nil
}
