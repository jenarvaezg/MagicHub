package team

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/jenarvaezg/MagicHub/utils"
	"gopkg.in/mgo.v2/bson"
)

type listResult struct {
	Nodes      []*models.Team `json:"nodes"`
	TotalCount int            `json:"totalCount"`
}

type controller struct {
	service interfaces.TeamService
	types   map[string]graphql.Output
	fields  map[string]*graphql.Field
}

var teamType *graphql.Object
var teamListQuery *graphql.Field
var teamByIDQuery *graphql.Field

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(service interfaces.TeamService, r interfaces.Registry) interfaces.GraphQLController {
	c := &controller{service: service}
	c.types = make(map[string]graphql.Output)
	c.fields = make(map[string]*graphql.Field)
	c.setTypes()

	c.types["team"] = teamType
	c.fields["teamList"] = teamListQuery
	c.fields["teamByID"] = teamByIDQuery

	r.RegisterController(c, "team")
	return c
}

func (c *controller) GetQueries() graphql.Fields {
	return graphql.Fields{"teams": teamListQuery, "team": teamByIDQuery}
}

func (c *controller) GetMutations() graphql.Fields {
	var createTeamMutation = &graphql.Field{
		Type: teamType,
		Args: graphql.FieldConfigArgument{
			"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String), Description: "Team name"},
			"image":       &graphql.ArgumentConfig{Type: graphql.String, Description: "Path to image of team"},
			"description": &graphql.ArgumentConfig{Type: graphql.String, Description: "Short description of team"},
		},
		Resolve: c.createTeamResolver,
	}

	return graphql.Fields{"createTeam": createTeamMutation}
}

func (c *controller) GetField(name string) *graphql.Field {
	return c.fields[name]
}

func (c *controller) GetOutputType(name string) graphql.Output {
	return c.types[name]
}

func (c *controller) OnAllControllersRegistered(r interfaces.Registry) {
	userController := r.GetController("user").(interfaces.GraphQLController)
	boxController := r.GetController("box").(interfaces.GraphQLController)

	userType := userController.GetOutputType("user")
	boxListField := boxController.GetField("boxList")

	teamType.AddFieldConfig("members", &graphql.Field{
		Type: graphql.NewList(userType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := user.RequireUser(p.Context)
			team := p.Source.(*models.Team)
			return c.service.GetTeamMembers(userID, team)
		},
	})

	teamType.AddFieldConfig("admins", &graphql.Field{
		Type: graphql.NewList(userType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := user.RequireUser(p.Context)
			team := p.Source.(*models.Team)
			return c.service.GetTeamAdmins(userID, team)
		},
	})

	teamType.AddFieldConfig("boxes", boxListField)

}

func (c *controller) setTypes() {
	teamType = c.teamType()
	teamByIDQuery = c.getTeamQuery()
	teamListQuery = utils.MakeListField(utils.MakeNodeListType("TeamList", teamType), c.teamListQuery, true)
}

func (c *controller) teamListQuery(p graphql.ResolveParams) (interface{}, error) {
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

			if !bson.IsObjectIdHex(id) {
				return nil, fmt.Errorf("%v is not a valid object id", id)
			}
			return c.service.FindByID(bson.ObjectIdHex(id))
		},
	}
}

func (c *controller) createTeamResolver(p graphql.ResolveParams) (interface{}, error) {
	userID := user.RequireUser(p.Context)
	name, _ := p.Args["name"].(string)
	image, _ := p.Args["image"].(string)
	description, _ := p.Args["description"].(string)

	return c.service.CreateTeam(userID, name, image, description)
}

func (c *controller) teamType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "Team",
		Description: "A team is the base organizational level, holds stuff like users, boxes, etc...",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.String, Resolve: utils.GetIDResolver},
			"name":        &graphql.Field{Type: graphql.String},
			"image":       &graphql.Field{Type: graphql.String},
			"routeName":   &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"memberCount": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				team := p.Source.(*models.Team)
				return c.service.GetTeamMembersCount(team)
			}},
		},
	})
}
