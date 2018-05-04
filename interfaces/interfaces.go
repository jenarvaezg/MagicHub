package interfaces

import "github.com/graphql-go/graphql"

// GraphQLController is an interface for Controllers that return mutations and queries
type GraphQLController interface {
	GetMutations() graphql.Fields
	GetQueries() graphql.Fields
}
