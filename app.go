package main

import (
	"log"
	"net/http"

	"github.com/jenarvaezg/magicbox/handlers"
	"github.com/jenarvaezg/magicbox/middleware"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

const (
	baseRoute  string = "/api/v1"
	boxRoute   string = "/box"
	notesRoute string = "/notes"
	userRoute  string = "/user"
	loginRoute string = "/login"
	register   string = "/register"
	idRoute    string = "/{id:[0-9a-f]+}"
	port       string = "8000"
)

var apiCommonMiddleware *negroni.Negroni

func getAPICommonMiddleware() *negroni.Negroni {
	optionsMiddleware := cors.AllowAll()
	return negroni.New(
		negroni.NewLogger(),
		optionsMiddleware,
		middleware.NewRequireJSONMiddleware(),
		middleware.NewUserFromJWTMiddleware(),
	)
}

func init() {
	apiCommonMiddleware = getAPICommonMiddleware()
}

func main() {
	log.Println("Setting up routes")
	middlewareRouter := mux.NewRouter()
	router := mux.NewRouter() //two routers are neccesary due to negroni
	// Token routes
	tokenRouter := router.PathPrefix(loginRoute).Subrouter()
	tokenRouter.HandleFunc("", handlers.LoginRequestHandler).Methods("POST")

	// API routes
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
	//Box register routes
	boxRegisterRouter := boxDetailRouter.PathPrefix(register).Subrouter()
	boxRegisterRouter.HandleFunc("", handlers.RegisterInBoxHandler).Methods("POST")
	boxRegisterRouter.HandleFunc("", handlers.RemoveFromBoxHandler).Methods("DELETE")
	// Note routes
	noteRouter := boxDetailRouter.PathPrefix(notesRoute).Subrouter()
	noteRouter.HandleFunc("", handlers.ListNotesHandler).Methods("GET")
	noteRouter.HandleFunc("", handlers.InsertNoteHandler).Methods("POST")
	noteRouter.HandleFunc("", handlers.DeleteNotesHandler).Methods("DELETE")
	// User routes
	userRouter := apiRouter.PathPrefix(userRoute).Subrouter()
	userRouter.HandleFunc("", handlers.ListUsersHandler).Methods("GET")
	userRouter.HandleFunc("", handlers.CreateUserHandler).Methods("POST").Name("create-user-url")
	// User detail routes
	userDetailRouter := userRouter.PathPrefix(idRoute).Subrouter()
	userDetailRouter.HandleFunc("", handlers.UserDetailHandler).Methods("GET")
	userDetailRouter.HandleFunc("", handlers.UserDeleteHandler).Methods("DELETE")
	userDetailRouter.HandleFunc("", handlers.UserPatchHandler).Methods("PATCH")
	// Middlewares
	// Order matters, we have to go from most to least specific routes

	middlewareRouter.PathPrefix(baseRoute + userRoute + idRoute).Handler(apiCommonMiddleware.With(
		middleware.NewRequireUserMiddleware(),
		negroni.Wrap(userDetailRouter),
	))
	middlewareRouter.PathPrefix(baseRoute + boxRoute + idRoute).Handler(apiCommonMiddleware.With(
		middleware.NewRequireBoxMiddleware(),
		negroni.Wrap(boxDetailRouter),
	))
	middlewareRouter.PathPrefix(baseRoute).Handler(apiCommonMiddleware.With(
		negroni.Wrap(apiRouter),
	))
	middlewareRouter.PathPrefix(loginRoute).Handler(negroni.New(
		negroni.NewLogger(),
		cors.AllowAll(),
		negroni.Wrap(tokenRouter),
	))

	log.Println("Server starting at port", port)
	log.Panic(http.ListenAndServe(":"+port, middlewareRouter))
}
