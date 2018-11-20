package main

import (
	"naeltok/go-blockchain/app"
	. "naeltok/go-blockchain/config"
)

func main() {
	var config = Config{}
	config.Read()

	app := app.NewApp(config)
	app.Server()
}
