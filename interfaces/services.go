package interfaces

import (
	"time"

	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

// Service is a interface of all the methods required to be a Service
type Service interface {
	OnAllServicesRegistered(r Registry)
}

// TeamService is a interface of all the methods required to be an interface for Team
type TeamService interface {
	Service
	GetRouteNameFromName(string) string
	FindFiltered(limit, offset int, search string) ([]*models.Team, error)
	CreateTeam(userID bson.ObjectId, name, image, description string) (*models.Team, error)
	FindByID(id bson.ObjectId) (*models.Team, error)
	GetTeamMembers(userID bson.ObjectId, team *models.Team) ([]*models.User, error)
	GetTeamMembersCount(team *models.Team) (int, error)
	GetTeamAdmins(userID bson.ObjectId, team *models.Team) ([]*models.User, error)
}

// AuthService is a interface of all the methods required to be an interface for Auth
type AuthService interface {
	Service
	GetAuthTokenByProvider(token, provider string) (*models.Token, error)
}

// UserService is a interface of all the methods required to be an interface for User
type UserService interface {
	Service
	FindByID(id bson.ObjectId) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	CreateUser(username, email, firstName, lastName, imageURL string) (*models.User, error)
}

// BoxService is a interface of all the methods required to be an interface for Box
type BoxService interface {
	Service
	FindByTeamFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error)
	CreateBox(userID, teamID bson.ObjectId, name string, openDate time.Time) (*models.Box, error)
	InsertNote(userID bson.ObjectId, boxID bson.ObjectId, text string) (*models.Box, error)
	GetNotes(userID bson.ObjectId, box *models.Box) ([]*models.Note, error)
}
