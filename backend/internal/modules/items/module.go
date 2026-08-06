package items

import (
	"github.com/go-chi/chi/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
	httpapi "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/presentation/http"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type Module struct {
	getItemUsecase  ports.GetItemUsecase
	getItemsUsecase ports.GetAllItemsUsecase
}

func New(pool tx.DBTX, logger ports.Logger) *Module {
	repo := postgres.NewItemsRepository(pool)

	return &Module{
		getItemsUsecase: application.NewGetAllItemsUsecase(repo, logger),
		getItemUsecase:  application.NewGetItemUsecase(repo, logger),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/products", httpapi.GetItems(m.getItemsUsecase))
	r.Get("/api/v1/products/{productId}", httpapi.GetItem(m.getItemUsecase))
}
