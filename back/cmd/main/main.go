package main

import (
	"game_of_life/api"
)

func main() {
	server := api.NewServer(":3000")

	server.Run()
}
