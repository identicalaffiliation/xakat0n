package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
	pgerrors "github.com/identicalaffiliation/xakat0n/backend/internal/shared/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type QueueRepository struct {
	pool   tx.DBTX
	logger ports.Logger
}

func NewQueueRepository(pool tx.DBTX, loggers ...ports.Logger) *QueueRepository {
	var logger ports.Logger
	if len(loggers) != 0 {
		logger = loggers[0]
	}

	return &QueueRepository{
		pool:   pool,
		logger: logger,
	}
}

func (repo *QueueRepository) wrapError(ctx context.Context, err error) error {
	if repo.logger == nil {
		return err
	}

	return repo.logger.WrapError(ctx, err)
}

func (repo *QueueRepository) dbtx(ctx context.Context) tx.DBTX {
	return tx.DBTXFromContext(ctx, repo.pool)
}

func (repo *QueueRepository) CreateQueue(ctx context.Context, queue *domain.Queue) (*domain.Queue, error) {
	const query string = `INSERT INTO
    queues (id, product_id, user_id)
		VALUES ($1, $2, $3) RETURNING
		id, product_id, user_id, status,
		created_at, updated_at, expires_at`

	var created domain.Queue
	err := repo.dbtx(ctx).QueryRow(
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
		if pgerrors.IsUniqueViolation(err) {
			return nil, repo.wrapError(ctx, domain.ErrUserAlreadyQueued)
		}

		return nil, repo.wrapError(ctx, fmt.Errorf("queue user: %w", err))
	}

	return &created, nil
}

func (repo *QueueRepository) GetActiveTicket(
	ctx context.Context,
	productID, userID uuid.UUID,
) (*domain.Queue, error) {
	const query = `
		SELECT id, product_id, user_id, status, created_at, updated_at, expires_at
		FROM queues
		WHERE product_id = $1 AND user_id = $2
		AND status IN ('QUEUED', 'OFFERED', 'CHECKOUT')
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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, repo.wrapError(ctx, fmt.Errorf("get active ticket: %w", err))
	}

	return &q, nil
}

func (repo *QueueRepository) IsSoldOut(ctx context.Context, productID uuid.UUID, stock int) (bool, error) {
	const query = `
		SELECT COUNT(*) FROM queues
		WHERE product_id = $1 AND status = 'PURCHASED'::queue_status`

	var purchased int
	if err := repo.dbtx(ctx).QueryRow(ctx, query, productID).Scan(&purchased); err != nil {
		return false, repo.wrapError(ctx, fmt.Errorf("is sold out: %w", err))
	}

	return purchased >= stock, nil
}

// CountPurchased возвращает количество PURCHASED-заявок по каждому из переданных
// товаров, сгруппированное одним запросом (используется items-модулем для soldOut
// в каталоге, где вызов IsSoldOut на каждый товар обернулся бы в N запросов).
// Товары без единой PURCHASED-заявки в карте отсутствуют.
func (repo *QueueRepository) CountPurchased(ctx context.Context, itemIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	counts := make(map[uuid.UUID]int, len(itemIDs))
	if len(itemIDs) == 0 {
		return counts, nil
	}

	const query = `
		SELECT product_id, COUNT(*) FROM queues
		WHERE product_id = ANY($1) AND status = 'PURCHASED'::queue_status
		GROUP BY product_id`

	rows, err := repo.dbtx(ctx).Query(ctx, query, itemIDs)
	if err != nil {
		return nil, repo.wrapError(ctx, fmt.Errorf("count purchased: %w", err))
	}
	defer rows.Close()

	for rows.Next() {
		var productID uuid.UUID
		var count int
		if err := rows.Scan(&productID, &count); err != nil {
			return nil, repo.wrapError(ctx, fmt.Errorf("scan count purchased: %w", err))
		}

		counts[productID] = count
	}

	if err := rows.Err(); err != nil {
		return nil, repo.wrapError(ctx, fmt.Errorf("count purchased: %w", err))
	}

	return counts, nil
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

	err := repo.dbtx(ctx).QueryRow(
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

		if pgerrors.IsUniqueViolation(err) {
			return false, nil, nil
		}

		return false, nil, repo.wrapError(ctx, fmt.Errorf("promote user: %w", err))
	}

	return true, &expiresAt, nil
}
