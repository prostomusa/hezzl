package main

import (
	"Project/internal/app"
)

func main() {
	config := app.NewConfig()
	serv := app.New(config)
	serv.Start()
}
