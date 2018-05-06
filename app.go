package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jenarvaezg/MagicHub/auth"
	"github.com/jenarvaezg/MagicHub/handlers"
	"github.com/jenarvaezg/MagicHub/middleware"
	"github.com/jenarvaezg/MagicHub/team"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/jenarvaezg/MagicHub/utils"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

const (
	baseRoute  string = "/api/v1"
	boxRoute   string = "/box"
	notesRoute string = "/notes"
	register   string = "/register"
	idRoute    string = "/{id:[0-9a-f]+}"
	port       string = "8000"
)

var apiCommonMiddleware *negroni.Negroni

func getAPICommonMiddleware() *negroni.Negroni {
	optionsMiddleware := cors.AllowAll()
	return negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		optionsMiddleware,
		middleware.NewUserFromJWTMiddleware(),
	)
}

func getGraphQLSchema() *graphql.Schema {
	teamRepo := team.NewMongoRepository()
	teamController := team.NewGraphQLController(teamRepo, team.NewService(teamRepo))
	userRepo := user.NewMongoRepository()
	userController := user.NewGraphQLController(userRepo, user.NewService(userRepo))
	authController := auth.NewGraphQLController(auth.NewService(user.NewService(userRepo)))

	queryFields := utils.MergeGraphQLFields(
		teamController.GetQueries(),
		userController.GetQueries(),
		authController.GetQueries(),
	)
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queryFields}

	mutationFields := utils.MergeGraphQLFields(
		teamController.GetMutations(),
		userController.GetMutations(),
		authController.GetMutations(),
	)
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Panic(err)
	}

	return &schema
}

func init() {
	apiCommonMiddleware = getAPICommonMiddleware()
}

func main() {

	log.Println("Setting up routes")
	middlewareRouter := mux.NewRouter()
	router := mux.NewRouter() //two routers are neccesary due to negroni

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

	// Middlewares
	// Order matters, we have to go from most to least specific routes

	middlewareRouter.PathPrefix(baseRoute + boxRoute + idRoute).Handler(apiCommonMiddleware.With(
		middleware.NewRequireBoxMiddleware(),
		negroni.Wrap(boxDetailRouter),
	))
	middlewareRouter.PathPrefix("/").Handler(apiCommonMiddleware.With(
		negroni.Wrap(router),
	))

	graphHandler := handler.New(&handler.Config{
		Schema: getGraphQLSchema(),
		Pretty: true,
	})

	router.Handle("/graphql", graphHandler)

	log.Println("Server starting at port", port)
	log.Panic(http.ListenAndServe(":"+port, middlewareRouter))
}
