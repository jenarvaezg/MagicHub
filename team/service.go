package team

import (
	"fmt"
	"strings"

	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/user"
	"gopkg.in/mgo.v2/bson"
)

// Service is a interface of all the methods required to be an interface for Team
type Service interface {
	GetRouteNameFromName(string) string
	FindFiltered(limit, offset int, search string) ([]*models.Team, error)
	CreateTeam(userID bson.ObjectId, name, image, description string) (*models.Team, error)
	GetTeamByID(id string) (*models.Team, error)
	GetTeamMembers(userID bson.ObjectId, team *models.Team) ([]*user.User, error)
	GetTeamMembersCount(team *models.Team) (int, error)
	GetTeamAdmins(userID bson.ObjectId, team *models.Team) ([]*user.User, error)
}

type service struct {
	repo Repository
}

// NewService returns an object that implements the Service interface
func NewService(repo Repository) Service {
	return &service{repo}
}

// GetRouteNameFromName returns a route name from a name
func (s *service) GetRouteNameFromName(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "", -1))
}

// FindFiltered returns a list of teams filtered by limit offset and search params
func (s *service) FindFiltered(limit, offset int, search string) ([]*models.Team, error) {
	return s.repo.FindFiltered(limit, offset, search)
}

// CreateTeam creates a team :)
func (s *service) CreateTeam(creatorID bson.ObjectId, name, image, description string) (*models.Team, error) {
	team := &models.Team{Name: name, Image: image, Description: description}
	team.RouteName = s.GetRouteNameFromName(team.Name)
	team.Members = []interface{}{creatorID}
	team.Admins = []interface{}{creatorID}

	_, err := s.repo.Store(team)

	return team, err
}

// GetTeamByID returns a team by its ID or error if not found
func (s *service) GetTeamByID(id string) (*models.Team, error) {

	return s.repo.FindByID(bson.ObjectId(id))
}

// GetTeamMembers returns the list of members that belong to a team or an error if the user is not in the team
func (s *service) GetTeamMembers(userID bson.ObjectId, team *models.Team) ([]*user.User, error) {
	if team.IsUserMember(userID) {
		return team.Members.([]*user.User), nil
	}
	return nil, fmt.Errorf("you must be in the team to see members")
}

// GetTeamMembersCount returns number members that belong to a team or an error if the user is not in the team
func (s *service) GetTeamMembersCount(team *models.Team) (int, error) {
	return len(team.Members.([]*user.User)), nil
}

// user returns the list of admin that belong to a team or an error if the user is not in the team
func (s *service) GetTeamAdmins(userID bson.ObjectId, team *models.Team) ([]*user.User, error) {
	if team.IsUserMember(userID) {
		return team.Members.([]*user.User), nil
	}
	return nil, fmt.Errorf("you must be in the team to see admins")
}
