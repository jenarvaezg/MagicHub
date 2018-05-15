package user

import (
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

// Repository is an interface that contains all required methods that fetch data for a User
type Repository interface {
	FindByID(id bson.ObjectId) (*models.User, error)
	FindBy(findMap map[string]interface{}) ([]*models.User, error)
	Store(user *models.User) (bson.ObjectId, error)
}
