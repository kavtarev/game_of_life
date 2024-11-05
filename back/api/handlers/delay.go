package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func HandleDelay(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(200+rand.Intn(50)) * time.Millisecond)
	fmt.Println("done++")
	w.Write([]byte("done"))
}
