package httpx

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	bearerScheme        = "Bearer"
)

type userIDKey struct{}

type TokenVerifier interface {
	Verify(token string) (uuid.UUID, error)
}

func JWTAuth(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r.Header.Get(authorizationHeader))
			if !ok {
				writeUnauthorized(w, "JWT отсутствует или передан не по схеме Bearer")
				return
			}

			userID, err := verifier.Verify(token)
			if err != nil {
				writeUnauthorized(w, "JWT повреждён, просрочен или не прошёл проверку")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], bearerScheme) {
		return "", false
	}

	if parts[1] == "" {
		return "", false
	}

	return parts[1], true
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	EncodeJSON(w, struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}{
		Error:   "unauthorized",
		Message: message,
	}, http.StatusUnauthorized)
}

// UserID достаёт user_id, положенный JWTAuth.
func UserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return userID, ok
}
