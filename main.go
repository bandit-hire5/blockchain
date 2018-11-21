package main

import (
	"github.com/bandit/blockchain/app"
	. "github.com/bandit/blockchain/config"
)

func main() {
	var config = Config{}
	config.Read()

	app := app.NewApp(config)
	app.Server()
}
