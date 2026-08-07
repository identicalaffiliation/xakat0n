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

	authModule "github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth"
	itemsModule "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/application"
	postgres2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	queueModule "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpserver"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/logger"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
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

	pool, cleanup, err := postgres.NewPool(context.Background(), &cfg.PostgresConfig)
	if err != nil {
		slogger.Error(
			"error", err,
		)
		os.Exit(1)
	}

	defer cleanup()

	txManager := tx.NewManager(pool, slogger)

	items := itemsModule.New(pool, slogger)
	queue := queueModule.New(pool, txManager, slogger, cfg.CheckoutTimer)

	application.AddSeedData(context.Background(), postgres2.NewItemsRepository(pool), slogger)
	auth, err := authModule.New(cfg.JWTConfig.PrivateKeyPath)
	if err != nil {
		slogger.Error(
			"error", err,
		)
		return
	}

	router := httpserver.NewRouter()
	queue.RegisterRoutes(router)
	items.RegisterRoutes(router)
	auth.RegisterRoutes(router)

	server := httpserver.New(&cfg.ServerConfig, router)
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
