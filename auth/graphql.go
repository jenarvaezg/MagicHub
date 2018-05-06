package auth

import (
	"github.com/graphql-go/graphql"
	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/user"
)

type controller struct {
	//repo    Repository
	service Service
}

var loginResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "LoginResponse",
	Description: "A LoginResponse contains an auth token and the data of the user that signed in",
	Fields: graphql.Fields{
		"jwt":  &graphql.Field{Type: graphql.String, Description: "JWT token that can be used for auth"},
		"user": &graphql.Field{Type: user.UType, Description: "User that logged in"},
	},
})

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(service Service) interfaces.GraphQLController {
	return &controller{service: service}
}

func (c *controller) GetQueries() graphql.Fields {
	return graphql.Fields{}
}

func (c *controller) GetMutations() graphql.Fields {

	return graphql.Fields{"login": c.loginMutation()}
}

func (c *controller) loginMutation() *graphql.Field {
	return &graphql.Field{
		Type: loginResponseType,
		Args: graphql.FieldConfigArgument{
			"token":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String), Description: "Auth provider token"},
			"provider": &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "google", Description: "Provider for the token"},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			token := p.Args["token"].(string)
			provider := p.Args["provider"].(string)
			return c.service.GetAuthTokenByProvider(token, provider)
		},
	}
}
