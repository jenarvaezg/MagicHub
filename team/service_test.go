package team_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"

	interfaceMocks "github.com/jenarvaezg/MagicHub/interfaces/mocks"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/registry"
	"github.com/jenarvaezg/MagicHub/team"
	"github.com/jenarvaezg/MagicHub/team/mocks"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	repo := new(mocks.Repository)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, mock.AnythingOfType("string")).Return()

	s := team.NewService(repo, mockRegistry)
	mockRegistry.AssertCalled(t, "RegisterService", s, "team")

}

func TestGetRouteNameFromName(t *testing.T) {
	var testCases = []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"basic string", "Chameleon", "chameleon"},
		{"string with spaces", "A team with spaces", "ateamwithspaces"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := registry.NewRegistry()
			mockRepository := new(mocks.Repository)
			service := team.NewService(mockRepository, r)

			result := service.GetRouteNameFromName(tc.input)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFindFiltered(t *testing.T) {
	mockTeam := &models.Team{}
	mockTeamList := []*models.Team{mockTeam}
	var testCases = []struct {
		name     string
		mocks    []*models.Team
		expected []*models.Team
		err      error
	}{
		{"call ok", mockTeamList, mockTeamList, nil},
		{"call fails", nil, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockCall := mockRepository.On("FindFiltered", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string"))
			mockCall.Return(tc.mocks, tc.err)
			r := registry.NewRegistry()

			service := team.NewService(mockRepository, r)
			result, err := service.FindFiltered(1, 1, "")

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestCreateTeam(t *testing.T) {
	creatorID := bson.NewObjectId()
	users := []interface{}{creatorID}
	name, image, description := "name", "image", "description"
	expectedTeamOK := &models.Team{Name: name, Image: image, Description: description, Members: users, Admins: users, RouteName: name}
	var testCases = []struct {
		testName    string
		name        string
		image       string
		description string
		expected    *models.Team
		err         error
	}{
		{"call ok", name, image, description, expectedTeamOK, nil},
		{"call fails", "", "", "", nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockRepository.On("Store", mock.Anything).Return(bson.NewObjectId(), tc.err)
			r := registry.NewRegistry()

			service := team.NewService(mockRepository, r)
			result, err := service.CreateTeam(creatorID, tc.name, tc.image, tc.description)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestGetTeamByID(t *testing.T) {
	var testCases = []struct {
		testName  string
		id        string
		callsRepo bool
		expected  *models.Team
		err       error
	}{
		{"call ok", bson.NewObjectId().Hex(), true, &models.Team{}, nil},
		{"bad objectid", "This is a string", false, nil, errors.New("This is a string is not a valid ID")},
		{"call fails", bson.NewObjectId().Hex(), true, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			if tc.callsRepo {
				mockRepository.On("FindByID", mock.AnythingOfType("bson.ObjectId")).Return(tc.expected, tc.err)
			}
			r := registry.NewRegistry()

			service := team.NewService(mockRepository, r)
			result, err := service.GetTeamByID(tc.id)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
			if !tc.callsRepo {
				mockRepository.AssertNotCalled(t, "FindByID", mock.AnythingOfType("bson.ObjectId"))
			}
		})
	}
}
