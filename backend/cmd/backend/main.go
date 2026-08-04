package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/identicalaffiliation/xakat0n/backend/internal/adapters/database"
	"github.com/identicalaffiliation/xakat0n/backend/internal/config"
	"github.com/identicalaffiliation/xakat0n/backend/pkg/logger"
	"github.com/identicalaffiliation/xakat0n/backend/pkg/psqlpool"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	slogger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}

	pool, cleanup, err := psqlpool.NewPool(context.Background(), &cfg.PostgresConfig)
	if err != nil {
		slogger.Error(
			"error", err,
		)
		os.Exit(1)
	}

	defer cleanup()

	txManager := database.NewTxManager(pool, slogger)
}
