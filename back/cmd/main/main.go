package main

import (
	"game_of_life/api"
	"game_of_life/db"
	"log"
)

func main() {
	storage, err := db.NewDb()
	if err != nil {
		log.Fatal("error in create storage", err)
	}
	defer storage.Con.Close()

	server := api.NewServer(":3000", &storage)

	err = server.Run()
	if err != nil {
		log.Fatal("error in run", err)
	}
}
