package user

import (
	"gopkg.in/mgo.v2/bson"
)

// Repository is an interface that contains all required methods that fetch data for a User
type Repository interface {
	FindByID(id bson.ObjectId) (*User, error)
	FindBy(findMap map[string]interface{}) ([]*User, error)
	Store(user *User) (bson.ObjectId, error)
}
