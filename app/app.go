package app

import (
	"log"
	"net/http"

	"github.com/bandit/blockchain-core"
	c "github.com/bandit/blockchain/config"
)

type App struct {
	config c.Config
	Ledger *core.Ledger
}

func NewApp(config c.Config) *App {
	app := &App{
		config: config,
	}

	app.initLedger()

	return app
}

func (a *App) Server() {
	server := NewServer("/entry", a.Ledger)
	go server.Listen()

	//router := r.Router(a.Ledger)
	//http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":"+a.config.Port, nil))
}

func (a *App) initLedger() {
	ledger := core.NewLedger(a.config.LedgerPath)
	ledger.CreateWithGenesisBlock()

	a.Ledger = ledger
}
