package main

import (
	"log"
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi"
)

func main() {
	h := httpapi.NewHandler()
	mux := httpapi.NewRouter(h)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
