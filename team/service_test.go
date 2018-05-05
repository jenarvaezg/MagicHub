package team_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"

	"github.com/jenarvaezg/MagicHub/team"
	"github.com/jenarvaezg/MagicHub/team/mocks"
	"github.com/stretchr/testify/mock"
)

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
			mockRepository := new(mocks.Repository)
			service := team.NewService(mockRepository)

			result := service.GetRouteNameFromName(tc.input)

			assert.Equal(t, result, tc.expected)
		})
	}
}

func TestFindFiltered(t *testing.T) {
	mockTeam := &team.Team{}
	mockTeamList := []*team.Team{mockTeam}
	var testCases = []struct {
		name     string
		mocks    []*team.Team
		expected []*team.Team
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

			service := team.NewService(mockRepository)
			result, err := service.FindFiltered(1, 1, "")

			assert.Equal(t, result, tc.expected)
			assert.Equal(t, err, tc.err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestCreateTeam(t *testing.T) {
	var testCases = []struct {
		testName    string
		name        string
		image       string
		description string
		expected    *team.Team
		err         error
	}{
		{"call ok", "name", "image", "description", &team.Team{Name: "name", Image: "image", Description: "description"}, nil},
		{"call fails", "", "", "", &team.Team{}, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockRepository.On("Store", mock.Anything).Return(bson.NewObjectId(), tc.err)

			service := team.NewService(mockRepository)
			result, err := service.CreateTeam(tc.name, tc.image, tc.description)

			assert.Equal(t, result, tc.expected)
			assert.Equal(t, err, tc.err)
			mockRepository.AssertExpectations(t)
		})
	}
}
