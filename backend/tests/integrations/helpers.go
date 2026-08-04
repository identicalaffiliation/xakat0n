package integrations

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

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
