package routes

import (
	"encoding/json"
	"net/http"

	"github.com/bandit/blockchain-core"
	m "github.com/bandit/blockchain/middlewares"
)

func NewBlock(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ledger := m.Ledger(r)

	var data core.Data

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	block, err := ledger.NextBlock(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	err = ledger.AddBlock(block)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusCreated, block)
}
