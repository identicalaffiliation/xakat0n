package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/identicalaffiliation/xakat0n/migrator/internal/config"
	"github.com/identicalaffiliation/xakat0n/migrator/pkg/logger"
)

const (
	dialect              = "postgres"
	driver               = "pgx"
	defaultMigrationsDir = "./migrations"
	cmdUp                = "up"
	cmdDown              = "down"
	cmdReset             = "reset"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		configPath string
		command    string
		dir        string
	)

	flag.StringVar(&configPath, "config", "", "path to config")
	flag.StringVar(&command, "command", "up", "migration command: up, down, reset")
	flag.StringVar(&dir, "dir", defaultMigrationsDir, "migrations directory")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	slogger, err := logger.NewLogger(cfg)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	db, err := sql.Open(driver, cfg.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			slogger.Error("close database", "error", err)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	switch command {
	case cmdUp:
		err = goose.UpContext(ctx, db, dir)
	case cmdDown:
		err = goose.DownContext(ctx, db, dir)
	case cmdReset:
		err = goose.ResetContext(ctx, db, dir)
	default:
		return fmt.Errorf("unknown command %q (available: up, down, reset)", command)
	}

	if err != nil {
		return fmt.Errorf("execute %q migrations: %w", command, err)
	}

	slogger.Debug("migration completed successfully", "command", command)
	return nil
}
