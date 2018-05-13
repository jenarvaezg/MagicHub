package team

import (
	"github.com/graphql-go/graphql"
	"github.com/jenarvaezg/MagicHub/box"
	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/jenarvaezg/MagicHub/utils"
)

type listResult struct {
	Nodes      []*models.Team `json:"nodes"`
	TotalCount int            `json:"totalCount"`
}

type controller struct {
	repo       Repository
	service    Service
	boxService box.Service
}

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(repo Repository, service Service, boxService box.Service) interfaces.GraphQLController {
	c := &controller{repo: repo, service: service, boxService: boxService}
	c.setTeamType()
	return c
}

func (c *controller) GetQueries() graphql.Fields {
	teamsQuery := utils.MakeListField(utils.MakeNodeListType("TeamList", teamType), c.queryTeams, true)
	teamByIDQuery := c.getTeamQuery()
	return graphql.Fields{"teams": teamsQuery, "team": teamByIDQuery}
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

func (c *controller) queryTeams(p graphql.ResolveParams) (interface{}, error) {
	user.RequireUser(p.Context)
	limit, _ := p.Args["limit"].(int)
	offset, _ := p.Args["offset"].(int)
	search, _ := p.Args["search"].(string)

	var result listResult
	var err error

	result.Nodes, err = c.service.FindFiltered(limit, offset, search)
	result.TotalCount = len(result.Nodes)
	return result, err
}

func (c *controller) getTeamQuery() *graphql.Field {
	return &graphql.Field{
		Type: teamType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID), Description: "Team identifier"},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, _ := p.Args["id"].(string)

			return c.service.GetTeamByID(id)
		},
	}
}

func (c *controller) createTeam(p graphql.ResolveParams) (interface{}, error) {
	userID := user.RequireUser(p.Context)
	name, _ := p.Args["name"].(string)
	image, _ := p.Args["image"].(string)
	description, _ := p.Args["description"].(string)

	return c.service.CreateTeam(userID, name, image, description)
}

var teamType *graphql.Object

func (c *controller) setTeamType() {
	teamType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Team",
		Description: "A team is the base organizational level, holds stuff like users, boxes, etc...",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				team := p.Source.(*models.Team)
				return team.GetId().Hex(), nil
			}},
			"name":        &graphql.Field{Type: graphql.String},
			"image":       &graphql.Field{Type: graphql.String},
			"routeName":   &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"memberCount": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				team := p.Source.(*models.Team)
				return c.service.GetTeamMembersCount(team)
			}},
			"members": &graphql.Field{Type: graphql.NewList(user.UType), Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := user.RequireUser(p.Context)
				team := p.Source.(*models.Team)
				return c.service.GetTeamMembers(userID, team)
			}},
			"admins": &graphql.Field{Type: graphql.NewList(user.UType), Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := user.RequireUser(p.Context)
				team := p.Source.(*models.Team)
				return c.service.GetTeamAdmins(userID, team)
			}},
			"boxes": box.BListField,
		},
	})
}
