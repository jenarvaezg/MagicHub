package auth

import (
	"fmt"

	"github.com/jenarvaezg/MagicHub/user"
)

// Service is a interface of all the methods required to be an interface for Auth
type Service interface {
	GetAuthTokenByProvider(token, provider string) (*Token, error)
}

type authProvider interface {
	GetUserFromToken(token string) (*user.User, error)
}

type service struct {
	userService user.Service
}

// NewService returns an object that implements the Service interface
func NewService(userService user.Service) Service {
	return &service{userService: userService}
}

// GetAuthTokenByProvider returns a jwt if everything ok else, returns an error
func (s *service) GetAuthTokenByProvider(inToken, provider string) (outToken *Token, err error) {
	var authProvider authProvider
	switch provider {
	case "google":
		authProvider = newGoogleAuthProvider(s.userService)
	default:
		err = fmt.Errorf("Provider %s is not supported", provider)
	}

	user, err := authProvider.GetUserFromToken(inToken)
	if err != nil {
		return nil, err
	}

	user, err = s.userService.FindByEmail(user.Email)
	if err != nil {
		user, err = s.userService.CreateUser(user.Username, user.Email, user.FirstName, user.LastName, user.ImageURL)
		if err != nil {
			return nil, err
		}
	}
	// TODO update user data from google

	outToken, err = generateToken(user)
	return outToken, err
}
