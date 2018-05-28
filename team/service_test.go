package team_test

import (
	"errors"
	"log"
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

func TestRequestTeamInvite(t *testing.T) {
	userID := bson.NewObjectId()
	user := &models.User{}
	user.SetId(userID)
	emptyTeam := func() *models.Team {
		return &models.Team{Members: []*models.User{}, JoinRequests: []*models.User{}}
	}
	teamWithUser := func() *models.Team {
		return &models.Team{Members: []*models.User{user}}
	}
	teamWithRequest := func() *models.Team {
		return &models.Team{Members: []*models.User{}, JoinRequests: []*models.User{user}}
	}
	var testCases = []struct {
		testName     string
		team         *models.Team
		getTeamError error
		getUserError error
		storeError   error
		expected     *models.Team
		err          error
	}{
		{"Team not found", emptyTeam(), errors.New("fail"), nil, nil, nil, errors.New("could not get team: fail")},
		{"User in member list", teamWithUser(), nil, nil, nil, nil, errors.New("you are already in the team")},
		{"User not found", emptyTeam(), nil, errors.New("fail"), nil, nil, errors.New("could not fetch user: fail")},
		{"User already requested", teamWithRequest(), nil, nil, nil, nil, errors.New("could not add invite request: user already requested to join")},
		{"Store call fails", emptyTeam(), nil, nil, errors.New("fail"), nil, errors.New("could not save team: fail")},
		{"Everything OK", emptyTeam(), nil, nil, nil, teamWithRequest(), nil},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockUserService := new(interfaceMocks.UserService)
			mockUserService.On("OnAllServicesRegistered", mock.Anything)
			mockRepository.On("FindByID", tc.team.GetId()).Return(tc.team, tc.getTeamError)
			if tc.getTeamError == nil {
				mockUserService.On("FindByID", user.GetId()).Return(user, tc.getUserError)
				mockRepository.On("Store", tc.team).Return(tc.team.GetId(), tc.storeError)
			}
			r := registry.NewRegistry()
			r.RegisterService(mockUserService, "user")

			service := team.NewService(mockRepository, r)
			r.AllServicesRegistered()
			result, err := service.RequestTeamInvite(user.GetId(), tc.team.GetId())
			log.Println(tc.testName, tc.team.JoinRequests)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)

		})
	}

}

func TestAcceptInviteRequest(t *testing.T) {
	userID, requesterID := bson.NewObjectId(), bson.NewObjectId()
	user, requester := &models.User{}, &models.User{}
	user.SetId(userID)
	requester.SetId(requesterID)

	emptyTeam := func() *models.Team {
		return &models.Team{Members: []*models.User{}, Admins: []*models.User{}, JoinRequests: []*models.User{}}
	}
	teamWithUser := func() *models.Team {
		return &models.Team{Members: []*models.User{user}, Admins: []*models.User{user}, JoinRequests: []*models.User{}}
	}
	teamWithRequest := func() *models.Team {
		return &models.Team{Members: []*models.User{user}, Admins: []*models.User{user}, JoinRequests: []*models.User{requester}}
	}
	teamWithUserAndRequester := func() *models.Team {
		return &models.Team{Members: []*models.User{user, requester}, Admins: []*models.User{user}, JoinRequests: []*models.User{}}
	}

	_, _ = teamWithUser(), teamWithRequest()
	var testCases = []struct {
		testName     string
		team         *models.Team
		getTeamError error
		getUserError error
		storeError   error
		expected     *models.Team
		err          error
	}{
		{"Team not found", emptyTeam(), errors.New("fail"), nil, nil, nil, errors.New("could not get team: fail")},
		{"User not admin", emptyTeam(), nil, nil, nil, nil, errors.New("you are not an admin of the team")},
		{"User not found", teamWithUser(), nil, errors.New("fail"), nil, nil, errors.New("could not fetch user: fail")},
		{"User did not request", teamWithUser(), nil, nil, nil, nil, errors.New("could not add invite request: user is not in the join request list")},
		{"Store call fails", teamWithRequest(), nil, nil, errors.New("fail"), nil, errors.New("could not save team: fail")},
		{"Everything OK", teamWithRequest(), nil, nil, nil, teamWithUserAndRequester(), nil},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockUserService := new(interfaceMocks.UserService)
			mockUserService.On("OnAllServicesRegistered", mock.Anything)
			mockRepository.On("FindByID", tc.team.GetId()).Return(tc.team, tc.getTeamError)
			if tc.getTeamError == nil {
				mockUserService.On("FindByID", requester.GetId()).Return(requester, tc.getUserError)
				mockRepository.On("Store", tc.team).Return(tc.team.GetId(), tc.storeError)
			}
			r := registry.NewRegistry()
			r.RegisterService(mockUserService, "user")

			service := team.NewService(mockRepository, r)
			r.AllServicesRegistered()
			result, err := service.AcceptInviteRequest(user.GetId(), requester.GetId(), tc.team.GetId())

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)

		})
	}

}
func TestOnAllServicesRegistered(t *testing.T) {
	t.Parallel()
	mockRepository := new(mocks.Repository)
	mockUserService := new(interfaceMocks.UserService)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, "team")
	mockRegistry.On("GetService", "user").Return(mockUserService)

	service := team.NewService(mockRepository, mockRegistry)
	service.OnAllServicesRegistered(mockRegistry)

	mockRegistry.AssertExpectations(t)

}
