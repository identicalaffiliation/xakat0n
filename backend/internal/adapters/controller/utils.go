package controller

import (
	"encoding/json"
	"net/http"
)

func encodeJSON(w http.ResponseWriter, val any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(val)
}
