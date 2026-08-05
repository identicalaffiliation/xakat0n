package queue

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
	httpapi "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/presentation/http"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type Module struct {
	createUsecase ports.CreateUsecase
	quitUsecase   ports.QuitUsecase
}

func New(pool tx.DBTX, txManager ports.TxManager, logger ports.Logger, ttl time.Duration) *Module {
	repo := postgres.NewQueueRepository(pool)

	return &Module{
		createUsecase: application.NewCreateQueueUsecase(repo, txManager, logger, ttl),
		quitUsecase:   application.NewQuitQueueUsecase(repo, txManager, logger, ttl),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	route := fmt.Sprintf("/api/v1/products/{%s}/queue", httpapi.ProductIdMuxPattern)
	r.With(httpx.SessionAuth).Post(route, httpapi.PutUserInQueue(m.createUsecase))
	r.With(httpx.SessionAuth).Delete(route+"/me", httpapi.QuitQueue(m.quitUsecase))
}
