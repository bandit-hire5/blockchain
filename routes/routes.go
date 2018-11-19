package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct{}

func (ro Router) Init() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/block", NewBlock).Methods("POST")

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
