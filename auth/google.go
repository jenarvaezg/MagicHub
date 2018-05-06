package auth

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/mendsley/gojwk"
)

const (
	googleCertURL    = "https://www.googleapis.com/oauth2/v3/certs"
	expectedAudience = "144467579021-gmpp7n1a9m3b82svfs51eqjbs0bhidkk.apps.googleusercontent.com"
	expectedIssuer   = "accounts.google.com"
)

type googleToken = string

type googleAuthProvider struct {
	userService user.Service
}

func newGoogleAuthProvider(userService user.Service) authProvider {
	return &googleAuthProvider{userService: userService}
}

func getGoogleKeys() (keys []*gojwk.Key, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //TODO REMOVE
	r, err := http.Get(googleCertURL)
	if err != nil {
		return
	}
	defer r.Body.Close()

	var key gojwk.Key

	err = json.NewDecoder(r.Body).Decode(&key)
	keys = key.Keys
	return
}

func parseToken(googleToken googleToken) (*jwt.Token, error) {
	keys, err := getGoogleKeys()
	if err != nil {
		return &jwt.Token{}, err
	}

	return jwt.Parse(googleToken, func(token *jwt.Token) (interface{}, error) {
		var keyToUse *gojwk.Key
		for _, key := range keys {
			if key.Kid == token.Header["kid"] {
				keyToUse = key
				break
			}
		}

		if keyToUse == nil {
			return nil, fmt.Errorf("Didn't find any appropiate key at google for the provided JWT")
		}

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return keyToUse.DecodePublicKey()
	})
}

func (p *googleAuthProvider) GetUserFromToken(token string) (*user.User, error) {
	parsedToken, err := parseToken(googleToken(token))
	if err != nil {
		return nil, err
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	if err := claims.Valid(); !claims.VerifyAudience(expectedAudience, true) ||
		!claims.VerifyIssuer(expectedIssuer, true) || err != nil {
		return nil, errors.New("Unvalid jwt internals")
	}

	return userFromClaims(claims), nil
}

func userFromClaims(claims jwt.MapClaims) *user.User {
	return &user.User{
		Username:  strings.Split(claims["email"].(string), "@")[0],
		Email:     claims["email"].(string),
		FirstName: claims["given_name"].(string),
		LastName:  claims["family_name"].(string),
		ImageURL:  claims["picture"].(string),
	}
}
