package user

import (
	"github.com/graphql-go/graphql"
	"github.com/jenarvaezg/MagicHub/interfaces"
)

type controller struct {
	repo    Repository
	service Service
}

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(repo Repository, service Service) interfaces.GraphQLController {
	return &controller{repo: repo, service: service}
}

func (c *controller) GetQueries() graphql.Fields {
	return graphql.Fields{"user": c.getUserQuery()}
}

func (c *controller) GetMutations() graphql.Fields {

	return graphql.Fields{}
}

// UType is the type that holds a user for GraphQL
var UType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A user is an user of MagicHub",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			user := p.Source.(*User)
			return user.ID.Hex(), nil
		}},
		"email":     &graphql.Field{Type: graphql.String},
		"firstName": &graphql.Field{Type: graphql.String},
		"lastName":  &graphql.Field{Type: graphql.String},
		"username":  &graphql.Field{Type: graphql.String},
		"imageUrl":  &graphql.Field{Type: graphql.String},
	},
})

func (c *controller) getUserQuery() *graphql.Field {
	return &graphql.Field{
		Type: UType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID), Description: "User ID"},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Args["id"].(string)
			return c.service.FindByID(id)
		},
	}
}
