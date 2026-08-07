package integrations

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	checkout "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout"
	checkoutapplication "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/application"
	checkoutdomain "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
	checkoutdto "github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
)

func newPaymentCallbackUsecase(ttl time.Duration) *checkoutapplication.PaymentCallbackUsecase {
	queueRepo := postgres.NewQueueRepository(db)
	advance := checkout.NewAdvanceAdapter(newAdvanceQueueUsecase())

	return checkoutapplication.NewPaymentCallbackUsecase(advance, queueRepo, slog.Default(), ttl)
}

func paidRequest(ticketID uuid.UUID) *checkoutdto.PaymentCallbackRequest {
	return &checkoutdto.PaymentCallbackRequest{TicketID: ticketID, Result: "paid"}
}

func failedRequest(ticketID uuid.UUID) *checkoutdto.PaymentCallbackRequest {
	return &checkoutdto.PaymentCallbackRequest{TicketID: ticketID, Result: "failed"}
}

func TestPaymentCallbackUsecase_Paid(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	userID := uuid.New()
	ticketID := uuid.New()
	expiresAt := time.Now().UTC().Add(30 * time.Second)
	seedTicket(t, ticketID, itemID, userID, domain.QueueStatusCheckout, time.Now().UTC(), &expiresAt)

	usecase := newPaymentCallbackUsecase(3 * time.Second)
	ticket, err := usecase.HandleCallback(context.Background(), itemID, userID, paidRequest(ticketID))
	require.NoError(t, err)

	assert.Equal(t, checkoutdomain.QueueStatusPurchased, ticket.Status)
	assert.Nil(t, ticket.ExpiresInSeconds)
}

func TestPaymentCallbackUsecase_Failed_DoesNotChangeStatusOrWindow(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	userID := uuid.New()
	ticketID := uuid.New()
	expiresAt := time.Now().UTC().Add(30 * time.Second)
	seedTicket(t, ticketID, itemID, userID, domain.QueueStatusCheckout, time.Now().UTC(), &expiresAt)

	usecase := newPaymentCallbackUsecase(3 * time.Second)
	ticket, err := usecase.HandleCallback(context.Background(), itemID, userID, failedRequest(ticketID))
	require.NoError(t, err)

	assert.Equal(t, checkoutdomain.QueueStatusCheckout, ticket.Status)
	assert.WithinDuration(t, expiresAt, *ticket.ExpiresAt, time.Millisecond)
}

// TestPaymentCallbackUsecase_FailedRetriesDoNotExtendWindow — NFR1: 10 подряд failed внутри
// окна, статус остаётся CHECKOUT, expiresAt не меняется, слот не освобождается.
func TestPaymentCallbackUsecase_FailedRetriesDoNotExtendWindow(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	userID := uuid.New()
	ticketID := uuid.New()
	expiresAt := time.Now().UTC().Add(30 * time.Second)
	seedTicket(t, ticketID, itemID, userID, domain.QueueStatusCheckout, time.Now().UTC(), &expiresAt)

	usecase := newPaymentCallbackUsecase(3 * time.Second)
	for i := 0; i < 10; i++ {
		ticket, err := usecase.HandleCallback(context.Background(), itemID, userID, failedRequest(ticketID))
		require.NoError(t, err)
		assert.Equal(t, checkoutdomain.QueueStatusCheckout, ticket.Status)
		assert.WithinDuration(t, expiresAt, *ticket.ExpiresAt, time.Millisecond, "attempt %d must not move expiresAt", i)
	}

	assert.Equal(t, 1, countTicketsByStatus(t, itemID, domain.QueueStatusCheckout))
}

func TestPaymentCallbackUsecase_TooLate(t *testing.T) {
	itemID := uuid.New()
	usecase := newPaymentCallbackUsecase(3 * time.Second)

	t.Run("paid after expiresAt", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		userID := uuid.New()
		ticketID := uuid.New()
		expiredAt := time.Now().UTC().Add(-time.Second)
		seedTicket(t, ticketID, itemID, userID, domain.QueueStatusCheckout, time.Now().UTC().Add(-time.Minute), &expiredAt)

		_, err := usecase.HandleCallback(context.Background(), itemID, userID, paidRequest(ticketID))
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrTooLate)
	})

	t.Run("failed after expiresAt — symmetric with paid", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		userID := uuid.New()
		ticketID := uuid.New()
		expiredAt := time.Now().UTC().Add(-time.Second)
		seedTicket(t, ticketID, itemID, userID, domain.QueueStatusCheckout, time.Now().UTC().Add(-time.Minute), &expiredAt)

		_, err := usecase.HandleCallback(context.Background(), itemID, userID, failedRequest(ticketID))
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrTooLate)
	})

	t.Run("paid retried after already purchased", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		userID := uuid.New()
		ticketID := uuid.New()
		seedTicket(t, ticketID, itemID, userID, domain.QueueStatusPurchased, time.Now().UTC(), nil)

		_, err := usecase.HandleCallback(context.Background(), itemID, userID, paidRequest(ticketID))
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrTooLate)
	})

	t.Run("callback on ticket still in OFFERED, never reached CHECKOUT", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		userID := uuid.New()
		ticketID := uuid.New()
		expiresAt := time.Now().UTC().Add(30 * time.Second)
		seedTicket(t, ticketID, itemID, userID, domain.QueueStatusOffered, time.Now().UTC(), &expiresAt)

		_, err := usecase.HandleCallback(context.Background(), itemID, userID, paidRequest(ticketID))
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrTooLate)
	})

	// TestPaymentCallbackUsecase_TooLate/callback_about_an_old_expired_ticket_must_not_touch_a_newer_active_one
	// — регрессия: T1 давно EXPIRED, T2 (та же пара item+user) сейчас в CHECKOUT. Callback,
	// адресованный T1 (по её ticketId), не должен ни воскресить T1, ни задеть T2 — раньше
	// матчинг шёл только по (item, user, status=CHECKOUT), и опоздавший сигнал про T1
	// ошибочно применялся бы к T2.
	t.Run("callback about an old expired ticket must not touch a newer active one", func(t *testing.T) {
		truncate(db, t)
		seedItem(t, itemID, 1)

		userID := uuid.New()
		oldExpiredAt := time.Now().UTC().Add(-time.Hour)
		oldTicketID := uuid.New()
		seedTicket(t, oldTicketID, itemID, userID, domain.QueueStatusExpired, time.Now().UTC().Add(-2*time.Hour), &oldExpiredAt)

		newTicketID := uuid.New()
		newExpiresAt := time.Now().UTC().Add(30 * time.Second)
		seedTicket(t, newTicketID, itemID, userID, domain.QueueStatusCheckout, time.Now().UTC(), &newExpiresAt)

		_, err := usecase.HandleCallback(context.Background(), itemID, userID, paidRequest(oldTicketID))
		require.Error(t, err)
		assert.ErrorIs(t, err, checkoutdomain.ErrTooLate, "old ticket exists but isn't CHECKOUT — same bucket as any other non-CHECKOUT ticket")

		newTicket, err := postgres.NewQueueRepository(db).GetLatestTicket(context.Background(), itemID, userID)
		require.NoError(t, err)
		assert.Equal(t, newTicketID, newTicket.ID, "callback about the old ticket must not resolve to the new one")
		assert.Equal(t, domain.QueueStatusCheckout, newTicket.Status, "new ticket must stay untouched")
	})
}

func TestPaymentCallbackUsecase_TicketNotFound(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	usecase := newPaymentCallbackUsecase(3 * time.Second)
	_, err := usecase.HandleCallback(context.Background(), itemID, uuid.New(), paidRequest(uuid.New()))
	require.Error(t, err)
	assert.ErrorIs(t, err, checkoutdomain.ErrTicketNotFound)
}

// TestPaymentCallbackUsecase_ConcurrentPaid_ExactlyOneWins — атомарность FinalizeCheckoutResult:
// число задетых строк условного UPDATE и есть результат проверки (см. checkout-plan.md).
// Из N параллельных paid-подтверждений одного и того же тикета ровно одно должно применить
// переход в PURCHASED, остальные обязаны увидеть, что тикет уже не в CHECKOUT.
func TestPaymentCallbackUsecase_ConcurrentPaid_ExactlyOneWins(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1)

	userID := uuid.New()
	ticketID := uuid.New()
	expiresAt := time.Now().UTC().Add(30 * time.Second)
	seedTicket(t, ticketID, itemID, userID, domain.QueueStatusCheckout, time.Now().UTC(), &expiresAt)

	usecase := newPaymentCallbackUsecase(3 * time.Second)

	const attempts = 20
	var wg sync.WaitGroup
	errs := make(chan error, attempts)
	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := usecase.HandleCallback(context.Background(), itemID, userID, paidRequest(ticketID))
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	successes, tooLate := 0, 0
	for err := range errs {
		if err == nil {
			successes++
			continue
		}

		assert.ErrorIs(t, err, checkoutdomain.ErrTooLate)
		tooLate++
	}

	assert.Equal(t, 1, successes)
	assert.Equal(t, attempts-1, tooLate)
	assert.Equal(t, 1, countTicketsByStatus(t, itemID, domain.QueueStatusPurchased))
}
