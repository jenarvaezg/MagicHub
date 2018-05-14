package team

import (
	"fmt"
	"strings"

	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	repo Repository
}

// NewService returns an object that implements the Service interface
func NewService(repo Repository, r interfaces.Registry) interfaces.TeamService {
	s := &service{repo: repo}

	r.RegisterService(s, "team")
	return s
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

	if _, err := s.repo.Store(team); err != nil {
		return nil, err
	}

	return team, nil
}

// GetTeamByID returns a team by its ID or error if not found
func (s *service) GetTeamByID(id string) (*models.Team, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("%s is not a valid ID", id)
	}
	return s.repo.FindByID(bson.ObjectIdHex(id))
}

// GetTeamMembers returns the list of members that belong to a team or an error if the user is not in the team
func (s *service) GetTeamMembers(userID bson.ObjectId, team *models.Team) ([]*models.User, error) {
	if team.IsUserMember(userID) {
		return team.Members.([]*models.User), nil
	}
	return nil, fmt.Errorf("you must be in the team to see members")
}

// GetTeamMembersCount returns number members that belong to a team or an error if the user is not in the team
func (s *service) GetTeamMembersCount(team *models.Team) (int, error) {
	return len(team.Members.([]*models.User)), nil
}

// user returns the list of admin that belong to a team or an error if the user is not in the team
func (s *service) GetTeamAdmins(userID bson.ObjectId, team *models.Team) ([]*models.User, error) {
	if team.IsUserMember(userID) {
		return team.Members.([]*models.User), nil
	}
	return nil, fmt.Errorf("you must be in the team to see admins")
}

func (s *service) OnAllServicesRegistered(r interfaces.Registry) {
	// As of now Team service does not need other services
}
