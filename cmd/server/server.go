package main

import "github.com/NoSoundLeR/basket.git/server"

func main() {
	server := server.NewServer()
	server.Run()
}
