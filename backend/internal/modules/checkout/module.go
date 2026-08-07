package checkout

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/ports"
	httpapi "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/presentation/http"
	itemspostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	queueapplication "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/application"
	queuepostgres "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
	queueports "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

type Module struct {
	checkoutUsecase ports.CheckoutUsecase
	paymentUsecase  ports.PaymentCallbackUsecase
}

// txManager нужен только для сборки внутреннего queue/application.AdvanceQueueUsecase.
func New(pool tx.DBTX, txManager queueports.TxManager, logger ports.Logger, ttl time.Duration) *Module {
	queueRepo := queuepostgres.NewQueueRepository(pool)
	itemsRepo := itemspostgres.NewItemsRepository(pool)

	advanceUsecase := queueapplication.NewAdvanceQueueUsecase(itemsRepo, queueRepo, txManager, logger)
	advance := NewAdvanceAdapter(advanceUsecase)

	return &Module{
		checkoutUsecase: application.NewCheckoutUsecase(advance, itemsRepo, queueRepo, logger, ttl),
		paymentUsecase:  application.NewPaymentCallbackUsecase(advance, queueRepo, logger, ttl),
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	checkoutRoute := fmt.Sprintf("/api/v1/items/{%s}/checkout", httpapi.ItemIdMuxPattern)
	r.With(httpx.SessionAuth).Post(checkoutRoute, httpapi.StartCheckout(m.checkoutUsecase))

	callbackRoute := fmt.Sprintf("/api/v1/items/{%s}/payment/callback", httpapi.ItemIdMuxPattern)
	r.With(httpx.SessionAuth).Post(callbackRoute, httpapi.PaymentCallback(m.paymentUsecase))
}
