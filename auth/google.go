package auth

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jenarvaezg/magicbox/models"
	"github.com/mendsley/gojwk"
	"golang.org/x/oauth2"
)

const (
	googleCertURL    = "https://www.googleapis.com/oauth2/v3/certs"
	expectedAudience = "144467579021-gmpp7n1a9m3b82svfs51eqjbs0bhidkk.apps.googleusercontent.com"
	expectedIssuer   = "accounts.google.com"
)

// GoogleUserProfile is the representation of the profile given by Google
type GoogleUserProfile struct {
	ID         string `json:"googleId"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	ImageURL   string `json:"imageUrl"`
}

// GoogleToken is the  basic OAuth token with google's jwt
type GoogleToken struct {
	oauth2.Token `json:",inline"`
	JWT          string `json:"id_token"`
}

// GoogleFrontendRequest is the request sent from frontend to validate google auth
type GoogleFrontendRequest struct {
	Token   GoogleToken       `json:"tokenObj"`
	Profile GoogleUserProfile `json:"profileObj"`
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

func parseToken(googleToken GoogleToken) (*jwt.Token, error) {
	var token *jwt.Token
	keys, err := getGoogleKeys()
	if err != nil {
		log.Println(err)
		return &jwt.Token{}, err
	}

	for _, key := range keys {
		token, err = jwt.Parse(googleToken.JWT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return key.DecodePublicKey()
		})
		if err == nil {
			return token, nil
		}
	}
	return nil, errors.New("No google cert match found")
}

func validateGoogleToken(token GoogleToken) error {

	parsedToken, err := parseToken(token)
	if err != nil {
		return err
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	if err := claims.Valid(); !claims.VerifyAudience(expectedAudience, true) ||
		!claims.VerifyIssuer(expectedIssuer, true) || err != nil {
		return errors.New("Unvalid jwt internals")
	}

	return nil
}

func userRequestFromGoogleProfile(profile GoogleUserProfile) models.UserRequest {
	pass := ""
	req := models.UserRequest{
		Username:   strings.Split(profile.Email, "@")[0],
		Email:      profile.Email,
		FirstName:  profile.GivenName,
		LastName:   profile.FamilyName,
		Password:   &pass,
		FromGoogle: true,
	}

	return req
}
