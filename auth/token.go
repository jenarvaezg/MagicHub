package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jenarvaezg/MagicHub/user"
)

// Token is a type that holds a JWT inside and the user that owns the token
type Token struct {
	JWT  string     `json:"jwt"`
	User *user.User `json:"user"`
}

// tokenClaims is a struct for the JWT claims
type tokenClaims struct {
	User user.User `json:"user"`
	jwt.StandardClaims
}

var mySigningKey = []byte("AllYourBase") //TODO use real key

func generateToken(u *user.User) (*Token, error) {

	claims := tokenClaims{*u, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 5).Unix(),
		Issuer:    "magichub.auh",
		IssuedAt:  time.Now().Unix(),
		Subject:   u.Username,
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedJWT, err := token.SignedString(mySigningKey)
	return &Token{JWT: signedJWT, User: u}, err
}

// GetUserFromToken returns the user held inside a JWT or error if not a valid JWT
func GetUserFromToken(tokenString string) (*user.User, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	}

	claims := tokenClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)

	return &claims.User, err
}
