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

// FindFiltered returns a list of pointer to boxes from mongodb filtered by limit and offset
func (r *repo) FindFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error) {
	model := r.getModel()
	query := utils.QueryLimitAndOffset(limit, offset, model.Find())

	var boxes []*models.Box
	err := query.Exec(&boxes)

	return boxes, err
}

func (r repo) getModel() *mongodm.Model {
	return r.connection.Model("Box")
}
