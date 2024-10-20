package main

import (
	"fmt"
	"game_of_life/api"
	"os"
)

func main() {
	env := os.Getenv("DB_USER")
	fmt.Println(env)
	server := api.NewServer(":3000")

	server.Run()
}
