package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo *QueueRepository) ExpireStale(ctx context.Context, productID uuid.UUID) error {
	const query = `
		UPDATE queues SET status = 'EXPIRED'::queue_status, updated_at = now()
		WHERE product_id = $1
		AND status IN ('OFFERED', 'CHECKOUT')
		AND expires_at < now()`

	if _, err := repo.dbtx(ctx).Exec(ctx, query, productID); err != nil {
		return repo.wrapError(ctx, fmt.Errorf("expire stale: %w", err))
	}

	return nil
}

func (repo *QueueRepository) CountTaken(ctx context.Context, productID uuid.UUID) (int, error) {
	const query = `
		SELECT COUNT(*) FROM queues
		WHERE product_id = $1
		AND status IN ('OFFERED', 'CHECKOUT', 'PURCHASED')`

	var count int
	if err := repo.dbtx(ctx).QueryRow(ctx, query, productID).Scan(&count); err != nil {
		return 0, repo.wrapError(ctx, fmt.Errorf("count taken: %w", err))
	}

	return count, nil
}

func (repo *QueueRepository) PromoteNext(
	ctx context.Context,
	productID uuid.UUID,
	freeSlots int,
	ttl time.Duration,
) (int64, error) {
	const query = `
		UPDATE queues
		SET status = 'OFFERED'::queue_status, expires_at = now() + $3::interval, updated_at = now()
		WHERE id IN (
			SELECT id FROM queues
			WHERE product_id = $1 AND status = 'QUEUED'::queue_status
			ORDER BY created_at
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		)`

	ttlStr := fmt.Sprintf("%d seconds", int(ttl.Seconds()))

	tag, err := repo.dbtx(ctx).Exec(ctx, query, productID, freeSlots, ttlStr)
	if err != nil {
		return 0, repo.wrapError(ctx, fmt.Errorf("promote next: %w", err))
	}

	return tag.RowsAffected(), nil
}

func (repo *QueueRepository) MarkSoldOut(ctx context.Context, productID uuid.UUID, stock int) error {
	const query = `
		UPDATE queues SET status = 'SOLD_OUT'::queue_status, updated_at = now()
		WHERE product_id = $1
		AND status = 'QUEUED'::queue_status
		AND (
			SELECT COUNT(*) FROM queues
			WHERE product_id = $1 AND status = 'PURCHASED'::queue_status
		) >= $2`

	if _, err := repo.dbtx(ctx).Exec(ctx, query, productID, stock); err != nil {
		return repo.wrapError(ctx, fmt.Errorf("mark sold out: %w", err))
	}

	return nil
}
