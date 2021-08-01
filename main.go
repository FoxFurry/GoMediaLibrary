package main

import (
	"github.com/foxfurry/simple-rest/app/app_config"
)

func main(){
	server := app_config.NewApp()
	server.Start()
}