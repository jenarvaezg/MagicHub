package user

import (
	"fmt"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *bongo.Collection

// User is a document which holds information about a user
type User struct {
	ID        bson.ObjectId `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	ImageURL  string        `json:"imageUrl"`
}

// List is a list of User Documents
type List []User

func (u *User) String() string {
	return fmt.Sprintf("User: %q id %s email %q", u.Username, u.ID, u.Email)
}

func newUserList() List {
	return make([]User, 0)
}

//ListUsers returns all boxes in the box collection
func ListUsers() (users List) {
	users = newUserList()
	results := userCollection.Find(bson.M{})

	user := User{}
	for results.Next(&user) {
		users = append(users, user)
	}

	return
}
