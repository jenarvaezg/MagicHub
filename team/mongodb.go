package team

import (
	"github.com/jenarvaezg/MagicHub/db"
	"github.com/jenarvaezg/MagicHub/utils"
	"github.com/zebresel-com/mongodm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const collectionName = "team"

type repo struct {
	connection *mongodm.Connection
}

type teamDocument struct {
	mongodm.DocumentBase `bson:",inline"`

	Image       string `json:"image"`
	Name        string `json:"name"`
	RouteName   string `json:"routeName"`
	Description string `json:"description"`
}

// NewMongoRepository returns a object that implements the repository interface using mongodb
func NewMongoRepository() Repository {
	connection := db.GetMongoConnection()
	connection.Register(&teamDocument{}, collectionName)
	index := mgo.Index{
		Key: []string{"$text:name", "$text:description"},
	}
	connection.Session.DB(db.DatabaseName).C(collectionName).EnsureIndex(index)

	return &repo{connection}
}

// Store saves a team to mongodb and returns a pointer to the team with updated fields
func (r *repo) Store(t *Team) (bson.ObjectId, error) {
	teamDoc := &teamDocument{Image: t.Image, Name: t.Name, RouteName: t.RouteName, Description: t.Description}
	model := r.getModel()

	model.New(teamDoc)
	if err := teamDoc.Save(); err != nil {
		return bson.NewObjectId(), err
	}

	t.ID = teamDoc.Id

	return t.ID, nil
}

// FindFiltered returns a list of pointer to teams from mongodb filtered by limit offset and search parameter
func (r *repo) FindFiltered(limit, offset int, search string) ([]*Team, error) {
	model := r.getModel()
	var query *mongodm.Query

	if search != "" {
		regex := bson.RegEx{Pattern: search, Options: "i"}
		query = model.Find(bson.M{"$or": []bson.M{bson.M{"name": regex}, bson.M{"description": regex}}})
	} else {
		query = model.Find()
	}
	query = utils.QueryLimitAndOffset(limit, offset, query)

	var teams []*teamDocument
	var teamInstances []*Team

	err := query.Exec(&teams)

	for _, teamP := range teams {
		teamInstances = append(teamInstances, teamP.instanceFromModel())
	}

	// Run the query
	return teamInstances, err
}

// FindByID returns a matching team by ID or error if not found
func (r *repo) FindByID(id bson.ObjectId) (*Team, error) {
	model := r.getModel()
	teamDoc := teamDocument{}

	if err := model.FindId(id).Exec(&teamDoc); err != nil {
		return nil, err
	}

	return teamDoc.instanceFromModel(), nil
}

func (t teamDocument) instanceFromModel() *Team {
	return &Team{
		ID:          t.Id,
		RouteName:   t.RouteName,
		Name:        t.Name,
		Image:       t.Image,
		Description: t.Description,
	}
}

func (r repo) getModel() *mongodm.Model {
	return r.connection.Model("teamDocument")
}
