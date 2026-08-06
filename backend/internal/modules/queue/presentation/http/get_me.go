package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

const ItemIdMuxPattern = "itemId"

func GetMyTicket(usecase ports.GetMeUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := uuid.Parse(chi.URLParam(r, ItemIdMuxPattern))
		if err != nil {
			http.Error(w, "invalid item id", http.StatusBadRequest)
			return
		}

		userID, ok := httpx.UserID(r.Context())
		if !ok {
			http.Error(w, "failed to parse user id", http.StatusUnauthorized)
			return
		}

		ticket, err := usecase.GetMyTicket(r.Context(), itemID, userID)
		if err != nil {
			if errors.Is(err, domain.ErrTicketNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		httpx.EncodeJSON(w, ticket, http.StatusOK)
	}
}
