package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Connection struct {
	ws               *websocket.Conn
	field            *Field
	stopCh           chan struct{}
	isGameInProgress bool
}

type Server struct {
	port  string
	conns map[string]*Connection
}

var upgrader = websocket.Upgrader{
	// Разрешаем подключения с любого источника
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewServer(port string) *Server {
	return &Server{
		port:  port,
		conns: make(map[string]*Connection),
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/next", s.HandleComputeNextForm)
	mux.HandleFunc("/ws", s.handleConnections)
	mux.Handle("/", http.FileServer(http.Dir("../front")))

	http.ListenAndServe(s.port, mux)
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Ошибка при апгрейде соединения: %v", err)
	}
	defer ws.Close()

	fmt.Println("Новое WebSocket соединение")

	id := fmt.Sprintf("%v", time.Now().Nanosecond())

	conn := &Connection{ws: ws, field: NewField(100), stopCh: make(chan struct{})}
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

func (s *Server) HandleComputeNextForm(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Data string `json:"data"`
	}

	var b Body

	if r.Method != http.MethodPost {
		w.WriteHeader(404)
		w.Write([]byte("Only Post"))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&b); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid Data"))
		return
	}

	f := NewField(100)
	f.update(b.Data)
	f.run()

	type Resp struct {
		Data string `json:"data"`
	}

	resp := Resp{Data: f.getStateString()}

	m, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Write(m)
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
			conn.field.update(msg.Data)
			go func() {
				ticker := time.NewTicker(1 * time.Second)
				defer func() {
					ticker.Stop()
					fmt.Println("exit goroutine")
				}()

				for {
					select {
					case <-ticker.C:
						conn.field.run()
						msg, _ := json.Marshal(Message{Data: conn.field.getStateString(), Event: "update"})
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
