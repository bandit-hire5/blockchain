package middlewares

import (
	"context"
	"net/http"

	"github.com/bandit/blockchain-core"
)

type handler func(http.ResponseWriter, *http.Request)

const ctxLedgerKey = "ledger"

func AddContextLedger(ledger *core.Ledger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxLedgerKey, ledger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Ledger(r *http.Request) *core.Ledger {
	return r.Context().Value(ctxLedgerKey).(*core.Ledger)
}
