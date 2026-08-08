package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_MiddlewareRecordsRequest(t *testing.T) {
	router := chi.NewRouter()
	router.With(Metrics).Get("/api/v1/items/{itemId}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	router.Get("/metrics", MetricsHandler)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/items/some-uuid", nil)
	router.ServeHTTP(httptest.NewRecorder(), request)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	body := recorder.Body.String()
	assert.Contains(t, body, `http_requests_total{url="/api/v1/items/{itemId}",status="404",method="GET"} 1`)
	assert.Contains(t, body, `http_request_duration_seconds_bucket{url="/api/v1/items/{itemId}",method="GET",vmrange=`)
	assert.Contains(t, body, "http_requests_in_flight 0")
}

func TestMetrics_MiddlewareDefaultsToStatusOK(t *testing.T) {
	router := chi.NewRouter()
	router.With(Metrics).Get("/api/v1/items", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	router.Get("/metrics", MetricsHandler)

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/v1/items", nil))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	assert.True(t, strings.Contains(
		recorder.Body.String(),
		`http_requests_total{url="/api/v1/items",status="200",method="GET"} 1`,
	))
}

func TestMetricsHandler_ExposesRegisteredMetrics(t *testing.T) {
	recorder := httptest.NewRecorder()
	MetricsHandler(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	assert.Equal(t, "text/plain; version=0.0.4; charset=utf-8", recorder.Header().Get("Content-Type"))
	require.Contains(t, recorder.Body.String(), "http_requests_total")
}
