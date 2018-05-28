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
	userService interfaces.UserService
}

// NewService returns a service for Box
func NewService(repo Repository, r interfaces.Registry) interfaces.BoxService {
	s := &service{repo: repo}

	r.RegisterService(s, "box")
	return s
}

// FindByTeamFiltered returns a list of boxes filtered by limit and offset
func (s *service) FindByTeamFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error) {
	return s.repo.FindByTeamFiltered(limit, offset, teamID)
}

// CreateBox creates a box and saves it
func (s *service) CreateBox(userID, teamID bson.ObjectId, name string, openDate time.Time) (*models.Box, error) {
	team, err := s.teamService.FindByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("could not create box, finding team: %v", err)
	}

	if !team.IsUserMember(userID) {
		return nil, fmt.Errorf("you are not in the team %s so you can't create boxes", teamID)
	}

	box := &models.Box{Name: name, OpenDate: openDate, Notes: []*models.Note{}, Team: teamID}
	_, err = s.repo.Store(box)
	if err != nil {
		return nil, err
	}

	return box, nil
}

// InsertNote inserts a note into a box and updates it
func (s *service) InsertNote(userID, boxID bson.ObjectId, text string) (*models.Box, error) {
	box, err := s.repo.FindByID(boxID)
	if err != nil {
		return nil, err
	}

	if team := box.Team.(*models.Team); !team.IsUserMember(userID) {
		return nil, fmt.Errorf("you are not in the team %s so you can't add notes", team.GetId())
	}

	if err := box.AddNote(models.Note{Text: text, From: userID}); err != nil {
		return nil, err
	}

	if _, err := s.repo.Store(box); err != nil {
		return nil, err
	}

	return box, nil
}

// GetNotes returns all the notes an user can see in a given moment, it can be all notes if the box is open or only the notes submitted by the user if
// the box is not open yet.
func (s *service) GetNotes(userID bson.ObjectId, box *models.Box) ([]*models.Note, error) {
	var err error
	notes := []*models.Note{}
	isOpen := box.IsOpen()
	for _, n := range box.Notes {
		if isOpen || userID == n.From {
			n.From, err = s.userService.FindByID(n.From.(bson.ObjectId))
			if err != nil {
				return nil, err
			}
			notes = append(notes, n)
		}
	}
	return notes, nil
}

// OnAllServicesRegistered is a method called when all services are registered, used to solve cyclic dependencies
func (s *service) OnAllServicesRegistered(sr interfaces.Registry) {
	s.teamService = sr.GetService("team").(interfaces.TeamService)
	s.userService = sr.GetService("user").(interfaces.UserService)
}
