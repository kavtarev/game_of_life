package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func unpackBits(data []byte) []int {
	bits := make([]int, len(data)*8)
	for i, b := range data {
		for j := 0; j < 8; j++ {
			bits[i*8+j] = (int(b) >> (7 - j)) & 1 // Получаем каждый бит из байта
		}
	}
	return bits
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Распаковка байтов в биты
	bits := unpackBits(body)
	fmt.Printf("Received checkbox states: %v\n", bits)

	// Отправляем ответ клиенту
	response := "Checkbox states received successfully"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/receive", handler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
