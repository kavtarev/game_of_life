package main

import (
	"fmt"
	"game_of_life/api"
	"game_of_life/db"
	"os"
)

func main() {
	env := os.Getenv("DB_USER")
	fmt.Println(env)
	db.NewDb()
	server := api.NewServer(":3000")

	server.Run()
}
