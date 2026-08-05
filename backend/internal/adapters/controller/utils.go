package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func encodeJSON(w http.ResponseWriter, val any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(val)
}

func parseQueueIDs(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	productID, err := uuid.Parse(chi.URLParam(r, ProductIdMuxPattern))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)

		return uuid.Nil, uuid.Nil, false
	}

	userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
	if err != nil {
		http.Error(w, "failed to parse user id", http.StatusUnauthorized)

		return uuid.Nil, uuid.Nil, false
	}

	return productID, userID, true
}
