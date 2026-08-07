package queue

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"

	itemspostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
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
	getMeUsecase  ports.GetMeUsecase
}

func New(pool tx.DBTX, txManager ports.TxManager, logger ports.Logger, ttl time.Duration) *Module {
	repo := postgres.NewQueueRepository(pool)
	itemsRepo := itemspostgres.NewItemsRepository(pool)
	advanceUsecase := application.NewAdvanceQueueUsecase(itemsRepo, repo, txManager, logger)

	return &Module{
		createUsecase: application.NewCreateQueueUsecase(repo, txManager, logger, ttl),
		quitUsecase:   application.NewQuitQueueUsecase(repo, advanceUsecase, txManager, logger, ttl),
		getMeUsecase:  application.NewGetMyTicketUsecase(advanceUsecase, repo, logger, ttl),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	route := fmt.Sprintf("/api/v1/items/{%s}/queue", httpapi.ItemIdMuxPattern)
	r.With(httpx.SessionAuth).Post(route, httpapi.PutUserInQueue(m.createUsecase))
	r.With(httpx.SessionAuth).Delete(route+"/me", httpapi.QuitQueue(m.quitUsecase))

	createRoute := fmt.Sprintf(route+"/{%s}/queue", httpapi.ItemIdMuxPattern)
	r.With(httpx.SessionAuth).Post(createRoute, httpapi.PutUserInQueue(m.createUsecase))

	getMeRoute := fmt.Sprintf(route+"/{%s}/queue/me", httpapi.ItemIdMuxPattern)
	r.With(httpx.SessionAuth).Get(getMeRoute, httpapi.GetMyTicket(m.getMeUsecase))
}
