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

	"github.com/go-chi/chi/v5"

	authModule "github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth"
	authjwt "github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/infrastructure/jwt"
	checkoutModule "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout"
	itemsModule "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/application"
	postgres2 "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	queueModule "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpserver"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/logger"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

func main() {
	ctx := context.Background()
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

	pool, cleanup, err := postgres.NewPool(ctx, &cfg.PostgresConfig)
	if err != nil {
		slogger.ErrorContext(ctx,
			"failed to connect to postgres",
			"error", err,
		)
		os.Exit(1)
	}

	defer cleanup()

	txManager := tx.NewManager(pool, slogger)

	items := itemsModule.New(pool, slogger)
	queue := queueModule.New(pool, txManager, slogger, cfg.CheckoutTimer)
	checkout := checkoutModule.New(pool, txManager, slogger, cfg.CheckoutTimer)

	application.AddSeedData(ctx, postgres2.NewItemsRepository(pool, slogger), slogger)
	auth, err := authModule.New(pool, cfg.JWTConfig.PrivateKeyPath, authModule.Config{
		Issuer:   cfg.JWTConfig.Issuer,
		Audience: cfg.JWTConfig.Audience,
		KeyID:    cfg.JWTConfig.KeyID,
		TTL:      cfg.JWTConfig.TTL,
	}, slogger)
	if err != nil {
		slogger.ErrorContext(ctx,
			"failed to initialize auth module",
			"error", err,
		)
		return
	}

	publicKey, err := authjwt.LoadPublicKey(cfg.JWTConfig.PublicKeyPath)
	if err != nil {
		slogger.ErrorContext(ctx, "failed to load JWT public key", "error", err)
		return
	}

	verifier := authjwt.NewVerifier(
		publicKey,
		cfg.JWTConfig.Issuer,
		cfg.JWTConfig.Audience,
		cfg.JWTConfig.KeyID,
	)

	router := httpserver.NewRouter()

	auth.RegisterRoutes(router)
	items.RegisterRoutes(router)

	router.Group(func(r chi.Router) {
		r.Use(httpx.JWTAuth(slogger, verifier))

		queue.RegisterRoutes(r)
		checkout.RegisterRoutes(r)
	})

	server := httpserver.New(&cfg.ServerConfig, router)
	notifyChan := make(chan os.Signal, 1)
	signal.Notify(notifyChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogger.ErrorContext(ctx, "failed to listen server", "error", err)
		}
	}()

	<-notifyChan

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slogger.ErrorContext(shutdownCtx, "failed to shutdown server", "error", err)
	}

	slogger.DebugContext(ctx, "server stopped gracefully")
}
