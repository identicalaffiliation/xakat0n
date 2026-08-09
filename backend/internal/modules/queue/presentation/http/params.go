package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/httpx"
)

const ItemIDParam = "itemId"

func parseQueueIDs(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	itemID, err := uuid.Parse(chi.URLParam(r, ItemIDParam))
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
