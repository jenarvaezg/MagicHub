package models

import (
	"errors"

	"github.com/jenarvaezg/mongodm"
	"gopkg.in/mgo.v2/bson"
)

// A Team is the entity for the team :)
type Team struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Image                string      `json:"image"`
	Name                 string      `json:"name"`
	RouteName            string      `json:"routeName"`
	Description          string      `json:"description"`
	Members              interface{} `bson:"members" json:"members" model:"User" relation:"1n" autosave:"true"`
	Admins               interface{} `bson:"admins" json:"admins" model:"User" relation:"1n" autosave:"true"`
	JoinRequests         interface{} `bson:"join_requests" json:"joinRequests" model:"User" relation:"1n" autosave:"true"`
}

// IsUserMember resturn a boolean that determines wheter a userID is in the team members list
func (t Team) IsUserMember(userID bson.ObjectId) bool {
	if len(t.Members.([]*User)) == 0 {
		return false
	}
	for _, m := range t.Members.([]*User) {
		if userID == m.GetId() {
			return true
		}
	}

	return false
}

// IsUserAdmin resturn a boolean that determines wheter a userID is in the team admin list
func (t Team) IsUserAdmin(userID bson.ObjectId) bool {
	if len(t.Admins.([]*User)) == 0 {
		return false
	}
	for _, m := range t.Admins.([]*User) {
		if userID == m.GetId() {
			return true
		}
	}

	return false
}

// AddInviteRequest adds a user to a team's requests list if not already
func (t *Team) AddInviteRequest(user *User) error {
	for _, req := range t.JoinRequests.([]*User) {
		if user.GetId() == req.GetId() {
			return errors.New("user already requested to join")
		}
	}
	t.JoinRequests = append(t.JoinRequests.([]*User), user)
	return nil
}

// AcceptInviteRequest removes an user from the requests list and adds it to members list
func (t *Team) AcceptInviteRequest(user *User) error {
	reqs := t.JoinRequests.([]*User)
	for i, req := range t.JoinRequests.([]*User) {
		if user.GetId() == req.GetId() {
			reqs = append(reqs[:i], reqs[i+1:]...)
			t.JoinRequests = reqs
			t.Members = append(t.Members.([]*User), user)
			return nil
		}
	}
	return errors.New("user is not in the join request list")
}
