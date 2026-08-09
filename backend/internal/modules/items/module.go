package items

import (
	"github.com/go-chi/chi/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
	httpapi "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/presentation/http"
	queuepostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type Module struct {
	getItemUsecase  ports.GetItemUsecase
	getItemsUsecase ports.GetAllItemsUsecase
	logger          ports.Logger
}

func New(pool tx.DBTX, logger ports.Logger) *Module {
	repo := postgres.NewItemsRepository(pool, logger)
	soldOutChecker := queuepostgres.NewQueueRepository(pool, logger)

	return &Module{
		getItemsUsecase: application.NewGetAllItemsUsecase(repo, soldOutChecker, logger),
		getItemUsecase:  application.NewGetItemUsecase(repo, soldOutChecker, logger),
		logger:          logger,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.With(httpx.Metrics).Get("/api/v1/items", httpapi.GetItems(m.getItemsUsecase))
	r.With(httpx.Metrics).Get("/api/v1/items/{itemId}", httpapi.GetItem(m.logger, m.getItemUsecase))
}
