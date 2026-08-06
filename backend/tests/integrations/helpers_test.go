package integrations

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

// seedTicket вставляет строку queues напрямую SQL — минуя QueueRepository.CreateQueue,
// чтобы контролировать created_at/status/expires_at, которые репозиторий не позволяет
// задать явно (нужно для детерминированного порядка FIFO и для тестовых терминальных
// статусов вроде PURCHASED, которые в проде некому выставлять до появления checkout-модуля).
func seedTicket(
	t *testing.T,
	id, productID, userID uuid.UUID,
	status domain.QueueStatus,
	createdAt time.Time,
	expiresAt *time.Time,
) {
	t.Helper()

	_, err := db.Exec(
		context.Background(),
		`INSERT INTO queues (id, product_id, user_id, status, created_at, updated_at, expires_at)
		 VALUES ($1, $2, $3, $4::queue_status, $5, $5, $6)`,
		id,
		productID,
		userID,
		status,
		createdAt,
		expiresAt,
	)
	require.NoError(t, err)
}

func truncate(db *pgxpool.Pool, t *testing.T) {
	t.Helper()

	rows, err := db.Query(context.Background(), `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
			AND table_name != 'schema_migrations';
	`)
	require.NoError(t, err)
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		require.NoError(t, err)
		tables = append(tables, table)
	}
	require.NoError(t, rows.Err())
	if len(tables) == 0 {
		return
	}
	query := fmt.Sprintf(
		"TRUNCATE TABLE %s RESTART IDENTITY CASCADE;",
		strings.Join(tables, ", "))
	_, err = db.Exec(context.Background(), query)
	require.NoError(t, err)
}
