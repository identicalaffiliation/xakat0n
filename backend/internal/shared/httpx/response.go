package httpx

import (
	"encoding/json"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, val any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(val)
}
