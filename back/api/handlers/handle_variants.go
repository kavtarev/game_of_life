package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HandleString(w http.ResponseWriter, r *http.Request) {
	var b Dto

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

}

func HandleBytes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post", http.StatusNotFound)
		return
	}

	maxSize := 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize))
	bReader := bufio.NewReader(r.Body)
	var res []byte

	for {
		temp := make([]byte, 10)

		n, err := bReader.Read(temp)
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		if n == 0 {
			break
		}
		res = append(res, temp[:n]...)
	}
	defer r.Body.Close()

	s := strings.Builder{}
	s.Grow(len(res) * 8)

	for i := 0; i < len(res); i++ {
		for j := 0; j < 8; j++ {
			s.WriteString(fmt.Sprintf("%v", (int(res[i])>>(7-j))&1))
		}
	}
}

func HandleBytesImproved(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post", http.StatusNotFound)
		return
	}

	const maxSize = 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	s := strings.Builder{}
	s.Grow(len(body) * 8)

	for _, b := range body {
		s.WriteString(fmt.Sprintf("%08b", b))
	}
}
