package interfaces

import "github.com/graphql-go/graphql"

// Controller is an interface for Controllers
type Controller interface {
	OnAllControllersRegistered(r Registry)
}

// GraphQLController is an interface for Controllers that return mutations and queries
type GraphQLController interface {
	Controller
	GetMutations() graphql.Fields
	GetQueries() graphql.Fields
	GetOutputType(name string) graphql.Output
	GetField(name string) *graphql.Field
}
