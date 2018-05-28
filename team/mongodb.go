package team

import (
	"fmt"

	"github.com/jenarvaezg/MagicHub/db"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/utils"
	"github.com/jenarvaezg/mongodm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const collectionName = "team"

type repo struct {
	connection *mongodm.Connection
}

// NewMongoRepository returns a object that implements the repository interface using mongodb
func NewMongoRepository() Repository {
	connection := db.GetMongoConnection()
	connection.Register(&models.Team{}, collectionName)
	index := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	}
	connection.Session.DB(db.DatabaseName).C(collectionName).EnsureIndex(index)

	return &repo{connection}
}

// Store saves a team.team to mongodb and returns a pointer to the team with updated fields
func (r *repo) Store(t *models.Team) (bson.ObjectId, error) {
	model := r.getModel()
	model.New(t)

	if err := t.Save(); err != nil {
		if _, ok := err.(*mongodm.DuplicateError); ok {
			err = fmt.Errorf("There is already a team with the provided name: %q", t.Name)
		}
		return bson.NewObjectId(), err

	}

	return t.GetId(), t.Populate("Admins", "Members")
}

// FindFiltered returns a list of pointer to teams from mongodb filtered by limit offset and search parameter
func (r *repo) FindFiltered(limit, offset int, search string) ([]*models.Team, error) {
	model := r.getModel()
	var query *mongodm.Query

	if search != "" {
		regex := bson.RegEx{Pattern: search, Options: "i"}
		query = model.Find(bson.M{"$or": []bson.M{bson.M{"name": regex}, bson.M{"description": regex}}})
	} else {
		query = model.Find()
	}
	query = utils.QueryLimitAndOffset(limit, offset, query)

	var teams []*models.Team

	// Run the query
	err := query.Populate("Admins", "Members").Exec(&teams)

	return teams, err
}

// FindByID returns a matching team by ID or error if not found
func (r *repo) FindByID(id bson.ObjectId) (*models.Team, error) {
	model := r.getModel()
	team := &models.Team{}

	if err := model.FindId(id).Populate("Admins", "Members").Exec(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (r repo) getModel() *mongodm.Model {
	return r.connection.Model("Team")
}
