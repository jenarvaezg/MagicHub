package models

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *bongo.Collection

// UserStatus is a string that
type userStatus string

const (
	userActive   = userStatus("ACTIVE")
	userInactive = userStatus("INACTIVE")
)

// User is a document which holds information about a user
type User struct {
	bongo.DocumentBase `bson:",inline"`
	Username           string     `json:"username"`
	Password           string     `json:"-" bson:"password"`
	Email              string     `json:"email"`
	FirstName          string     `json:"firstName"`
	LastName           string     `json:"lastName"`
	Status             userStatus `json:"-" bson:"status"`
}

// NewUser returns an User instance, with status set to inactive
func NewUser() User {
	user := User{Status: userInactive}
	return user
}

// GetUserByEmail return an user from database if the email exists
func GetUserByEmail(email string) (User, error) {
	user := User{}
	err := userCollection.FindOne(bson.M{"email": email}, &user)
	log.Println(user, err)
	return user, err

}

// GetUserByUsername return an user from database if the email exists
func GetUserByUsername(username string) (User, error) {
	user := User{}
	err := userCollection.FindOne(bson.M{"username": username}, &user)

	return user, err

}

func (u *User) String() string {
	return fmt.Sprintf("User: %q id %s email %q", u.Username, u.Id, u.Email)
}

// Save saves a User instance into database
func (u *User) Save() error {
	if err := u.validate(); err != nil {
		return err
	}
	u.Status = userActive
	return userCollection.Save(u)
}

func (u *User) validate() error {
	if err := u.validateUsername(); err != nil {
		return err
	}
	if u.FirstName == "" {
		return errors.New("Field firstName is required")
	}
	if err := u.validateEmail(); err != nil {
		return err
	}
	if err := u.validatePassword(); err != nil {
		return err
	}
	return nil
}

func (u *User) validatePassword() error {
	if u.Password == "" {
		return errors.New("Field password is required")
	}
	if len(u.Password) < 8 {
		return errors.New("Password must have at least 8 characters")
	}

	return nil
}

func (u *User) validateEmail() error {
	emailRegexp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if u.Email == "" {
		return errors.New("Field email is required")
	}
	if !emailRegexp.MatchString(u.Email) {
		return errors.New("Invalid email format")
	}

	if _, err := GetUserByEmail(u.Email); err == nil { //ensure unique email
		return errors.New("Email already exists")
	}
	return nil
}

func (u *User) validateUsername() error {
	if u.Username == "" {
		return errors.New("Field username is required")
	}
	if _, err := GetUserByEmail(u.Username); err == nil { //ensure unique email
		return errors.New("Username already exists")
	}
	return nil
}

/*
// Delete deletes a box instance from database
func (b *Box) Delete() error {
	return boxCollection.DeleteDocument(b)
}

// Update updates a box instance from database
func (b *Box) Update(updateMap utils.JSONMap) error {
	updateMap = utils.RemoveForbiddenFields(updateMap)
	updateBytes, err := json.Marshal(updateMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(updateBytes, b)
	if err != nil {
		return err
	}

	return b.Save()
}*/

// UserList is a list of User Documents
type UserList = []User

func newUserList() UserList {
	return make([]User, 0)
}

//ListUsers returns all boxes in the box collection
func ListUsers() (users UserList) {
	users = newUserList()
	results := userCollection.Find(bson.M{})

	user := User{}
	for results.Next(&user) {
		users = append(users, user)
	}

	return
}
