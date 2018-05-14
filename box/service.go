package box

import (
	"github.com/jenarvaezg/MagicHub/interfaces"
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	repo Repository
}

// NewService returns a service for Box
func NewService(repo Repository, r interfaces.Registry) interfaces.BoxService {
	s := &service{repo: repo}

	r.RegisterService(s, "box")
	return s
}

// FindFiltered returns a list of boxes filtered by limit and offset
func (s *service) FindFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error) {
	return s.repo.FindFiltered(limit, offset, teamID)
}

func (s *service) OnAllServicesRegistered(sr interfaces.Registry) {
	// As of now Box service does not need other services
}
