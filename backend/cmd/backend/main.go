package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/identicalaffiliation/xakat0n/backend/internal/adapters/database"
	"github.com/identicalaffiliation/xakat0n/backend/internal/config"
	"github.com/identicalaffiliation/xakat0n/backend/internal/usecase"
	"github.com/identicalaffiliation/xakat0n/backend/pkg/httpserver"
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
	repo := database.NewQueueRepository(pool)

	createUsecase := usecase.NewCreateQueueUsecase(
		repo,
		txManager,
		slogger,
		time.Second*3,
	)
	quitUsecase := usecase.NewQuitQueueUsecase(
		repo,
		txManager,
		slogger,
		time.Second*3,
	)

	server := httpserver.SetupServer(&cfg.ServerConfig, createUsecase, quitUsecase)
	notifyChan := make(chan os.Signal, 1)
	signal.Notify(notifyChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogger.Error("failed to listen server", "error", err)
		}
	}()

	<-notifyChan

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slogger.Error("failed to shutdown server", "error", err)
	}

	slogger.Debug("server stopped gracefully..")
}
