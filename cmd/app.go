package main

import (
	"log"
	"net/http"

	"github.com/bandit/blockchain-core"
	"github.com/bandit/blockchain/config"
	ws "github.com/bandit/blockchain/websocket"
)

type App struct {
	config config.Config
	ledger *core.Ledger
}

func NewApp(cfg config.Config) *App {
	app := &App{
		config: cfg,
	}

	app.initLedger()

	return app
}

func (self *App) Server() {
	server := ws.NewServer()
	server.Pattern = self.config.Pattern
	server.Ledger = self.ledger

	go server.Listen()

	log.Fatal(http.ListenAndServe(":"+self.config.Port, nil))
}

func (self *App) initLedger() {
	ledger := core.NewLedger(self.config.LedgerPath)

	err := ledger.CreateWithGenesisBlock()
	if err != nil {
		panic(err)
	}

	self.ledger = ledger
}
