package api

import (
	"net/http"
)

func (s *Server) IncrementMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.metrics.gauge.Inc()
		fn(w, r)
	}
}
