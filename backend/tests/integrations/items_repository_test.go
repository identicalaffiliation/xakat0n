package integrations

import (
	"context"
	"testing"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemsRepository_CreateItem(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := postgres.NewItemsRepository(db)
	require.NoError(t, repo.CreateItem(ctx, domain.NewItem(
		"test title",
		"test desc",
		1000,
	)))
}

func TestItemsRepository_GetAll(t *testing.T) {
	truncate(db, t)

	ctx := context.Background()
	repo := postgres.NewItemsRepository(db)
	expected := domain.NewItem(
		"test title",
		"test desc",
		1000,
	)

	require.NoError(t, repo.CreateItem(ctx, expected))

	actual, err := repo.GetAll(ctx)
	require.NoError(t, err)

	assert.NotNil(t, actual)
	assert.Equal(t, expected.Title, actual[0].Title)
	assert.Equal(t, expected.Description, actual[0].Description)
	assert.Equal(t, expected.Price, actual[0].Price)
}
