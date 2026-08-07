package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionAuth(t *testing.T) {
	t.Parallel()

	t.Run("valid header puts user id into context", func(t *testing.T) {
		userID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(sessionHeader, userID.String())

		var got uuid.UUID
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got, _ = UserID(r.Context())
		})

		rec := httptest.NewRecorder()
		SessionAuth(next).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userID, got)
	})

	t.Run("missing header is unauthorized", func(t *testing.T) {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		rec := httptest.NewRecorder()
		SessionAuth(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.False(t, called)
	})

	t.Run("malformed header is unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(sessionHeader, "not-a-uuid")

		rec := httptest.NewRecorder()
		SessionAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next must not be called")
		})).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestUserID(t *testing.T) {
	t.Parallel()

	t.Run("missing value returns false", func(t *testing.T) {
		_, ok := UserID(context.Background())
		assert.False(t, ok)
	})

	t.Run("wrong type returns false", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey{}, "not-a-uuid")
		_, ok := UserID(ctx)
		assert.False(t, ok)
	})

	t.Run("stored uuid is returned", func(t *testing.T) {
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), userIDKey{}, userID)

		got, ok := UserID(ctx)
		require.True(t, ok)
		assert.Equal(t, userID, got)
	})
}
