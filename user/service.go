package user

import (
	"fmt"

	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	repo Repository
}

// NewService returns an object that implements the Service interface
func NewService(repo Repository, r interfaces.Registry) interfaces.UserService {
	s := &service{repo: repo}

	r.RegisterService(s, "user")
	return s
}

// FindByID returns a users matching an ID
func (s *service) FindByID(id string) (*models.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("%s is not a valid ID", id)
	}
	return s.repo.FindByID(bson.ObjectIdHex(id))
}

// FindByEmail returns a users matching an email
func (s *service) FindByEmail(email string) (*models.User, error) {
	users, err := s.repo.FindBy(map[string]interface{}{"email": email})

	if err != nil {
		return nil, fmt.Errorf("find user by email: %v", err.Error())
	} else if len(users) == 0 {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	return users[0], nil

}

// CreateUser creates a user :)
func (s *service) CreateUser(username, email, firstName, lastName, imageURL string) (*models.User, error) {
	user := &models.User{Username: username, Email: email, FirstName: firstName, LastName: lastName, ImageURL: imageURL}

	_, err := s.repo.Store(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) OnAllServicesRegistered(r interfaces.Registry) {
	// User service does not need any other service as of now
}
