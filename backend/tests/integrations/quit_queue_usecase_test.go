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
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/tx"
)

func TestQuitQueueUsecase_PromotesNextUser(t *testing.T) {
	truncate(db, t)

	const ttl = 2 * time.Minute
	itemID := uuid.New()
	seedItem(t, itemID, 1)
	base := time.Now().UTC().Add(-time.Minute)
	firstUserID := uuid.New()
	secondUserID := uuid.New()
	expiresAt := time.Now().UTC().Add(time.Minute)
	seedTicket(t, uuid.New(), itemID, firstUserID, domain.QueueStatusOffered, base, &expiresAt)
	seedTicket(t, uuid.New(), itemID, secondUserID, domain.QueueStatusQueued, base.Add(time.Second), nil)

	repo := postgres.NewQueueRepository(db)
	manager := tx.NewManager(db, slog.Default())
	advance := newAdvanceQueueUsecase()
	usecase := application.NewQuitQueueUsecase(repo, advance, manager, slog.Default(), ttl)

	response, err := usecase.QuitQueue(context.Background(), itemID, firstUserID)

	require.NoError(t, err)
	assert.Equal(t, domain.QueueStatusCancelled, response.Queue.Status)

	var (
		status        domain.QueueStatus
		actualExpires *time.Time
	)
	err = db.QueryRow(context.Background(),
		`SELECT status, expires_at FROM queues WHERE product_id = $1 AND user_id = $2`,
		itemID, secondUserID,
	).Scan(&status, &actualExpires)
	require.NoError(t, err)
	assert.Equal(t, domain.QueueStatusOffered, status)
	require.NotNil(t, actualExpires)
	assert.WithinDuration(t, time.Now().Add(ttl), *actualExpires, time.Second)
}
