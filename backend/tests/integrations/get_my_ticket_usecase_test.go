package integrations

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/application"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/infrastructure/postgres"
)

func newGetMyTicketUsecase(ttl time.Duration) *application.GetMyTicketUsecase {
	queueRepo := postgres.NewQueueRepository(db)
	advanceUsecase := newAdvanceQueueUsecase()

	return application.NewGetMyTicketUsecase(advanceUsecase, queueRepo, slog.Default(), ttl)
}

func TestGetMyTicketUsecase_Queued(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1, true)

	now := time.Now().UTC()
	holderExpiresAt := now.Add(30 * time.Second)
	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusOffered, now, &holderExpiresAt)

	firstInLine := uuid.New()
	seedTicket(t, uuid.New(), itemID, firstInLine, domain.QueueStatusQueued, now.Add(time.Second), nil)

	myUserID := uuid.New()
	seedTicket(t, uuid.New(), itemID, myUserID, domain.QueueStatusQueued, now.Add(2*time.Second), nil)

	usecase := newGetMyTicketUsecase(3 * time.Second)
	ticket, err := usecase.GetMyTicket(context.Background(), itemID, myUserID)
	require.NoError(t, err)

	assert.Equal(t, domain.QueueStatusQueued, ticket.Status)
	require.NotNil(t, ticket.Position)
	assert.Equal(t, 2, *ticket.Position)
	require.NotNil(t, ticket.NextSlotFreeInSeconds)
	assert.LessOrEqual(t, *ticket.NextSlotFreeInSeconds, int64(30))
	assert.Greater(t, *ticket.NextSlotFreeInSeconds, int64(0))
	assert.Nil(t, ticket.ExpiresInSeconds)
}

func TestGetMyTicketUsecase_Offered(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 5, true)

	userID := uuid.New()
	expiresAt := time.Now().UTC().Add(10 * time.Second)
	seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusOffered, time.Now().UTC(), &expiresAt)

	usecase := newGetMyTicketUsecase(3 * time.Second)
	ticket, err := usecase.GetMyTicket(context.Background(), itemID, userID)
	require.NoError(t, err)

	assert.Equal(t, domain.QueueStatusOffered, ticket.Status)
	assert.Nil(t, ticket.Position)
	require.NotNil(t, ticket.ExpiresInSeconds)
	assert.LessOrEqual(t, *ticket.ExpiresInSeconds, int64(10))
	assert.Greater(t, *ticket.ExpiresInSeconds, int64(0))
}

func TestGetMyTicketUsecase_NoTicket(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1, true)

	usecase := newGetMyTicketUsecase(3 * time.Second)
	_, err := usecase.GetMyTicket(context.Background(), itemID, uuid.New())
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrTicketNotFound)
}

func TestGetMyTicketUsecase_TerminalTicketIsStillReturned(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1, true)

	userID := uuid.New()
	seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusCancelled, time.Now().UTC(), nil)

	usecase := newGetMyTicketUsecase(3 * time.Second)
	ticket, err := usecase.GetMyTicket(context.Background(), itemID, userID)
	require.NoError(t, err)
	assert.Equal(t, domain.QueueStatusCancelled, ticket.Status)
}

func TestGetMyTicketUsecase_PromotesOwnTicketWithinSameCall(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1, true)

	userID := uuid.New()
	seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusQueued, time.Now().UTC(), nil)

	usecase := newGetMyTicketUsecase(3 * time.Second)
	ticket, err := usecase.GetMyTicket(context.Background(), itemID, userID)
	require.NoError(t, err)

	assert.Equal(t, domain.QueueStatusOffered, ticket.Status)
	assert.Nil(t, ticket.Position)
	require.NotNil(t, ticket.ExpiresInSeconds)
	assert.Greater(t, *ticket.ExpiresInSeconds, int64(0))
}

func TestGetMyTicketUsecase_SoldOut(t *testing.T) {
	truncate(db, t)

	itemID := uuid.New()
	seedItem(t, itemID, 1, true)

	seedTicket(t, uuid.New(), itemID, uuid.New(), domain.QueueStatusPurchased, time.Now().UTC(), nil)

	userID := uuid.New()
	seedTicket(t, uuid.New(), itemID, userID, domain.QueueStatusQueued, time.Now().UTC(), nil)

	usecase := newGetMyTicketUsecase(3 * time.Second)
	ticket, err := usecase.GetMyTicket(context.Background(), itemID, userID)
	require.NoError(t, err)
	assert.Equal(t, domain.QueueStatusSoldOut, ticket.Status)
}

func TestGetMyTicketUsecase_ItemNotFound(t *testing.T) {
	truncate(db, t)

	usecase := newGetMyTicketUsecase(3 * time.Second)
	_, err := usecase.GetMyTicket(context.Background(), uuid.New(), uuid.New())
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrTicketNotFound)
}
