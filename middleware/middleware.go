package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jenarvaezg/MagicHub/auth"
	"github.com/jenarvaezg/MagicHub/user"
)

//UserFromJWTMiddleware is a middleware that varifies a JWT in the Authorization header and sets the user in the conext
type UserFromJWTMiddleware struct {
}

// NewUserFromJWTMiddleware returns a RequireUserMiddleware
func NewUserFromJWTMiddleware() *UserFromJWTMiddleware {
	return &UserFromJWTMiddleware{}
}

func extractJWTFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Missing Authorization header")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

/*
UserFromJWTMiddleware's handler, extracts JWT from auth header, validates JWT and inserts user in the request context
*/
func (l *UserFromJWTMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := extractJWTFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		next(w, r)
		return
	}

	u, err := auth.GetUserFromToken(token)
	if err != nil {
		log.Println(u, err)
		next(w, r)
		return
	}

	r = r.WithContext(user.StoreUserIDInContext(r.Context(), u.GetId()))
	next(w, r)
}
