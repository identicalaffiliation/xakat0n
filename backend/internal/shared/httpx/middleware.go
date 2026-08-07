package httpx

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const sessionHeader = "X-User-ID"

type userIDKey struct{}

// SessionAuth — заглушка полноценной сессионной авторизации: значение заголовка
// парсится напрямую как user_id, без резолва через таблицу сессий. Реальные сессии
// (internal/auth) — отдельная будущая работа.
func SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.Header.Get(sessionHeader))
		if err != nil {
			http.Error(w, "failed to parse user id", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey{}, userID)))
	})
}

// UserID достаёт user_id, положенный SessionAuth.
func UserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return userID, ok
}
