package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

type quitUsecaseStub struct {
	response *dto.QueueResponse
	err      error
	called   bool
	itemID   uuid.UUID
	userID   uuid.UUID
}

func (s *quitUsecaseStub) QuitQueue(_ context.Context, itemID, userID uuid.UUID) (*dto.QueueResponse, error) {
	s.called = true
	s.itemID = itemID
	s.userID = userID
	return s.response, s.err
}

func quitRouter(usecase *quitUsecaseStub) http.Handler {
	router := chi.NewRouter()
	router.With(httpx.SessionAuth).Delete("/api/v1/items/{itemId}/queue/me", QuitQueue(usecase))
	return router
}

func quitRequest(itemID string, userID uuid.UUID) *http.Request {
	request := httptest.NewRequest(http.MethodDelete, "/api/v1/items/"+itemID+"/queue/me", nil)
	request.Header.Set("X-User-ID", userID.String())
	return request
}

func TestQuitQueueHandler_Success(t *testing.T) {
	itemID := uuid.New()
	userID := uuid.New()
	usecase := &quitUsecaseStub{response: dto.NewQueueResponse(&domain.Queue{
		ID:        uuid.New(),
		ProductID: itemID,
		UserID:    userID,
		Status:    domain.QueueStatusCancelled,
	})}
	recorder := httptest.NewRecorder()

	quitRouter(usecase).ServeHTTP(recorder, quitRequest(itemID.String(), userID))

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.Contains(t, recorder.Body.String(), `"status": "CANCELLED"`)
	assert.True(t, usecase.called)
	assert.Equal(t, itemID, usecase.itemID)
	assert.Equal(t, userID, usecase.userID)
}

func TestQuitQueueHandler_InvalidItemID(t *testing.T) {
	usecase := &quitUsecaseStub{}
	recorder := httptest.NewRecorder()

	quitRouter(usecase).ServeHTTP(recorder, quitRequest("not-a-uuid", uuid.New()))

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.False(t, usecase.called)
}

func TestQuitQueueHandler_Unauthorized(t *testing.T) {
	usecase := &quitUsecaseStub{}
	request := httptest.NewRequest(http.MethodDelete, "/api/v1/items/"+uuid.NewString()+"/queue/me", nil)
	recorder := httptest.NewRecorder()

	quitRouter(usecase).ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.False(t, usecase.called)
}

func TestQuitQueueHandler_UsecaseErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
	}{
		{name: "queue not found", err: domain.ErrQueueNotFound, statusCode: http.StatusNotFound},
		{name: "internal error", err: errors.New("unexpected"), statusCode: http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			usecase := &quitUsecaseStub{err: test.err}
			recorder := httptest.NewRecorder()

			quitRouter(usecase).ServeHTTP(recorder, quitRequest(uuid.NewString(), uuid.New()))

			assert.Equal(t, test.statusCode, recorder.Code)
			require.True(t, usecase.called)
		})
	}
}
