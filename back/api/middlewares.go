package api

import (
	"game_of_life/db"
	"net/http"
)

func (s *Server) mapperWithStorage(f func(w http.ResponseWriter, r *http.Request, db *db.Storage)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, s.storage)
	}
}
