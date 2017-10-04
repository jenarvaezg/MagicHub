package main

import (
	"log"
	"net/http"

	"github.com/jenarvaezg/magicbox/handlers"

	"github.com/gorilla/mux"
)

const (
	baseRoute string = "/api/v1"
	port      string = "8000"
)

func main() {
	r := mux.NewRouter()
	api_router := r.PathPrefix(baseRoute).Subrouter()
	api_router.HandleFunc("/box", handlers.ListBoxesHandler).Methods("GET")
	api_router.HandleFunc("/box", handlers.CreateBoxHandler).Methods("POST")

	log.Panic(http.ListenAndServe(":"+port, api_router))
}
