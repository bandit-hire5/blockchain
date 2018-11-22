package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	LedgerPath string
	Port       string
	Pattern 	string
}

func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
