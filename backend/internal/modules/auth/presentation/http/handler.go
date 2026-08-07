package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/ports"
	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

func Login(issuer ports.TokenIssuer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		username, err := domain.NewUsername(request.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := uuid.New()

		token, err := issuer.Issue(userID, username)
		if err != nil {
			http.Error(w, "failed to issue token", http.StatusInternalServerError)
			return
		}

		response := dto.NewLoginResponse(userID, request.Username, token)
		httpx.EncodeJSON(w, response, http.StatusOK)
	}
}
