package httpserver

import (
	"fmt"
	"net/http"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/config"
)

// New собирает http.Server поверх готового handler'а. Ничего не знает про
// конкретные роуты — их регистрируют модули через свой module.go, до вызова New.
func New(cfg *config.ServerConfig, handler http.Handler) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IddleTimeout,
	}
}
