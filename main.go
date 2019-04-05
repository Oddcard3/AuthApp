package main

import (
	"authapp/server"
)

func main() {
	s, _ := server.NewServer()
	s.Start()
}