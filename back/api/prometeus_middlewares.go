package api

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (s *Server) IncrementMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.metrics.gauge.Inc()
		fn(w, r)
	}
}

func (s *Server) HistogramMiddleware(fn http.HandlerFunc, name string) http.HandlerFunc {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "game_of_life",
		Name:      name,
		Help:      "Duration of the request.",
		// 4 times larger for apdex score
		// Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
		// Buckets: prometheus.LinearBuckets(0.1, 5, 5),
		Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
	}, []string{"status", "method"})
	s.prometheus.Register(histogram)
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		fn(w, r)
		histogram.With(prometheus.Labels{"method": "GET", "status": "200"}).Observe(time.Since(now).Seconds())
	}
}
