package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

func Login(usecase ports.LoginUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid request body")
			return
		}

		response, err := usecase.Login(r.Context(), request.Username)
		if err != nil {
			if errors.Is(err, domain.ErrInternal) {
				httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to issue token")
				return
			}

			httpx.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		httpx.EncodeJSON(w, response, http.StatusOK)
	}
}
