package main

import (
	"authapp/cmd"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("App starting")
	cmd.Execute()
}
