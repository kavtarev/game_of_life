package main

import ()

func main() {
	server := NewServer(":3000")
	server.Run()
}
