package auth

import (
	"fmt"
	"log"

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
		return nil, fmt.Errorf("Provider %s is not supported", provider)
	}

	tokenUser, err := authProvider.GetUserFromToken(inToken)
	if err != nil {
		return nil, err
	}

	log.Println(tokenUser.LastName)
	user, err := s.userService.FindByEmail(tokenUser.Email)
	log.Println(user.LastName)
	if err != nil {
		user, err = s.userService.CreateUser(tokenUser.Username, tokenUser.Email, tokenUser.FirstName, tokenUser.LastName, tokenUser.ImageURL)
		if err != nil {
			return nil, err
		}
	}

	return generateToken(user)
}
