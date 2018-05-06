package team

import (
	"strings"
)

// Service is a interface of all the methods required to be an interface for Team
type Service interface {
	GetRouteNameFromName(string) string
	FindFiltered(limit, offset int, search string) ([]*Team, error)
	CreateTeam(name, image, description string) (*Team, error)
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
func (s *service) FindFiltered(limit, offset int, search string) ([]*Team, error) {
	return s.repo.FindFiltered(limit, offset, search)
}

// CreateTeam creates a team :)
func (s *service) CreateTeam(name, image, description string) (*Team, error) {
	team := &Team{Name: name, Image: image, Description: description}

	_, err := s.repo.Store(team)

	return team, err
}
