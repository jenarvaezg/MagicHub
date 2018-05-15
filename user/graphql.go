package user

import (
	"github.com/graphql-go/graphql"

	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
)

type controller struct {
	service interfaces.UserService
	types   map[string]graphql.Output
	fields  map[string]*graphql.Field
}

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(service interfaces.UserService, r interfaces.Registry) interfaces.GraphQLController {
	c := &controller{service: service}
	c.types = make(map[string]graphql.Output)
	c.fields = make(map[string]*graphql.Field)

	c.types["user"] = userType

	r.RegisterController(c, "user")
	return c
}

func (c *controller) GetQueries() graphql.Fields {
	return graphql.Fields{"user": c.getUserQuery()}
}

func (c *controller) GetMutations() graphql.Fields {
	return graphql.Fields{}
}

func (c controller) GetOutputType(name string) graphql.Output {
	return c.types[name]
}

func (c controller) GetField(name string) *graphql.Field {
	return nil
}

func (c *controller) OnAllControllersRegistered(sr interfaces.Registry) {}

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "A user is an user of MagicHub",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			user := p.Source.(*models.User)
			return user.GetId().Hex(), nil
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
		Type: c.GetOutputType("user"),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID), Description: "User ID"},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			RequireUser(p.Context)
			id := p.Args["id"].(string)
			return c.service.FindByID(id)
		},
	}
}
