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
	t.Parallel()
	repo := new(mocks.Repository)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, mock.AnythingOfType("string")).Return()

	s := team.NewService(repo, mockRegistry)
	mockRegistry.AssertCalled(t, "RegisterService", s, "team")
}

func TestGetRouteNameFromName(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	creatorID := bson.NewObjectId()
	users := []bson.ObjectId{creatorID}
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

func TestFindByID(t *testing.T) {
	t.Parallel()
	var testCases = []struct {
		testName  string
		id        bson.ObjectId
		callsRepo bool
		expected  *models.Team
		err       error
	}{
		{"call ok", bson.NewObjectId(), true, &models.Team{}, nil},
		{"call fails", bson.NewObjectId(), true, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			if tc.callsRepo {
				mockRepository.On("FindByID", mock.AnythingOfType("bson.ObjectId")).Return(tc.expected, tc.err)
			}
			r := registry.NewRegistry()

			service := team.NewService(mockRepository, r)
			result, err := service.FindByID(tc.id)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestGetTeamMembers(t *testing.T) {
	t.Parallel()
	currentUserID := bson.NewObjectId()
	currentUser, otherUser := &models.User{}, &models.User{}
	currentUser.SetId(currentUserID)
	usersWithCurrentUser := []*models.User{currentUser, otherUser}
	theTeam := &models.Team{Admins: usersWithCurrentUser, Members: usersWithCurrentUser}

	var testCases = []struct {
		testName      string
		currentUserID bson.ObjectId
		expected      []*models.User
		err           error
	}{
		{"User in member list and found", currentUserID, usersWithCurrentUser, nil},
		{"User not in member list", bson.NewObjectId(), nil, errors.New("you must be in the team to see members")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			r := registry.NewRegistry()

			service := team.NewService(mockRepository, r)
			result, err := service.GetTeamMembers(tc.currentUserID, theTeam)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)

		})
	}
}

func TestGetTeamMembersCount(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		testName string
		users    []*models.User
		expected int
		err      error
	}{
		{"Some users in team", []*models.User{&models.User{}, &models.User{}, &models.User{}}, 3, nil},
		{"No users in team", []*models.User{}, 0, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			r := registry.NewRegistry()
			theTeam := &models.Team{Members: tc.users}

			service := team.NewService(mockRepository, r)
			result, err := service.GetTeamMembersCount(theTeam)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)

		})
	}
}

func TestGetTeamAdmins(t *testing.T) {
	t.Parallel()
	currentUserID := bson.NewObjectId()
	currentUser, otherUser := &models.User{}, &models.User{}
	currentUser.SetId(currentUserID)
	usersWithCurrentUser := []*models.User{currentUser, otherUser}
	admins := []*models.User{currentUser}
	theTeam := &models.Team{Admins: admins, Members: usersWithCurrentUser}

	var testCases = []struct {
		testName      string
		currentUserId bson.ObjectId
		expected      []*models.User
		err           error
	}{
		{"User in member list and found", currentUserID, usersWithCurrentUser, nil},
		{"User not in member list", bson.NewObjectId(), nil, errors.New("you must be in the team to see admins")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			r := registry.NewRegistry()

			service := team.NewService(mockRepository, r)
			result, err := service.GetTeamAdmins(tc.currentUserId, theTeam)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)

		})
	}
}

func TestOnAllServicesRegistered(t *testing.T) {
	t.Parallel()
	mockRepository := new(mocks.Repository)
	mockTeamService := new(interfaceMocks.TeamService)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, "team")
	mockRegistry.On("GetService", "team").Return(mockTeamService)

	service := team.NewService(mockRepository, mockRegistry)
	service.OnAllServicesRegistered(mockRegistry)

	mockRegistry.AssertExpectations(t)

}
