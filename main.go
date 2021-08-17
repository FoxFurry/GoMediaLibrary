package main

import (
	"github.com/foxfurry/simple-rest/app"
	"github.com/foxfurry/simple-rest/configs"
)

// Do I need to explain this?
func main() {
	configs.LoadConfig()
	server := app.NewApp()
	server.Start()
}

/*
integration tests/ unit tests/ functional tests
 */