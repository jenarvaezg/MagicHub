package user

import (
	"github.com/zebresel-com/mongodm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/jenarvaezg/MagicHub/db"
)

const collectionName = "user"

type repo struct {
	connection *mongodm.Connection
}

type userDocument struct {
	mongodm.DocumentBase `bson:",inline"`

	Username  string `bson:"username"`
	Email     string `bson:"email"`
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
	ImageURL  string `bson:"image_url"`
}

// NewMongoRepository returns a object that implements the repository interface using mongodb
func NewMongoRepository() Repository {
	connection := db.GetMongoConnection()
	connection.Register(&userDocument{}, collectionName)
	index := mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}
	connection.Session.DB(db.DatabaseName).C(collectionName).EnsureIndex(index)

	return &repo{connection}
}

// Store saves a user to mongodb and returns a pointer to the team with updated fields
func (r *repo) Store(u *User) (bson.ObjectId, error) {
	userDoc := &userDocument{Username: u.Username, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName, ImageURL: u.ImageURL}
	model := r.getModel()

	model.New(userDoc)
	if err := userDoc.Save(); err != nil {
		return bson.NewObjectId(), err
	}

	u.ID = userDoc.Id

	return u.ID, nil
}

// FindByID returns a matching team by ID or error if not found
func (r *repo) FindByID(id bson.ObjectId) (*User, error) {
	model := r.getModel()
	userDoc := userDocument{}

	if err := model.FindId(id).Exec(&userDoc); err != nil {
		return nil, err
	}

	return userDoc.instanceFromModel(), nil
}

// FindBy returns a list of Users y the provided fields or error if there is a problem
func (r *repo) FindBy(findMap map[string]interface{}) ([]*User, error) {
	model := r.getModel()
	userDocs := []*userDocument{}

	err := model.Find(findMap).Exec(&userDocs)

	users := []*User{}
	for _, userDoc := range userDocs {
		users = append(users, userDoc.instanceFromModel())
	}

	return users, err
}

func (u userDocument) instanceFromModel() *User {
	return &User{
		ID:        u.Id,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ImageURL:  u.ImageURL,
		Email:     u.Email,
	}
}

func (r repo) getModel() *mongodm.Model {
	return r.connection.Model("userDocument")
}
