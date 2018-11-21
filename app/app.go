package app

import (
	"log"
	c "naeltok/go-blockchain/config"
	l "naeltok/go-blockchain/ledger"
	"net/http"
)

type App struct {
	config c.Config
	Ledger *l.Ledger
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
	ledger := l.NewLedger(a.config.LedgerPath)
	ledger.Create()

	a.Ledger = ledger
}
