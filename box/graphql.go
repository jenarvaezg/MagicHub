package box

import (
	"github.com/graphql-go/graphql"

	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/jenarvaezg/MagicHub/utils"
)

type controller struct {
	service interfaces.BoxService
	types   map[string]graphql.Output
	fields  map[string]*graphql.Field
}

type listResult struct {
	Nodes      []*models.Box `json:"nodes"`
	TotalCount int           `json:"totalCount"`
}

var boxListQuery *graphql.Field
var boxType *graphql.Object

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(service interfaces.BoxService, r interfaces.Registry) interfaces.GraphQLController {
	c := &controller{service: service}
	c.setFields()

	c.types = make(map[string]graphql.Output)
	c.fields = make(map[string]*graphql.Field)

	c.types["box"] = boxType
	c.fields["boxList"] = boxListQuery

	r.RegisterController(c, "box")
	return c
}

func (c *controller) GetQueries() graphql.Fields {
	return graphql.Fields{}
}

func (c *controller) GetMutations() graphql.Fields {
	var createBoxMutation = &graphql.Field{
		Type: c.GetOutputType("box"),
		Args: graphql.FieldConfigArgument{
			"teamID":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID), Description: "Team ID for this box"},
			"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String), Description: "Box name"},
			"openDate": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int), Description: "Time when the box will open"},
		},
		Resolve: c.createBox,
	}

	return graphql.Fields{"createBox": createBoxMutation}
}

func (c controller) GetOutputType(name string) graphql.Output {
	return c.types[name]
}

func (c controller) GetField(name string) *graphql.Field {
	return c.fields[name]
}

func (c *controller) OnAllControllersRegistered(r interfaces.Registry) {
	// TODO fields stuff
}

func (c *controller) setFields() {
	boxType = c.boxType()
	boxListQuery = utils.MakeListField(utils.MakeNodeListType("BoxList", boxType), c.queryBoxesFromProject, true)
}

func (c *controller) boxType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "Box",
		Description: "A box is a magicbox of MagicHub",
		Fields: graphql.Fields{
			"name": &graphql.Field{Type: graphql.String},
		},
	})
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
