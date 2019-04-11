package main

import (
	"authapp/config"
	"authapp/logging"
	"authapp/server"
)

func main() {
	config.Init()
	logging.NewLogger()

	logging.Logger.Info("App starting")
	s, _ := server.NewServer()
	s.Start()
}
