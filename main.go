package main

import (
	"github.com/foxfurry/medialib/app"
	"github.com/foxfurry/medialib/configs"
)

func init() {
	configs.LoadConfig()
}

func main() {
	server := app.NewApp()
	server.Start()
}

/*
integration tests/ unit tests/ functional tests
*/
