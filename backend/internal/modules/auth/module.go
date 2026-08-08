package auth

import (
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/application"
	authjwt "github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/infrastructure/jwt"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/ports"
	httpapi "github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/presentation/http"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type Module struct {
	loginUsecase ports.LoginUsecase
}

// Config — подмножество shared/config.JWTConfig, нужное этому модулю. Отдельный
// тип вместо прямой зависимости на shared/config, чтобы не тянуть весь конфиг
// приложения ради четырёх полей.
type Config struct {
	Issuer   string
	Audience string
	KeyID    string
	TTL      time.Duration
}

func New(pool tx.DBTX, privateKeyPath string, cfg Config, logger ports.Logger) (*Module, error) {
	privateKey, err := authjwt.LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	issuer := authjwt.NewIssuer(
		privateKey,
		cfg.Issuer,
		cfg.Audience,
		cfg.KeyID,
		cfg.TTL,
	)

	users := postgres.NewUsersRepository(pool)

	return &Module{
		loginUsecase: application.NewLoginUsecase(users, issuer, logger),
	}, nil
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/auth/login", httpapi.Login(m.loginUsecase))
}
