package integrations

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/infrastructure/postgres"
)

func TestUsersRepository_GetOrCreate(t *testing.T) {
	t.Run("repeated login with the same username returns the same id", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewUsersRepository(db)
		username, err := domain.NewUsername("hunter_42")
		require.NoError(t, err)

		first, err := repo.GetOrCreate(ctx, username)
		require.NoError(t, err)

		second, err := repo.GetOrCreate(ctx, username)
		require.NoError(t, err)

		assert.Equal(t, first, second)
	})

	t.Run("different usernames get different ids", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewUsersRepository(db)
		alice, err := domain.NewUsername("alice")
		require.NoError(t, err)
		bob, err := domain.NewUsername("bob")
		require.NoError(t, err)

		aliceID, err := repo.GetOrCreate(ctx, alice)
		require.NoError(t, err)
		bobID, err := repo.GetOrCreate(ctx, bob)
		require.NoError(t, err)

		assert.NotEqual(t, aliceID, bobID)
	})

	t.Run("concurrent first-time logins with the same username converge on one id", func(t *testing.T) {
		truncate(db, t)

		ctx := context.Background()
		repo := postgres.NewUsersRepository(db)
		username, err := domain.NewUsername("racer")
		require.NoError(t, err)

		const n = 10
		ids := make([]string, n)
		var wg sync.WaitGroup
		wg.Add(n)
		for i := range n {
			go func() {
				defer wg.Done()
				id, err := repo.GetOrCreate(ctx, username)
				assert.NoError(t, err)
				ids[i] = id.String()
			}()
		}
		wg.Wait()

		for i := 1; i < n; i++ {
			assert.Equal(t, ids[0], ids[i])
		}
	})
}
