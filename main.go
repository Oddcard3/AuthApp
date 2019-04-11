package main

import (
	"authapp/config"
	"authapp/server"
)

func main() {
	config.Init()

	s, _ := server.NewServer()
	s.Start()
}
