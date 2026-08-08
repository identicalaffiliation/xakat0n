package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	checkoutdomain "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	queuedomain "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type ItemsRepository struct {
	pool tx.DBTX
}

func NewItemsRepository(pool tx.DBTX) *ItemsRepository {
	return &ItemsRepository{
		pool: pool,
	}
}

func (repo *ItemsRepository) dbtx(ctx context.Context) tx.DBTX {
	return tx.DBTXFromContext(ctx, repo.pool)
}

func (repo *ItemsRepository) CreateItem(ctx context.Context, item *domain.Item) error {
	const query string = `
			INSERT INTO items
     	(title, description, price, is_limited)
    	VALUES ($1, $2, $3, $4)
		`

	if _, err := repo.pool.Exec(ctx, query, item.Title, item.Description, item.Price, item.IsLimited); err != nil {
		return fmt.Errorf("create item: %w", err)
	}

	return nil
}

func (repo *ItemsRepository) GetAll(ctx context.Context) ([]*domain.Item, error) {
	const query string = `SELECT * FROM items ORDER BY created_at DESC`

	var items []*domain.Item
	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get items: %w", err)
	}

	for rows.Next() {
		var item domain.Item
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Description,
			&item.Price,
			&item.Category,
			&item.IsLimited,
			&item.Stock,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan items: %w", err)
	}

	return items, nil
}

func (repo *ItemsRepository) GetItemByID(ctx context.Context, itemID uuid.UUID) (*domain.Item, error) {
	const query = `SELECT * FROM items WHERE id = $1`
	var item domain.Item
	err := repo.pool.QueryRow(ctx, query, itemID).Scan(
		&item.ID,
		&item.Title,
		&item.Description,
		&item.Price,
		&item.Category,
		&item.IsLimited,
		&item.Stock,
		&item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrItemNotFound
		}

		return nil, fmt.Errorf("get item: %w", err)
	}

	return &item, nil
}

// IsLimited — реализует checkout/ports.ItemsRepository, без FOR UPDATE в отличие от LockStock.
func (repo *ItemsRepository) IsLimited(ctx context.Context, itemID uuid.UUID) (bool, error) {
	const query = `SELECT is_limited FROM items WHERE id = $1`

	var isLimited bool
	err := repo.dbtx(ctx).QueryRow(ctx, query, itemID).Scan(&isLimited)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, checkoutdomain.ErrItemNotFound
		}

		return false, fmt.Errorf("is limited: %w", err)
	}

	return isLimited, nil
}

// LockStock — узкое чтение для queue-модуля: реализует queue/ports.ItemsRepository.
// Отдельно от GetItemByID, потому что берёт FOR UPDATE (нужно advanceQueue) и не тянет
// каталожные поля, которые queue-модулю не нужны.
func (repo *ItemsRepository) LockStock(ctx context.Context, itemID uuid.UUID) (*queuedomain.Item, error) {
	const query = `SELECT id, stock, is_limited FROM items WHERE id = $1 FOR UPDATE`

	var item queuedomain.Item
	err := repo.dbtx(ctx).QueryRow(ctx, query, itemID).Scan(&item.ID, &item.Stock, &item.IsLimited)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, queuedomain.ErrItemNotFound
		}

		return nil, fmt.Errorf("lock stock: %w", err)
	}

	return &item, nil
}
