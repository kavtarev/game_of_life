package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"game_of_life/api/handlers"
	"game_of_life/db"
	"game_of_life/internal"

	"github.com/gorilla/websocket"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Connection struct {
	ws               *websocket.Conn
	field            *internal.Field
	stopCh           chan struct{}
	isGameInProgress bool
}

type Server struct {
	port       string
	conns      map[string]*Connection
	storage    *db.Storage
	prometheus *prometheus.Registry
	metrics    *metrics
}

var upgrader = websocket.Upgrader{
	// Разрешаем подключения с любого источника
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

func (s *Server) mapperWithStorage(f func(w http.ResponseWriter, r *http.Request, db *db.Storage)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, s.storage)
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Ошибка при апгрейде соединения: %v", err)
	}
	defer ws.Close()

	fmt.Println("Новое WebSocket соединение")

	id := fmt.Sprintf("%v", time.Now().Nanosecond())

	conn := &Connection{ws: ws, field: internal.NewField(100), stopCh: make(chan struct{})}
	s.conns[id] = conn

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			conn.stopCh <- struct{}{}
			delete(s.conns, id)
			close(conn.stopCh)
			break
		}

		s.parseMessage(ws, message, conn)
	}
}

func (s *Server) parseMessage(ws *websocket.Conn, b []byte, conn *Connection) {
	msg := Message{}

	if err := json.Unmarshal(b, &msg); err != nil {
		fmt.Println("in json unmarshal")
		panic(err)
	}

	switch msg.Event {
	case "start":
		if !conn.isGameInProgress {
			conn.isGameInProgress = true
			conn.field.Update(msg.Data)
			go func() {
				ticker := time.NewTicker(1 * time.Second)
				defer func() {
					ticker.Stop()
					fmt.Println("exit goroutine")
				}()

				for {
					select {
					case <-ticker.C:
						str := conn.field.Run()
						msg, _ := json.Marshal(Message{Data: str, Event: "update"})
						ws.WriteMessage(1, msg)
					case <-conn.stopCh:
						return
					}
				}
			}()
		}
	case "stop":
		conn.isGameInProgress = false
		conn.stopCh <- struct{}{}
	}

}
