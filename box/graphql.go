package box

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"gopkg.in/mgo.v2/bson"

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
var noteType *graphql.Object

// NewGraphQLController returns a GraphQLController
func NewGraphQLController(service interfaces.BoxService, r interfaces.Registry) interfaces.GraphQLController {
	c := &controller{service: service}
	c.setFields()

	c.types = make(map[string]graphql.Output)
	c.fields = make(map[string]*graphql.Field)

	c.types["box"] = boxType
	c.types["note"] = noteType
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
			"openDate": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.DateTime), Description: "Time when the box will open"},
		},
		Resolve: c.createBox,
	}

	var insertNoteMutation = &graphql.Field{
		Type: c.GetOutputType("box"),
		Args: graphql.FieldConfigArgument{
			"boxID": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID), Description: "ID of the box where the note will be inserted"},
			"text":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String), Description: "Text of the note"},
		},
		Resolve: c.insertNote,
	}

	return graphql.Fields{"createBox": createBoxMutation, "insertNote": insertNoteMutation}
}

func (c controller) GetOutputType(name string) graphql.Output {
	return c.types[name]
}

func (c controller) GetField(name string) *graphql.Field {
	return c.fields[name]
}

func (c *controller) OnAllControllersRegistered(r interfaces.Registry) {
	teamController := r.GetController("team").(interfaces.GraphQLController)
	userController := r.GetController("user").(interfaces.GraphQLController)

	boxType.AddFieldConfig("team", &graphql.Field{
		Type: teamController.GetOutputType("team"), Description: "Box's team",
	})

	noteType.AddFieldConfig("from", &graphql.Field{
		Type: userController.GetOutputType("user"), Description: "User who wrote the box", Resolve: c.userFromNoteResolver,
	})
}

func (c *controller) setFields() {
	noteType = c.noteType()
	boxType = c.boxType()
	boxListQuery = utils.MakeListField(utils.MakeNodeListType("BoxList", boxType), c.queryBoxesFromProject, true)
}

func (c *controller) boxType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "Box",
		Description: "A box is a magicbox of MagicHub",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.ID, Resolve: utils.GetIDResolver},
			"name":     &graphql.Field{Type: graphql.String},
			"openDate": &graphql.Field{Type: graphql.DateTime},
			"isOpen":   &graphql.Field{Type: graphql.Boolean, Resolve: c.isBoxOpenResolver},
			"notes":    &graphql.Field{Type: &graphql.List{OfType: noteType}, Resolve: c.getNotesResolver},
		},
	})
}

func (c *controller) noteType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "Note",
		Description: "A note is the content of a box",
		Fields: graphql.Fields{
			// "from": get user on all registered
			"text": &graphql.Field{Type: graphql.String},
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

	result.Nodes, err = c.service.FindByTeamFiltered(limit, offset, team.GetId())
	result.TotalCount = len(result.Nodes)
	return result, err
}

func (c *controller) createBox(p graphql.ResolveParams) (interface{}, error) {
	userID := user.RequireUser(p.Context)
	name, _ := p.Args["name"].(string)
	teamID, _ := p.Args["teamID"].(string)
	openDateStr, _ := p.Args["openDate"].(string)

	openDate, err := time.Parse(time.RFC3339, openDateStr)
	if err != nil {
		return nil, fmt.Errorf("date is expected in format 2018-05-31T08:00:00.000Z")
	}

	if !bson.IsObjectIdHex(teamID) {
		return nil, fmt.Errorf("%v is not a valid object id", teamID)
	}

	return c.service.CreateBox(userID, bson.ObjectIdHex(teamID), name, openDate)
}

func (c *controller) insertNote(p graphql.ResolveParams) (interface{}, error) {
	userID := user.RequireUser(p.Context)
	boxID := p.Args["boxID"].(string)
	text := p.Args["text"].(string)

	if !bson.IsObjectIdHex(boxID) {
		return nil, fmt.Errorf("%v is not a valid object id", boxID)
	}

	return c.service.InsertNote(userID, bson.ObjectIdHex(boxID), text)
}

func (c *controller) userFromNoteResolver(p graphql.ResolveParams) (interface{}, error) {
	note := p.Source.(*models.Note)
	return note.From, nil
}

func (c *controller) isBoxOpenResolver(p graphql.ResolveParams) (interface{}, error) {
	box := p.Source.(*models.Box)
	return box.IsOpen(), nil
}

func (c *controller) getNotesResolver(p graphql.ResolveParams) (interface{}, error) {
	userID := user.RequireUser(p.Context)
	box := p.Source.(*models.Box)
	return c.service.GetNotes(userID, box)
}
