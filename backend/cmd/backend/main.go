package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/identicalaffiliation/xakat0n/backend/internal/adapters/controller"
	"github.com/identicalaffiliation/xakat0n/backend/internal/adapters/database"
	"github.com/identicalaffiliation/xakat0n/backend/internal/config"
	"github.com/identicalaffiliation/xakat0n/backend/internal/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/usecase"
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

	server := setupServer(&cfg.ServerConfig, createUsecase)
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

func setupServer(cfg *config.ServerConfig, createUsecase ports.CreateUsecase) *http.Server {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)

	mux.Post("/api/v1/products/{productId}/queue", controller.PutUserInQueue(createUsecase))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IddleTimeout,
	}

	return server
}
