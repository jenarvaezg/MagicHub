package box

import (
	"fmt"
	"time"

	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	repo        Repository
	teamService interfaces.TeamService
}

// NewService returns a service for Box
func NewService(repo Repository, r interfaces.Registry) interfaces.BoxService {
	s := &service{repo: repo}

	r.RegisterService(s, "box")
	return s
}

// FindByTeamFiltered returns a list of boxes filtered by limit and offset
func (s *service) FindByTeamFiltered(limit, offset int, teamID string) ([]*models.Box, error) {
	if !bson.IsObjectIdHex(teamID) {
		return nil, fmt.Errorf("%s is no a valid objectid", teamID)
	}
	return s.repo.FindByTeamFiltered(limit, offset, bson.ObjectIdHex(teamID))
}

func (s *service) CreateBox(userID bson.ObjectId, name, teamID string, openDate time.Time) (*models.Box, error) {
	team, err := s.teamService.FindByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("could not create box, finding team: %v", err)
	}

	if !team.IsUserMember(userID) {
		return nil, fmt.Errorf("you are not in the team %s so you can't create boxes", teamID)
	}

	box := &models.Box{Name: name, OpenDate: openDate, Status: models.BoxStatusClosed, Notes: []models.Note{}, Team: bson.ObjectIdHex(teamID)}
	_, err = s.repo.Store(box)
	if err != nil {
		return nil, err
	}

	return box, nil
}

func (s *service) OnAllServicesRegistered(sr interfaces.Registry) {
	s.teamService = sr.GetService("team").(interfaces.TeamService)
}
