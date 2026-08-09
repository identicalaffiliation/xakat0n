package httpx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/logctx"
)

type tokenVerifierStub struct {
	userID uuid.UUID
	err    error
	token  string
	called bool
}

func (s *tokenVerifierStub) Verify(token string) (uuid.UUID, error) {
	s.called = true
	s.token = token
	return s.userID, s.err
}

func TestJWTAuth_AllowsValidBearerToken(t *testing.T) {
	userID := uuid.New()
	verifier := &tokenVerifierStub{userID: userID}
	logContext := logctx.NewLogCtx()
	nextCalled := false
	handler := JWTAuth(logContext, verifier)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		actualUserID, ok := UserID(r.Context())
		require.True(t, ok)
		assert.Equal(t, userID, actualUserID)
		assert.Equal(t, userID, logContext.FieldsFromContext(r.Context())["userID"])
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
	request.Header.Set("Authorization", "Bearer signed.jwt.token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.True(t, nextCalled)
	assert.True(t, verifier.called)
	assert.Equal(t, "signed.jwt.token", verifier.token)
}

func TestJWTAuth_RejectsMissingOrMalformedBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{name: "missing"},
		{name: "token without scheme", header: "signed.jwt.token"},
		{name: "wrong scheme", header: "Basic signed.jwt.token"},
		{name: "empty bearer", header: "Bearer"},
		{name: "extra fields", header: "Bearer signed.jwt.token extra"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier := &tokenVerifierStub{}
			handler := JWTAuth(logctx.NewLogCtx(), verifier)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				t.Fatal("protected handler must not be called")
			}))
			request := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
			request.Header.Set("Authorization", test.header)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			assert.JSONEq(t, `{
				"error": "unauthorized",
				"message": "JWT отсутствует или передан не по схеме Bearer"
			}`, recorder.Body.String())
			assert.False(t, verifier.called)
		})
	}
}

func TestJWTAuth_RejectsInvalidToken(t *testing.T) {
	verifier := &tokenVerifierStub{err: errors.New("invalid signature")}
	handler := JWTAuth(logctx.NewLogCtx(), verifier)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("protected handler must not be called")
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
	request.Header.Set("Authorization", "Bearer invalid.jwt.token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
		"error": "unauthorized",
		"message": "JWT повреждён, просрочен или не прошёл проверку"
	}`, recorder.Body.String())
	assert.True(t, verifier.called)
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
