package items

import (
	"github.com/go-chi/chi/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
	httpapi "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/presentation/http"
	queuepostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type Module struct {
	getItemUsecase  ports.GetItemUsecase
	getItemsUsecase ports.GetAllItemsUsecase
}

func New(pool tx.DBTX, logger ports.Logger) *Module {
	repo := postgres.NewItemsRepository(pool)
	soldOutChecker := queuepostgres.NewQueueRepository(pool)

	return &Module{
		getItemsUsecase: application.NewGetAllItemsUsecase(repo, soldOutChecker, logger),
		getItemUsecase:  application.NewGetItemUsecase(repo, soldOutChecker, logger),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/items", httpapi.GetItems(m.getItemsUsecase))
	r.Get("/api/v1/items/{itemId}", httpapi.GetItem(m.getItemUsecase))
}
