package middlewares

import (
	"context"
	l "naeltok/go-blockchain/ledger"
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)

const ctxLedgerKey = "ledger"

func AddContextLedger(ledger *l.Ledger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxLedgerKey, ledger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Ledger(r *http.Request) *l.Ledger {
	return r.Context().Value(ctxLedgerKey).(*l.Ledger)
}
