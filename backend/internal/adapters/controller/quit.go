package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/ports"
)

func QuitQueue(usecase ports.QuitUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, userID, ok := parseQueueIDs(w, r)
		if !ok {
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
