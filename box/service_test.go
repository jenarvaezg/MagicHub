package box_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jenarvaezg/MagicHub/box"
	"github.com/jenarvaezg/MagicHub/box/mocks"
	interfaceMocks "github.com/jenarvaezg/MagicHub/interfaces/mocks"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	repo := new(mocks.Repository)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, mock.AnythingOfType("string")).Return()

	s := box.NewService(repo, mockRegistry)
	mockRegistry.AssertCalled(t, "RegisterService", s, "box")
}

func TestFindByTeamFiltered(t *testing.T) {
	t.Parallel()
	mockBox := &models.Box{}
	mockBoxList := []*models.Box{mockBox}
	var testCases = []struct {
		name   string
		limit  int
		offset int
		teamID bson.ObjectId

		mocks    []*models.Box
		expected []*models.Box
		err      error
	}{
		{"call ok", 1, 1, bson.NewObjectId(), mockBoxList, mockBoxList, nil},
		{"call fails", 1, 1, bson.NewObjectId(), nil, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockCall := mockRepository.On("FindByTeamFiltered", tc.limit, tc.offset, tc.teamID)
			mockCall.Return(tc.mocks, tc.err)
			r := registry.NewRegistry()

			service := box.NewService(mockRepository, r)
			result, err := service.FindByTeamFiltered(tc.limit, tc.offset, tc.teamID)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestCreateBox(t *testing.T) {
	t.Parallel()
	userID, teamID, boxName, openDate := bson.NewObjectId(), bson.NewObjectId(), "box", time.Now()
	user := &models.User{}
	user.SetId(userID)
	expectedBoxOK := &models.Box{Team: teamID, Name: boxName, OpenDate: openDate, Notes: []*models.Note{}}
	expectedTeamOK := &models.Team{Members: []*models.User{user}}
	var testCases = []struct {
		testName     string
		expectedTeam *models.Team
		findTeamErr  error
		expected     *models.Box
		err          error
	}{
		{"team does not exist", nil, errors.New("fail"), nil, errors.New("could not create box, finding team: fail")},
		{"user not in team", &models.Team{Members: []*models.User{&models.User{}}}, nil, nil, fmt.Errorf("you are not in the team %s so you can't create boxes", teamID)},
		{"call ok", expectedTeamOK, nil, expectedBoxOK, nil},
		{"call fails", expectedTeamOK, nil, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockUserService := new(interfaceMocks.UserService)
			mockTeamService := new(interfaceMocks.TeamService)
			if tc.findTeamErr == nil && tc.expectedTeam == expectedTeamOK {
				mockRepository.On("Store", mock.Anything).Return(bson.NewObjectId(), tc.err)
			}

			mockUserService.On("OnAllServicesRegistered", mock.Anything)
			mockTeamService.On("OnAllServicesRegistered", mock.Anything)
			mockTeamService.On("FindByID", teamID).Return(tc.expectedTeam, tc.findTeamErr)

			r := registry.NewRegistry()
			r.RegisterService(mockTeamService, "team")
			r.RegisterService(mockUserService, "user")
			service := box.NewService(mockRepository, r)
			r.AllServicesRegistered()
			result, err := service.CreateBox(userID, teamID, boxName, openDate)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestInsertNote(t *testing.T) {
	t.Parallel()
	userID, boxID, boxName, openDate := bson.NewObjectId(), bson.NewObjectId(), "box", time.Now().Add(time.Duration(1*time.Minute))
	user := &models.User{}
	user.SetId(userID)
	note := &models.Note{From: userID, Text: "this is a test"}
	team := &models.Team{Members: []*models.User{user}}
	badTeam := &models.Team{Members: []*models.User{}}
	expectedBoxOK := &models.Box{Name: boxName, OpenDate: openDate, Notes: []*models.Note{note}, Team: team}
	openBox := &models.Box{Name: boxName, OpenDate: time.Now(), Notes: []*models.Note{note}, Team: team}
	var testCases = []struct {
		testName    string
		expectedBox *models.Box
		findBoxErr  error
		addNoteErr  error
		expected    *models.Box
		err         error
	}{
		{"box does not exist", nil, errors.New("fail"), nil, nil, errors.New("fail")},
		{"user not in team", &models.Box{Team: badTeam}, nil, nil, nil, fmt.Errorf("you are not in the team %s so you can't add notes", badTeam.GetId())},
		{"add note fails", openBox, nil, nil, nil, errors.New("can't add note because box is open")},
		{"call ok", expectedBoxOK, nil, nil, expectedBoxOK, nil},
		{"call fails", expectedBoxOK, nil, nil, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockUserService := new(interfaceMocks.UserService)
			mockTeamService := new(interfaceMocks.TeamService)
			if tc.findBoxErr == nil && tc.expectedBox == expectedBoxOK {
				mockRepository.On("Store", mock.Anything).Return(bson.NewObjectId(), tc.err)
			}

			mockRepository.On("FindByID", mock.AnythingOfType("bson.ObjectId")).Return(tc.expectedBox, tc.findBoxErr)
			mockUserService.On("OnAllServicesRegistered", mock.Anything)
			mockTeamService.On("OnAllServicesRegistered", mock.Anything)

			r := registry.NewRegistry()
			r.RegisterService(mockTeamService, "team")
			r.RegisterService(mockUserService, "user")
			service := box.NewService(mockRepository, r)
			r.AllServicesRegistered()
			result, err := service.InsertNote(userID, boxID, "this is a test")

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestGetNotes(t *testing.T) {
	userID := bson.NewObjectId()
	user := &models.User{}
	user.SetId(userID)

	var myNotes, otherNotes, allNotes []*models.Note
	resetNotes := func() {
		myNotes = []*models.Note{&models.Note{Text: "test1", From: userID}, &models.Note{Text: "test2", From: userID}}
		otherNotes = []*models.Note{&models.Note{Text: "test3", From: bson.NewObjectId()}, &models.Note{Text: "test2", From: bson.NewObjectId()}}
		allNotes = append(myNotes, otherNotes...)
	}

	openBox := &models.Box{OpenDate: time.Now()}
	closedBox := &models.Box{OpenDate: time.Now().Add(time.Minute * 60)}

	var pointerToNil []*models.Note

	var testCases = []struct {
		testName string
		box      *models.Box
		notes    *[]*models.Note
		expected *[]*models.Note
		err      error
	}{
		{"no notes", &models.Box{Notes: []*models.Note{}}, &[]*models.Note{}, &[]*models.Note{}, nil},
		{"all notes because open", openBox, &allNotes, &allNotes, nil},
		{"only my notes because closed", closedBox, &allNotes, &myNotes, nil},
		{"call fails", &models.Box{Notes: []*models.Note{}}, &allNotes, &pointerToNil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			resetNotes()
			tc.box.Notes = *tc.notes
			mockRepository := new(mocks.Repository)
			mockUserService := new(interfaceMocks.UserService)
			mockTeamService := new(interfaceMocks.TeamService)
			mockUserService.On("OnAllServicesRegistered", mock.Anything)
			mockTeamService.On("OnAllServicesRegistered", mock.Anything)
			if len(*tc.expected) != 0 {
				mockUserService.On("FindByID", mock.AnythingOfType("bson.ObjectId")).Return(user, nil)
			}

			if tc.err != nil {
				mockUserService.On("FindByID", mock.AnythingOfType("bson.ObjectId")).Return(nil, tc.err)
			}

			r := registry.NewRegistry()
			r.RegisterService(mockTeamService, "team")
			r.RegisterService(mockUserService, "user")
			service := box.NewService(mockRepository, r)
			r.AllServicesRegistered()
			result, err := service.GetNotes(userID, tc.box)

			assert.Equal(t, *tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}

}

func TestOnAllServicesRegistered(t *testing.T) {
	t.Parallel()
	mockRepository := new(mocks.Repository)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, "box")
	mockRegistry.On("GetService", "team").Return(new(interfaceMocks.TeamService))
	mockRegistry.On("GetService", "user").Return(new(interfaceMocks.UserService))

	s := box.NewService(mockRepository, mockRegistry)
	s.OnAllServicesRegistered(mockRegistry)

	mockRegistry.AssertExpectations(t)

}

/*


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
*/
