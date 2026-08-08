package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

func (repo *QueueRepository) QuitQueue(ctx context.Context, productID, userID uuid.UUID) (*domain.Queue, error) {
	const quitUserQuery string = `
		UPDATE queues
		SET status = 'CANCELLED', updated_at = now()
		WHERE product_id = $1
		AND user_id = $2
		AND status IN ('QUEUED', 'OFFERED', 'CHECKOUT')
		RETURNING
			id,
			product_id,
			user_id,
			status,
			created_at,
			updated_at,
			expires_at
	`
	var queue domain.Queue
	err := repo.dbtx(ctx).QueryRow(ctx, quitUserQuery, productID, userID).Scan(
		&queue.ID,
		&queue.ProductID,
		&queue.UserID,
		&queue.Status,
		&queue.CreatedAt,
		&queue.UpdatedAt,
		&queue.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrQueueNotFound
		}
		return nil, fmt.Errorf("quit queue: %w", err)
	}
	return &queue, nil
}
