package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/jenarvaezg/MagicHub/auth"
	"github.com/jenarvaezg/MagicHub/box"
	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/middleware"
	"github.com/jenarvaezg/MagicHub/registry"
	"github.com/jenarvaezg/MagicHub/team"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/jenarvaezg/MagicHub/utils"
)

const (
	defaultPort string = "8000"
)

func getGraphQLSchema(r interfaces.Registry) *graphql.Schema {
	boxService := box.NewService(box.NewMongoRepository(), r)
	userService := user.NewService(user.NewMongoRepository(), r)
	teamService := team.NewService(team.NewMongoRepository(), r)
	authService := auth.NewService(r)
	r.AllServicesRegistered()

	boxController := box.NewGraphQLController(boxService, r)
	teamController := team.NewGraphQLController(teamService, r)
	userController := user.NewGraphQLController(userService, r)
	authController := auth.NewGraphQLController(authService, r)
	r.AllControllersRegistered()

	queryFields := utils.MergeGraphQLFields(
		teamController.GetQueries(),
		userController.GetQueries(),
		authController.GetQueries(),
		boxController.GetQueries(),
	)
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queryFields}

	mutationFields := utils.MergeGraphQLFields(
		teamController.GetMutations(),
		userController.GetMutations(),
		authController.GetMutations(),
		boxController.GetMutations(),
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

func getCommonMiddleware() *negroni.Negroni {
	return negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		cors.AllowAll(),
		middleware.NewUserFromJWTMiddleware(),
	)
}

func main() {
	mux := http.NewServeMux()
	middlewareRouter := getCommonMiddleware()
	graphHandler := handler.New(&handler.Config{
		Schema: getGraphQLSchema(registry.NewRegistry()),
		Pretty: true,
	})
	mux.Handle("/graphql", graphHandler)

	middlewareRouter.UseHandler(mux)

	log.Println("Server starting at port", defaultPort)
	log.Panic(http.ListenAndServe(":"+defaultPort, middlewareRouter))
}
