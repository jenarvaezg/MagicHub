package models

import (
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

// A Team is the entity for the team :)
type Team struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Image                string      `json:"image"`
	Name                 string      `json:"name"`
	RouteName            string      `json:"routeName"`
	Description          string      `json:"description"`
	Members              interface{} `json:"members" model:"User" relation:"1n" autosave:"true"`
	Admins               interface{} `json:"admins" model:"User" relation:"1n" autosave:"true"`
}

// IsUserMember resturn a boolean that determines wheter a userID is in the team members list
func (t Team) IsUserMember(userID bson.ObjectId) bool {

	for _, m := range t.Members.([]*user.User) {
		if userID == m.GetId() {
			return true
		}
	}

	return false
}
