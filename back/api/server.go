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
	port    string
	conns   map[string]*Connection
	storage *db.Storage
}

var upgrader = websocket.Upgrader{
	// Разрешаем подключения с любого источника
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewServer(port string, s *db.Storage) *Server {
	return &Server{
		port:    port,
		conns:   make(map[string]*Connection),
		storage: s,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/next", handlers.HandleComputeNextForm)
	mux.HandleFunc("/ws", s.handleConnections)
	mux.Handle("/", http.FileServer(http.Dir("../front")))

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
		conn.stopCh <- struct{}{}
	}

}
