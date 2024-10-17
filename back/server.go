package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Connection struct {
	ws     *websocket.Conn
	field  *Field
	stopCh chan struct{}
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

	conn := &Connection{ws: ws, field: NewField(10), stopCh: make(chan struct{})}
	s.conns[id] = conn
	s.sendEvent(ws, "init", "init")
	s.listenToSocket(ws, conn)

}

func (s *Server) sendEvent(ws *websocket.Conn, data, event string) {
	msg, _ := json.Marshal(Message{Data: data, Event: event})
	ws.Write(msg)
}

func (s *Server) listenToSocket(ws *websocket.Conn, conn *Connection) {
	buf := make([]byte, 1000000)

	for {
		b, err := ws.Read(buf)

		if err != nil {
			id, findErr := s.findId(ws)
			if findErr != nil {
				return
			}

			delete(s.conns, id)

			if err == io.EOF {
				fmt.Println("eof error, client closed connection")
				break
			}

		}
		s.parseMessage(ws, buf[:b], conn)
	}

}

func (s *Server) findId(ws *websocket.Conn) (string, error) {
	for id, v := range s.conns {
		if v.ws == ws {
			return id, nil
		}
	}

	return "", errors.New("not found connection to delete")
}

func (s *Server) parseMessage(ws *websocket.Conn, b []byte, conn *Connection) {
	msg := Message{}

	if err := json.Unmarshal(b, &msg); err != nil {
		panic(err)
	}

	switch msg.Event {
	case "start":
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
					ws.Write(msg)
				case <-conn.stopCh:
					return
				}
			}
		}()
	case "stop":
		conn.stopCh <- struct{}{}
	}

}
