package box

import (
	"github.com/jenarvaezg/MagicHub/models"
	"gopkg.in/mgo.v2/bson"
)

// Service is a interface of all the methods required to be an interface for Box
type Service interface {
	FindFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error)
}

type service struct {
	repo Repository
}

// NewService returns a service for Box
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// FindFiltered returns a list of boxes filtered by limit and offset
func (s *service) FindFiltered(limit, offset int, teamID bson.ObjectId) ([]*models.Box, error) {
	return s.repo.FindFiltered(limit, offset, teamID)
}
