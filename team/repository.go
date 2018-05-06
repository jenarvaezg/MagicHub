package team

import (
	"gopkg.in/mgo.v2/bson"
)

// Repository is an interface that contains all required methods that fetch data
type Repository interface {
	FindByID(id bson.ObjectId) (*Team, error)
	FindFiltered(limit, offset int, search string) ([]*Team, error)
	Store(team *Team) (bson.ObjectId, error)
}
