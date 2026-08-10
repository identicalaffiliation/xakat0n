package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (repo *QueueRepository) IsSoldOut(ctx context.Context, itemID uuid.UUID, stock int) (bool, error) {
	const query = `
		SELECT COUNT(*) FROM queues
		WHERE item_id = $1 AND status = 'PURCHASED'::queue_status`

	var purchased int
	if err := repo.dbtx(ctx).QueryRow(ctx, query, itemID).Scan(&purchased); err != nil {
		return false, fmt.Errorf("is sold out: %w", err)
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
		SELECT item_id, COUNT(*) FROM queues
		WHERE item_id = ANY($1) AND status = 'PURCHASED'::queue_status
		GROUP BY item_id`

	rows, err := repo.dbtx(ctx).Query(ctx, query, itemIDs)
	if err != nil {
		return nil, fmt.Errorf("count purchased: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var itemID uuid.UUID
		var count int
		if err := rows.Scan(&itemID, &count); err != nil {
			return nil, fmt.Errorf("scan count purchased: %w", err)
		}

		counts[itemID] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("count purchased: %w", err)
	}

	return counts, nil
}
