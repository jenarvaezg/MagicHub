package auth

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jenarvaezg/MagicHub/models"
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
		Issuer:    "magichub.auh",
		IssuedAt:  time.Now().Unix(),
		Subject:   user.Username,
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(mySigningKey)
}

// GetAuthTokenFromForm returns an auth token for the requesting user and passoword, or an error
func GetAuthTokenFromForm(form url.Values) (token string, err error) {
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

		return getJWT(*user)

	}
	return token, errors.New("Unsupported grant type")
}

/*
GetAuthTokenFromGoogleToken returns an auth token from a frontend google auth requests
If user does not exist in database, it is created
*/
func GetAuthTokenFromGoogleToken(googleReq GoogleFrontendRequest) (token string, err error) {
	if err = validateGoogleToken(googleReq.Token); err != nil {
		return
	}

	email := googleReq.Profile.Email
	req := userRequestFromGoogleProfile(googleReq.Profile)
	user, err := models.GetUserByEmail(email)
	if err != nil {
		user, err = models.NewUser(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = user.Save(); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if err = user.Update(req); err != nil {
			fmt.Println(err)
			return
		}
	}

	return getJWT(*user)
}
