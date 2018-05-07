package team

import (
	"github.com/graphql-go/graphql"
	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/utils"
)

type listResult struct {
	Nodes      []*Team `json:"nodes"`
	TotalCount int     `json:"totalCount"`
}

type controller struct {
	repo    Repository
	service Service
}

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(repo Repository, service Service) interfaces.GraphQLController {
	return &controller{repo: repo, service: service}
}

func (c *controller) GetQueries() graphql.Fields {
	teamsQuery := utils.MakeListField(utils.MakeNodeListType("TeamList", teamType), c.queryTeams, true)
	return graphql.Fields{"teams": teamsQuery}
}

func (c *controller) GetMutations() graphql.Fields {
	var createTeamMutation = &graphql.Field{
		Type: teamType,
		Args: graphql.FieldConfigArgument{
			"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String), Description: "Team name"},
			"image":       &graphql.ArgumentConfig{Type: graphql.String, Description: "Path to image of team"},
			"description": &graphql.ArgumentConfig{Type: graphql.String, Description: "Short description of team"},
		},
		Resolve: c.createTeam,
	}

	return graphql.Fields{"createTeam": createTeamMutation}
}

var teamType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Team",
	Description: "A team is the base organizational level, holds stuff like users, boxes, etc...",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			team := p.Source.(*Team)
			return team.ID.Hex(), nil
		}},
		"name":        &graphql.Field{Type: graphql.String},
		"image":       &graphql.Field{Type: graphql.String},
		"routeName":   &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
	},
})

func (c *controller) queryTeams(p graphql.ResolveParams) (interface{}, error) {
	utils.RequireUser(p.Context)
	limit, _ := p.Args["limit"].(int)
	offset, _ := p.Args["offset"].(int)
	search, _ := p.Args["search"].(string)

	var result listResult
	var err error

	result.Nodes, err = c.service.FindFiltered(limit, offset, search)
	result.TotalCount = len(result.Nodes)
	return result, err
}

func (c *controller) createTeam(p graphql.ResolveParams) (interface{}, error) {
	utils.RequireUser(p.Context)
	name, _ := p.Args["name"].(string)
	image, _ := p.Args["image"].(string)
	description, _ := p.Args["description"].(string)

	return c.service.CreateTeam(name, image, description)
}
