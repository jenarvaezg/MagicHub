package box

import (
	"github.com/jenarvaezg/MagicHub/db"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/utils"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

const collectionName = "box"

type repo struct {
	connection *mongodm.Connection
}

// NewMongoRepository returns a object that implements the repository interface using mongodb
func NewMongoRepository() Repository {
	connection := db.GetMongoConnection()
	connection.Register(&models.Box{}, collectionName)

	return &repo{connection}
}

// Store saves a box to mongodb and returns its objectID and an error if any
func (r *repo) Store(b *models.Box) (bson.ObjectId, error) {
	model := r.getModel()
	model.New(b)

	if err := b.Save(); err != nil {
		return bson.NewObjectId(), err
	}

	return b.GetId(), b.Populate("Team")
}

// FindByTeam returns a list of pointers to boxes from mongodb filtered team, limit and offset parameters
func (r *repo) FindByTeamFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error) {
	model := r.getModel()
	query := model.Find(bson.M{"team": teamID})
	utils.QueryLimitAndOffset(limit, offset, query)

	var boxes []*models.Box
	err := query.Populate("Team").Exec(&boxes)

	return boxes, err
}

func (r repo) getModel() *mongodm.Model {
	return r.connection.Model("Box")
}
