package main

import (
	"log"
	"net/http"

	"naeltok/go-blockchain/block"
	. "naeltok/go-blockchain/routes"
)

var routes = Router{}

func main() {
	r := routes.Init()
	block.Init()

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
