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
			httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid item id")
			return
		}

		userID, ok := httpx.UserID(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unathorized", "invalid user id")
			return
		}

		rCtx := r.Context()
		response, err := usecase.CreateQueue(rCtx, dto.NewCreateRequest(productID, userID))
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrQueueNotApplicable),
				errors.Is(err, domain.ErrItemSoldOut):
				httpx.WriteError(w, http.StatusConflict, "queue_not_applicable", "queue not applicable")
			case errors.Is(err, domain.ErrItemNotFound):
				httpx.WriteError(w, http.StatusNotFound, "item_not_found", "item not found")
			default:
				httpx.WriteError(w, http.StatusInternalServerError, "internal_server_error", "internal server error")
			}
			return
		}

		httpx.EncodeJSON(w, response, http.StatusOK)
	}
}
