package httpx

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/go-chi/chi/v5"
)

var inFlight atomic.Int64

func init() {
	metrics.NewCounter(`http_requests_total{url="",status="",method=""}`)
	metrics.NewHistogram(`http_request_duration_seconds{url="",method=""}`)
	metrics.NewGauge(`http_requests_in_flight`, func() float64 {
		return float64(inFlight.Load())
	})
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := chi.RouteContext(r.Context()).RoutePattern()
		if url == "" {
			url = r.URL.Path
		}
		method := r.Method

		inFlight.Add(1)
		start := time.Now()

		sw := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(sw, r)

		inFlight.Add(-1)

		metrics.GetOrCreateCounter(fmt.Sprintf(
			`http_requests_total{url=%q,status=%q,method=%q}`, url, strconv.Itoa(sw.code()), method,
		)).Inc()
		metrics.GetOrCreateHistogram(fmt.Sprintf(
			`http_request_duration_seconds{url=%q,method=%q}`, url, method,
		)).Update(time.Since(start).Seconds())
	})
}

func MetricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	metrics.WritePrometheus(w, false)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusRecorder) code() int {
	if w.status == 0 {
		return http.StatusOK
	}

	return w.status
}
