package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
)

type authProvider interface {
	GetUserFromToken(token string) (*models.User, error)
}

type service struct {
	userService interfaces.UserService
}

// NewService returns an object that implements the Service interface
func NewService(r interfaces.Registry) interfaces.AuthService {
	s := &service{}

	r.RegisterService(s, "auth")
	return s
}

// GetAuthTokenByProvider returns a jwt if everything ok else, returns an error
func (s *service) GetAuthTokenByProvider(inToken, provider string) (outToken *models.Token, err error) {
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

	user, err := s.userService.FindByEmail(tokenUser.Email)
	if err != nil {
		user, err = s.userService.CreateUser(tokenUser.Username, tokenUser.Email, tokenUser.FirstName, tokenUser.LastName, tokenUser.ImageURL)
		if err != nil {
			return nil, err
		}
	}

	return s.generateToken(user)
}

func (s *service) OnAllServicesRegistered(r interfaces.Registry) {
	s.userService = r.GetService("user").(interfaces.UserService)
}

func (s *service) generateToken(u *models.User) (*models.Token, error) {
	claims := tokenClaims{*u, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 5).Unix(),
		Issuer:    "magichub.auh",
		IssuedAt:  time.Now().Unix(),
		Subject:   u.Username,
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedJWT, err := token.SignedString(mySigningKey)
	return &models.Token{JWT: signedJWT, User: u}, err
}

// tokenClaims is a struct for the JWT claims
type tokenClaims struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

var mySigningKey = []byte("AllYourBase") //TODO use real key

// GetUserFromToken returns the user held inside a JWT or error if not a valid JWT
func GetUserFromToken(tokenString string) (*models.User, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	}

	claims := tokenClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)

	return &claims.User, err
}
