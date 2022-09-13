package main

import (
	"reconmap/agent/internal"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	app := internal.NewApp()
	if err := app.Run(); err != nil {
		log.Error(*err)
	}
}
