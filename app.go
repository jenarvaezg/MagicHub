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
	baseRoute  string = "/api/v1"
	boxRoute   string = "/box"
	notesRoute string = "/notes"
	userRoute  string = "/user"
	idRoute    string = "/{id:[0-9a-f]+}"
	port       string = "8000"
)

var apiCommonMiddleware *negroni.Negroni

func getAPICommonMiddleware() *negroni.Negroni {
	return negroni.New(
		negroni.NewLogger(),
		middleware.NewRequireJSONMiddleware(),
	)
}

func init() {
	apiCommonMiddleware = getAPICommonMiddleware()
}

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
	noteRouter := boxDetailRouter.PathPrefix(notesRoute).Subrouter()
	noteRouter.HandleFunc("", handlers.ListNotesHandler).Methods("GET")
	noteRouter.HandleFunc("", handlers.InsertNoteHandler).Methods("POST")
	noteRouter.HandleFunc("", handlers.DeleteNotesHandler).Methods("DELETE")
	// User routes
	userRouter := apiRouter.PathPrefix(userRoute).Subrouter()
	userRouter.HandleFunc("", handlers.ListUsersHandler).Methods("GET")
	userRouter.HandleFunc("", handlers.CreateUserHandler).Methods("POST")

	// Middlewares
	// Order matters, we have to go from most to least specific routes
	middlewareRouter.PathPrefix(baseRoute + boxRoute + idRoute).Handler(apiCommonMiddleware.With(
		middleware.NewRequireBoxMiddleware(),
		negroni.Wrap(boxDetailRouter),
	))
	middlewareRouter.PathPrefix(baseRoute).Handler(apiCommonMiddleware.With(
		negroni.Wrap(apiRouter),
	))

	log.Panic(http.ListenAndServe(":"+port, middlewareRouter))
}
