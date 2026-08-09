package handler

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

func QuitQueue(usecase ports.QuitUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID, userID, ok := parseQueueIDs(w, r)
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
