package user

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Service is a interface of all the methods required to be an interface for User
type Service interface {
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	CreateUser(username, email, firstName, lastName, imageURL string) (*User, error)
}

type service struct {
	repo Repository
}

// NewService returns an object that implements the Service interface
func NewService(repo Repository) Service {
	return &service{repo}
}

// FindByID returns a users matching an ID
func (s *service) FindByID(id string) (*User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("%s is not a valid ID", id)
	}
	return s.repo.FindByID(bson.ObjectIdHex(id))
}

// FindByEmail returns a users matching an email
func (s *service) FindByEmail(email string) (*User, error) {
	users, err := s.repo.FindBy(map[string]interface{}{"email": email})

	if err != nil {
		return nil, fmt.Errorf("FindByEmail: %v", err.Error())
	} else if len(users) == 0 {
		return nil, fmt.Errorf("FindByEmail: User with email %s not found", email)
	}
	return users[0], nil

}

// CreateUser creates a user :)
func (s *service) CreateUser(username, email, firstName, lastName, imageURL string) (*User, error) {
	user := &User{Username: username, Email: email, FirstName: firstName, LastName: lastName, ImageURL: imageURL}

	_, err := s.repo.Store(user)

	return user, err
}
