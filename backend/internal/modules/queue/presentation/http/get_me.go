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

func GetMyTicket(usecase ports.GetMeUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, err := uuid.Parse(chi.URLParam(r, ItemIdMuxPattern))
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid item id")
			return
		}

		userID, ok := httpx.UserID(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid session")
			return
		}

		ticket, err := usecase.GetMyTicket(r.Context(), itemID, userID)
		if err != nil {
			if errors.Is(err, domain.ErrTicketNotFound) {
				httpx.WriteError(w, http.StatusNotFound, "ticket_not_found", "user has no ticket for this item")
				return
			}

			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			return
		}

		httpx.EncodeJSON(w, ticket, http.StatusOK)
	}
}
