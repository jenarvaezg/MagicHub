package auth

import (
	"errors"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jenarvaezg/magicbox/models"
)

const (
	grantTypePassword = "password"
)

// TokenClaims is a struct for the JWT claims
type TokenClaims struct {
	User models.UserResponse `json:"user"`
	jwt.StandardClaims
}

var mySigningKey = []byte("AllYourBase") //TODO use real key

func getJWT(user models.User) (string, error) {

	claims := TokenClaims{user.GetResponse(), jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    "magicbox.auh",
		IssuedAt:  time.Now().Unix(),
		Subject:   user.Username,
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(mySigningKey)
}

// GetAuthToken returns an auth token for the requesting user and passoword, or an error
func GetAuthToken(form url.Values) (token string, err error) {
	grantType := form.Get("grant_type")
	if grantType == grantTypePassword {
		username, password := form.Get("username"), form.Get("password")
		user, err := models.GetUserByUsername(username)
		if err != nil {
			return token, errors.New("Invalid username or password")
		}
		if ok := user.ChallengePassword(password); !ok {
			return token, errors.New("Invalid username or password")
		}

		return getJWT(user)

	}
	return token, errors.New("Unsupported grant type")

}
