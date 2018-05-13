package box

import (
	"github.com/graphql-go/graphql"

	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/jenarvaezg/MagicHub/utils"
)

type controller struct {
	repo    Repository
	service Service
}

type listResult struct {
	Nodes      []*models.Box `json:"nodes"`
	TotalCount int           `json:"totalCount"`
}

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(repo Repository, service Service) interfaces.GraphQLController {
	c := &controller{repo: repo, service: service}
	c.setBListField()
	return c
}

func (c *controller) GetQueries() graphql.Fields {
	return graphql.Fields{}
}

func (c *controller) GetMutations() graphql.Fields {
	var createBoxMutation = &graphql.Field{
		Type: BType,
		Args: graphql.FieldConfigArgument{
			"teamID":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID), Description: "Team ID for this box"},
			"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String), Description: "Box name"},
			"openDate": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int), Description: "Time when the box will open"},
		},
		Resolve: c.createBox,
	}

	return graphql.Fields{"createBox": createBoxMutation}
}

// BListField is a list of boxes for a given project
var BListField *graphql.Field

func (c *controller) setBListField() {
	BListField = utils.MakeListField(utils.MakeNodeListType("BoxList", BType), c.queryBoxesFromProject, true)
}

func (c *controller) queryBoxesFromProject(p graphql.ResolveParams) (interface{}, error) {
	user.RequireUser(p.Context)
	limit, _ := p.Args["limit"].(int)
	offset, _ := p.Args["offset"].(int)
	team, _ := p.Source.(*models.Team)

	var result listResult
	var err error

	result.Nodes, err = c.service.FindFiltered(limit, offset, team.GetId())
	result.TotalCount = len(result.Nodes)
	return result, err
}

func (c *controller) createBox(p graphql.ResolveParams) (interface{}, error) {
	user.RequireUser(p.Context)
	// name, _ := p.Args["name"].(string)
	// openDate, _ := p.Args["openDate"].(int)
	// teamID, _ := p.Args["teamID"].(string)

	return nil, nil
}

// BType is the type that holds a box for GraphQL
var BType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Box",
	Description: "A box is a magicbox of MagicHub",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})
