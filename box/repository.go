package box

import (
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

// Repository is a interface of all the methods required to be a repository for Box
type Repository interface {
	FindFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error)
}
