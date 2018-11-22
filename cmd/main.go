package main

import (
	. "github.com/bandit/blockchain/config"
)

func main() {
	var config = Config{}
	config.Read()

	app := NewApp(config)
	app.Server()
}
