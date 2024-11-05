package handlers

import "net/http"

func HandleCount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("count++"))
}
