package auth

import (
	"github.com/graphql-go/graphql"
	"github.com/jenarvaezg/MagicHub/interfaces"
)

type controller struct {
	service interfaces.AuthService
}

var loginType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "LoginResponse",
	Description: "A LoginResponse contains an auth token and the data of the user that signed in",
	Fields: graphql.Fields{
		"jwt": &graphql.Field{Type: graphql.String, Description: "JWT token that can be used for auth"},
	},
})

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(service interfaces.AuthService, r interfaces.Registry) interfaces.GraphQLController {
	c := &controller{service: service}
	r.RegisterController(c, "auth")

	return c
}

func (c *controller) GetQueries() graphql.Fields {
	return graphql.Fields{}
}

func (c *controller) GetMutations() graphql.Fields {

	return graphql.Fields{"login": c.loginMutation()}
}

func (c *controller) OnAllControllersRegistered(r interfaces.Registry) {
	userController := r.GetController("user").(interfaces.GraphQLController)

	loginType.AddFieldConfig("user", &graphql.Field{
		Type: userController.GetOutputType("user"), Description: "User that logged in",
	})
}

func (c *controller) GetOutputType(name string) graphql.Output {
	return loginType
}

func (c *controller) GetField(name string) *graphql.Field {
	return nil
}

func (c *controller) loginMutation() *graphql.Field {
	return &graphql.Field{
		Type: loginType,
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
