package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/ports"
)

func QuitQueue(usecase ports.QuitUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, err := uuid.Parse(chi.URLParam(r, ProductIdMuxPattern))
		if err != nil {
			http.Error(w, "invalid product id", http.StatusBadRequest)
			return
		}

		userIDS := r.Header.Get("X-User-ID")

		userID, err := uuid.Parse(userIDS)
		if err != nil {
			http.Error(w, "failed to parse user id", http.StatusUnauthorized)
			return
		}

		rCtx := r.Context()
		response, err := usecase.QuitQueue(rCtx, dto.NewQuitQueueRequest(productID, userID))
		if err != nil {
			if errors.Is(err, domain.ErrQueueNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encodeJSON(w, response, http.StatusOK)
	}
}
