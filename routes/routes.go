package routes

import (
	"encoding/json"
	"net/http"

	"github.com/bandit/blockchain-core"
	m "github.com/bandit/blockchain/middlewares"

	"github.com/gorilla/mux"
)

func Router(l *core.Ledger) *mux.Router {
	r := mux.NewRouter()

	//r.HandleFunc("/block", m.AddContext(context, NewBlock)).Methods("POST")
	r.Handle("/block", m.AddContextLedger(l, http.HandlerFunc(NewBlock))).Methods("POST")

	return r
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
