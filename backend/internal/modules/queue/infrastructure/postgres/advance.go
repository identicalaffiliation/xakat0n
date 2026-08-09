package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo *QueueRepository) ExpireStale(ctx context.Context, itemID uuid.UUID) error {
	const query = `
		UPDATE queues SET status = 'EXPIRED'::queue_status, updated_at = now()
		WHERE item_id = $1
		AND status IN ('OFFERED', 'CHECKOUT')
		AND expires_at < now()`

	if _, err := repo.dbtx(ctx).Exec(ctx, query, itemID); err != nil {
		return fmt.Errorf("expire stale: %w", err)
	}

	return nil
}

func (repo *QueueRepository) CountTaken(ctx context.Context, itemID uuid.UUID) (int, error) {
	const query = `
		SELECT COUNT(*) FROM queues
		WHERE item_id = $1
		AND status IN ('OFFERED', 'CHECKOUT', 'PURCHASED')`

	var count int
	if err := repo.dbtx(ctx).QueryRow(ctx, query, itemID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count taken: %w", err)
	}

	return count, nil
}

func (repo *QueueRepository) PromoteNext(
	ctx context.Context,
	itemID uuid.UUID,
	freeSlots int,
	ttl time.Duration,
) (int64, error) {
	const query = `
		UPDATE queues
		SET status = 'OFFERED'::queue_status, expires_at = now() + $3::interval, updated_at = now()
		WHERE id IN (
			SELECT id FROM queues
			WHERE item_id = $1 AND status = 'QUEUED'::queue_status
			ORDER BY created_at
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		)`

	ttlStr := fmt.Sprintf("%d seconds", int(ttl.Seconds()))

	tag, err := repo.dbtx(ctx).Exec(ctx, query, itemID, freeSlots, ttlStr)
	if err != nil {
		return 0, fmt.Errorf("promote next: %w", err)
	}

	return tag.RowsAffected(), nil
}

func (repo *QueueRepository) MarkSoldOut(ctx context.Context, itemID uuid.UUID, stock int) error {
	const query = `
		UPDATE queues SET status = 'SOLD_OUT'::queue_status, updated_at = now()
		WHERE item_id = $1
		AND status = 'QUEUED'::queue_status
		AND (
			SELECT COUNT(*) FROM queues
			WHERE item_id = $1 AND status = 'PURCHASED'::queue_status
		) >= $2`

	if _, err := repo.dbtx(ctx).Exec(ctx, query, itemID, stock); err != nil {
		return fmt.Errorf("mark sold out: %w", err)
	}

	return nil
}
