package integrations

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	PostgresDockerImage = "postgres:17"
	PostgresDb          = "testdb"
	PostgresUser        = "testuser"
	PostgresPass        = "testpassword"
	PostgresSsl         = "sslmode=disable"
	MigrationsDir       = "./../../../migrator/migrations"
)

var (
	db *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	var (
		cleanups []func()
		err      error
		cleanUp  func()
	)

	db, cleanUp, err = setupPostgres(ctx)
	if err != nil {
		log.Fatalf("setup postgres err: %v", err)
	}
	cleanups = append(cleanups, cleanUp)

	code := m.Run()

	for i := len(cleanups) - 1; i >= 0; i-- {
		cleanups[i]()
	}

	os.Exit(code)
}

func setupPostgres(ctx context.Context) (*pgxpool.Pool, func(), error) {
	container, err := postgres.Run(
		ctx,
		PostgresDockerImage,
		postgres.WithDatabase(PostgresDb),
		postgres.WithUsername(PostgresUser),
		postgres.WithPassword(PostgresPass),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, err
	}

	connStr, err := container.ConnectionString(ctx, PostgresSsl)
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, nil, err
	}

	if err := goose.UpContext(ctx, stdlib.OpenDBFromPool(pool), MigrationsDir); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			log.Printf("terminate container: %v", err)
		}
	}

	return pool, cleanup, nil
}
