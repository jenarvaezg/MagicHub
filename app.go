package main

import (
	"log"
	"net/http"

	"github.com/jenarvaezg/magicbox/handlers"
	"github.com/jenarvaezg/magicbox/middleware"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

const (
	baseRoute string = "/api/v1"
	boxRoute  string = "/box"
	idRoute   string = "/{id:[0-9a-f]+}"
	port      string = "8000"
)

func main() {
	middlewareRouter := mux.NewRouter()
	router := mux.NewRouter() //two routers are neccesary due to negroni
	apiRouter := router.PathPrefix(baseRoute).Subrouter()
	// Box router
	boxRouter := apiRouter.PathPrefix(boxRoute).Subrouter()
	boxRouter.HandleFunc("", handlers.ListBoxesHandler).Methods("GET")
	boxRouter.HandleFunc("", handlers.CreateBoxHandler).Methods("POST")
	//Box detail routes
	boxDetailRouter := boxRouter.PathPrefix(idRoute).Subrouter()
	boxDetailRouter.HandleFunc("", handlers.BoxDetailHandler).Methods("GET")
	boxDetailRouter.HandleFunc("", handlers.BoxDeleteHandler).Methods("DELETE")
	boxDetailRouter.HandleFunc("", handlers.BoxPatchHandler).Methods("PATCH")
	// Note routes
	noteRouter := boxRouter.PathPrefix("/notes").Subrouter()
	noteRouter.HandleFunc("", handlers.ListNotesHandler).Methods("GET")
	noteRouter.HandleFunc("", handlers.InsertNoteHandler).Methods("PUT")

	// Middlewares
	apiCommonMiddleware := negroni.New(
		middleware.NewRequireJSONMiddleware(),
	)

	// Order matters, we have to go from most to least specific routes
	middlewareRouter.PathPrefix(baseRoute + boxRoute + idRoute).Handler(apiCommonMiddleware.With(
		middleware.NewRequireBoxMiddleware(),
		negroni.Wrap(boxDetailRouter),
	))
	middlewareRouter.PathPrefix(baseRoute).Handler(apiCommonMiddleware.With(
		negroni.Wrap(boxRouter),
	))

	log.Panic(http.ListenAndServe(":"+port, middlewareRouter))
}
