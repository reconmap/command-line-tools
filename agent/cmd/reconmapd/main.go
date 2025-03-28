package main

import (
	"reconmap/agent/internal"
)

func main() {
	app := internal.NewApp()
	if err := app.Run(); err != nil {
		app.Logger.Error(*err)
	}
}
