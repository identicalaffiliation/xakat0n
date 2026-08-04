package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/identicalaffiliation/xakat0n/backend/internal/adapters/controller"
	"github.com/identicalaffiliation/xakat0n/backend/internal/config"
	"github.com/identicalaffiliation/xakat0n/backend/internal/ports"
)

func SetupServer(cfg *config.ServerConfig, createUsecase ports.CreateUsecase, quitUsecase ports.QuitUsecase) *http.Server {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)

	mux.Post("/api/v1/items/{productId}/queue", controller.PutUserInQueue(createUsecase))
	mux.Delete("/api/v1/items/{productId}/queue/me", controller.QuitQueue(quitUsecase))

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
