package auth

import (
	"time"

	"github.com/go-chi/chi/v5"

	authjwt "github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/infrastructure/jwt"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/ports"
	httpapi "github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/presentation/http"
)

type Module struct {
	issuer ports.TokenIssuer
}

func New(privateKeyPath string) (*Module, error) {
	privateKey, err := authjwt.LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	issuer := authjwt.NewIssuer(
		privateKey,
		"mock-auth",
		"xakat0n-api",
		"mock-key-1",
		time.Hour,
	)

	return &Module{issuer: issuer}, nil
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/auth/login", httpapi.Login(m.issuer))
}
