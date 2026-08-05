package controller

import (
	"errors"
	"net/http"

	"github.com/identicalaffiliation/xakat0n/backend/internal/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/ports"
)

const (
	ProductIdMuxPattern = "productId"
)

func PutUserInQueue(usecase ports.CreateUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, userID, ok := parseQueueIDs(w, r)
		if !ok {
			return
		}

		rCtx := r.Context()
		response, err := usecase.CreateQueue(rCtx, dto.NewCreateRequest(productID, userID))
		if err != nil {
			if errors.Is(err, domain.ErrUserAlreadyQueued) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodeJSON(w, response, http.StatusCreated)
	}
}
