package routes

import (
	"encoding/json"
	b "naeltok/go-blockchain/block"
	"net/http"
)

func NewBlock(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data b.Data

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	block, err := b.GenerateNextBlock(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	err = b.WriteNewBlock(block)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusCreated, block)
}
