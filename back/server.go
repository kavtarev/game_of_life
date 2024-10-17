package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Connection struct {
	ws *websocket.Conn
}

type Server struct {
	port  string
	conns map[string]*Connection
}

func NewServer(port string) *Server {
	return &Server{
		port:  port,
		conns: make(map[string]*Connection),
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ss", s.HandleStatic)
	mux.Handle("/ws", websocket.Handler(s.HandleConn))
	mux.Handle("/", http.FileServer(http.Dir("../front")))

	http.ListenAndServe(s.port, mux)
}

func (s *Server) HandleStatic(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("works"))
}

func (s *Server) HandleConn(ws *websocket.Conn) {
	id := fmt.Sprintf("%v", time.Now().Nanosecond())

	s.conns[id] = &Connection{ws: ws}
	fmt.Println(66666, s.conns)

	s.sendEvent(ws, "init", "init")
	s.listenToSocket(ws)

}

func (s *Server) sendEvent(ws *websocket.Conn, data, event string) {
	msg, _ := json.Marshal(Message{Data: data, Event: event})
	ws.Write(msg)
}

func (s *Server) listenToSocket(ws *websocket.Conn) {
	buf := make([]byte, 1000000)

	b, err := ws.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	s.parseMessage(ws, buf[:b])

}

func (s *Server) parseMessage(ws *websocket.Conn, b []byte) {
	msg := Message{}

	if err := json.Unmarshal(b, &msg); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", msg)
}
