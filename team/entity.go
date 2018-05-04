package team

import (
	"gopkg.in/mgo.v2/bson"
)

// A Team is the entity for the team :)
type Team struct {
	ID          bson.ObjectId `json:"id"`
	Image       string        `json:"image"`
	Name        string        `json:"name"`
	RouteName   string        `json:"routeName"`
	Description string        `json:"description"`
	// Users
	// Is user registered
}
