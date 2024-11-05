package api

import (
	"net/http"

	"game_of_life/api/handlers"
	"game_of_life/db"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	port       string
	conns      map[string]*Connection
	storage    *db.Storage
	prometheus *prometheus.Registry
	metrics    *metrics
}

type metrics struct {
	gauge prometheus.Gauge
	info  *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		gauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "game_of_life", Name: "default_gauge", Help: "default descriptor of gauge",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "game_of_life",
			Name:      "info",
			Help:      "Information about the My App environment.",
		},
			[]string{"version"}),
	}
	reg.MustRegister(m.gauge, m.info)
	return m
}

func NewServer(port string, s *db.Storage) *Server {
	server := &Server{
		port:       port,
		conns:      make(map[string]*Connection),
		storage:    s,
		prometheus: prometheus.NewRegistry(),
	}

	server.metrics = NewMetrics(server.prometheus)
	return server
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	s.metrics.info.With(prometheus.Labels{"version": "1.0.0"}).Set(1)

	mux.HandleFunc("/next", handlers.HandleComputeNextForm)
	mux.HandleFunc("/ws", s.handleConnections)
	mux.Handle("/", http.FileServer(http.Dir("../front")))

	mux.HandleFunc("/count", s.IncrementMiddleware(handlers.HandleCount))

	// Expose /metrics HTTP endpoint using the created custom registry.
	mux.Handle(
		"/metrics", promhttp.HandlerFor(
			s.prometheus,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)

	http.ListenAndServe(s.port, mux)
}
