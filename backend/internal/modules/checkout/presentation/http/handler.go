package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

const (
	ItemIdMuxPattern = "itemId"
)

func StartCheckout(usecase ports.CheckoutUsecase) http.HandlerFunc {
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

		result, err := usecase.StartCheckout(r.Context(), itemID, userID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrItemNotFound):
				httpx.WriteError(w, http.StatusNotFound, "item_not_found", "item not found")
			case errors.Is(err, domain.ErrNoActiveRight):
				httpx.WriteError(w, http.StatusForbidden, "no_active_right", "you have no active right to purchase this item")
			default:
				httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
			return
		}

		httpx.EncodeJSON(w, result, http.StatusOK)
	}
}

func PaymentCallback(usecase ports.PaymentCallbackUsecase) http.HandlerFunc {
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

		var in dto.PaymentCallbackRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || !in.Valid() {
			httpx.WriteError(w, http.StatusBadRequest, "bad_request", "ticketId is required, result must be either \"paid\" or \"failed\"")
			return
		}

		ticket, err := usecase.HandleCallback(r.Context(), itemID, userID, &in)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrTicketNotFound):
				httpx.WriteError(w, http.StatusNotFound, "ticket_not_found", "user has no ticket for this item")
			case errors.Is(err, domain.ErrTooLate):
				httpx.WriteError(w, http.StatusConflict, "checkout_result_too_late", "window already expired, item was given to the next in line")
			default:
				httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
			return
		}

		httpx.EncodeJSON(w, ticket, http.StatusOK)
	}
}
