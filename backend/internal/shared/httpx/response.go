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

// ErrorResponse — JSON-тело ошибки по схеме Error из api-contract.yaml.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// WriteError кодирует ошибку в контрактный JSON {error, message} вместо голого текста.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	EncodeJSON(w, ErrorResponse{Error: code, Message: message}, status)
}
