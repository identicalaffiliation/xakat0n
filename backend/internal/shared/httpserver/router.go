package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

// NewRouter собирает chi.Router с общими для всех модулей middleware. Модули
// регистрируют на нём свои роуты через RegisterRoutes, вызывается из cmd/api/main.go.
func NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)

	mux.Get("/metrics", httpx.MetricsHandler)

	return mux
}
