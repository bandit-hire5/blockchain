package routes

import (
	"encoding/json"
	l "naeltok/go-blockchain/ledger"
	m "naeltok/go-blockchain/middlewares"
	"net/http"
)

func NewBlock(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ledger := m.Ledger(r)

	var data l.Data

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	block, err := ledger.GenerateNextBlock(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	err = ledger.WriteLedgerFile(block)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusCreated, block)
}
