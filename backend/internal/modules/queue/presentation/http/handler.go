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
		itemID, err := uuid.Parse(chi.URLParam(r, ItemIdMuxPattern))
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid item id")
			return
		}

		userID, ok := httpx.UserID(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid user id")
			return
		}

		response, err := usecase.CreateQueue(r.Context(), dto.NewCreateRequest(itemID, userID))
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrQueueNotApplicable):
				httpx.WriteError(w, http.StatusConflict, "queue_not_applicable", "item is not limited, queue does not apply")
			case errors.Is(err, domain.ErrItemSoldOut):
				httpx.WriteError(w, http.StatusConflict, "item_sold_out", "item is completely sold out")
			case errors.Is(err, domain.ErrItemNotFound):
				httpx.WriteError(w, http.StatusNotFound, "item_not_found", "item not found")
			default:
				httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
			return
		}

		httpx.EncodeJSON(w, response, http.StatusOK)
	}
}

func QuitQueue(usecase ports.QuitUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, userID, ok := queueIDs(w, r)
		if !ok {
			return
		}

		response, err := usecase.QuitQueue(r.Context(), itemID, userID)
		if err != nil {
			if errors.Is(err, domain.ErrQueueNotFound) {
				httpx.WriteError(w, http.StatusNotFound, "ticket_not_found", "user has no ticket for this item")

				return
			}

			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")

			return
		}

		httpx.EncodeJSON(w, response, http.StatusOK)
	}
}

func queueIDs(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	itemID, err := uuid.Parse(chi.URLParam(r, ItemIdMuxPattern))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid item id")

		return uuid.Nil, uuid.Nil, false
	}

	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid session")

		return uuid.Nil, uuid.Nil, false
	}

	return itemID, userID, true
}
