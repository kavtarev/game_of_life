package main

import (
	"game_of_life/api"
	"game_of_life/db"
	"log"
)

func main() {
	storage, err := db.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Con.Close()

	server := api.NewServer(":3000", &storage)

	server.Run()
}
