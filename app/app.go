package app

import (
	"log"
	c "naeltok/go-blockchain/config"
	l "naeltok/go-blockchain/ledger"
	r "naeltok/go-blockchain/routes"
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
	router := r.Router(a.Ledger)

	if err := http.ListenAndServe(":"+a.config.Port, router); err != nil {
		log.Fatal(err)
	}
}

func (a *App) initLedger() {
	ledger := l.NewLedger(a.config.LedgerPath)
	ledger.Create()

	a.Ledger = ledger
}
