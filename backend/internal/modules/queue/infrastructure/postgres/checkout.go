package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	checkoutdomain "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
)

// TryStartCheckout — реализует checkout/ports.QueueRepository; число задетых строк и есть результат проверки права.
func (repo *QueueRepository) TryStartCheckout(ctx context.Context, itemID, userID uuid.UUID) (*checkoutdomain.Queue, bool, error) {
	const query = `
		UPDATE queues SET status = 'CHECKOUT'::queue_status, updated_at = now()
		WHERE product_id = $1 AND user_id = $2
		AND status = 'OFFERED'::queue_status AND expires_at > now()
		RETURNING id, product_id, user_id, status, created_at, updated_at, expires_at`

	var q checkoutdomain.Queue
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
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("try start checkout: %w", err)
	}

	return &q, true, nil
}

// FinalizeCheckoutResult — один метод на paid и failed, обе ветки проверяют одно и то же условие.
// Матчит конкретный тикет по id, а не "текущий CHECKOUT для пары (item, user)": пользователь
// может успеть создать новую попытку T2, пока старая T1 уже EXPIRED, но ещё существует
// в базе — без id опоздавший сигнал про T1 применился бы к T2 (см. checkout-plan.md).
func (repo *QueueRepository) FinalizeCheckoutResult(ctx context.Context, itemID, userID, ticketID uuid.UUID, paid bool) (*checkoutdomain.Queue, bool, error) {
	const query = `
		UPDATE queues
		SET status = CASE WHEN $4 THEN 'PURCHASED'::queue_status ELSE status END,
		    updated_at = now()
		WHERE id = $3 AND product_id = $1 AND user_id = $2
		AND status = 'CHECKOUT'::queue_status AND expires_at > now()
		RETURNING id, product_id, user_id, status, created_at, updated_at, expires_at`

	var q checkoutdomain.Queue
	err := repo.dbtx(ctx).QueryRow(ctx, query, itemID, userID, ticketID, paid).Scan(
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
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("finalize checkout result: %w", err)
	}

	return &q, true, nil
}

// FindTicket — отдельный метод от GetLatestTicket, возвращает домен потребителя (checkout).
// Матчит конкретный тикет по id (не "последний по created_at" для пары) — та же причина,
// что и у FinalizeCheckoutResult: нужно диагностировать ошибку именно про тот тикет,
// про который спросили, а не про случайно оказавшийся текущим. product_id/user_id
// в WHERE остаются — без них чужой ticketId позволил бы получить чужой тикет.
func (repo *QueueRepository) FindTicket(ctx context.Context, itemID, userID, ticketID uuid.UUID) (*checkoutdomain.Queue, error) {
	const query = `
		SELECT id, product_id, user_id, status, created_at, updated_at, expires_at
		FROM queues
		WHERE id = $3 AND product_id = $1 AND user_id = $2`

	var q checkoutdomain.Queue
	err := repo.dbtx(ctx).QueryRow(ctx, query, itemID, userID, ticketID).Scan(
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
			return nil, checkoutdomain.ErrTicketNotFound
		}

		return nil, fmt.Errorf("find ticket: %w", err)
	}

	return &q, nil
}
