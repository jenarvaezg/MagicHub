package user

import (
	"context"
	"log"

	"gopkg.in/mgo.v2/bson"
)

// ContextKey is //TODO
type ContextKey string

//ContextKeyCurrentUser is a key used for indexing a user in a context
var ContextKeyCurrentUser = ContextKey("current-user")

// RequireUser returns user from http context or panics if not found
func RequireUser(c context.Context) bson.ObjectId {
	userID, ok := c.Value(ContextKeyCurrentUser).(bson.ObjectId)
	if !ok {
		log.Println("Require user got", c.Value(ContextKeyCurrentUser))
		panic("User not authenticated")
	}

	return userID
}

// StoreUserIDInContext stores a user id in a context and returs a context with that value
func StoreUserIDInContext(c context.Context, userID bson.ObjectId) context.Context {
	return context.WithValue(c, ContextKeyCurrentUser, userID)
}
