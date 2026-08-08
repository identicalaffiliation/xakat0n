package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeJSON(t *testing.T) {
	t.Parallel()

	t.Run("writes status, content type and body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		payload := map[string]string{"status": "OFFERED"}

		EncodeJSON(rec, payload, http.StatusCreated)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var decoded map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &decoded))
		assert.Equal(t, "OFFERED", decoded["status"])
	})

	t.Run("indent is applied", func(t *testing.T) {
		rec := httptest.NewRecorder()

		EncodeJSON(rec, map[string]string{"a": "b"}, http.StatusOK)

		assert.Contains(t, rec.Body.String(), "\n")
		assert.Contains(t, rec.Body.String(), "\"a\": \"b\"")
	})

	t.Run("nil value encodes as null", func(t *testing.T) {
		rec := httptest.NewRecorder()

		EncodeJSON(rec, nil, http.StatusOK)

		assert.Equal(t, "null\n", rec.Body.String())
	})
}
