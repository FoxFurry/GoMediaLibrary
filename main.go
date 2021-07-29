package main

import "github.com/foxfurry/simple-rest/app"

func main(){
	server := app.NewApp()
	server.Start()
}