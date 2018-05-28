package user

import (
	"github.com/jenarvaezg/mongodm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/jenarvaezg/MagicHub/db"
	"github.com/jenarvaezg/MagicHub/models"
)

const collectionName = "user"

type repo struct {
	connection *mongodm.Connection
}

// NewMongoRepository returns a object that implements the repository interface using mongodb
func NewMongoRepository() Repository {
	connection := db.GetMongoConnection()
	connection.Register(&models.User{}, collectionName)
	index := mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}
	connection.Session.DB(db.DatabaseName).C(collectionName).EnsureIndex(index)

	return &repo{connection}
}

// Store saves a user to mongodb and returns a pointer to the team with updated fields
func (r *repo) Store(u *models.User) (bson.ObjectId, error) {
	model := r.getModel()

	model.New(u)
	if err := u.Save(); err != nil {
		return bson.NewObjectId(), err
	}

	return u.GetId(), nil
}

// FindByID returns a matching team by ID or error if not found
func (r *repo) FindByID(id bson.ObjectId) (*models.User, error) {
	model := r.getModel()
	user := &models.User{}

	if err := model.FindId(id).Exec(user); err != nil {
		return nil, err
	}

	return user, nil
}

// FindBy returns a list of Users y the provided fields or error if there is a problem
func (r *repo) FindBy(findMap map[string]interface{}) ([]*models.User, error) {
	model := r.getModel()
	users := []*models.User{}

	err := model.Find(findMap).Exec(&users)

	return users, err
}

func (r repo) getModel() *mongodm.Model {
	return r.connection.Model("User")
}
