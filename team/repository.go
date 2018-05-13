package team

import (
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

// Repository is an interface that contains all required methods that fetch data
type Repository interface {
	FindByID(id bson.ObjectId) (*models.Team, error)
	FindFiltered(limit, offset int, search string) ([]*models.Team, error)
	Store(team *models.Team) (bson.ObjectId, error)
}
