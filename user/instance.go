package user

import (
	"fmt"

	"github.com/go-bongo/bongo"
	"github.com/zebresel-com/mongodm"
)

var userCollection *bongo.Collection

// User is a document which holds information about a user
type User struct {
	mongodm.DocumentBase `bson:",inline"`

	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	ImageURL  string `json:"imageUrl" bson:"imageUrl"`
}

func (u *User) String() string {
	return fmt.Sprintf("User: %q, id: %s, email %q", u.Username, u.GetId(), u.Email)
}
