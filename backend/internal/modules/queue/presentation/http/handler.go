package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

const (
	ItemIdMuxPattern = "itemId"
)

func PutUserInQueue(usecase ports.CreateUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, err := uuid.Parse(chi.URLParam(r, ItemIdMuxPattern))
		if err != nil {
			http.Error(w, "invalid product id", http.StatusBadRequest)
			return
		}

		userID, ok := httpx.UserID(r.Context())
		if !ok {
			return
		}

		response, err := usecase.CreateQueue(r.Context(), dto.NewCreateRequest(productID, userID))
		if err != nil {
			if errors.Is(err, domain.ErrUserAlreadyQueued) {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		httpx.EncodeJSON(w, response, http.StatusCreated)
	}
}

func QuitQueue(usecase ports.QuitUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, userID, ok := queueIDs(w, r)
		if !ok {
			return
		}

		response, err := usecase.QuitQueue(r.Context(), dto.NewQuitQueueRequest(productID, userID))
		if err != nil {
			if errors.Is(err, domain.ErrQueueNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)

				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		httpx.EncodeJSON(w, response, http.StatusOK)
	}
}

func queueIDs(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	productID, err := uuid.Parse(chi.URLParam(r, ItemIdMuxPattern))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)

		return uuid.Nil, uuid.Nil, false
	}

	userID, ok := httpx.UserID(r.Context())
	if !ok {
		http.Error(w, "failed to parse user id", http.StatusUnauthorized)

		return uuid.Nil, uuid.Nil, false
	}

	return productID, userID, true
}
