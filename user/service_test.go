package user_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"

	interfaceMocks "github.com/jenarvaezg/MagicHub/interfaces/mocks"
	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/registry"
	"github.com/jenarvaezg/MagicHub/user"
	"github.com/jenarvaezg/MagicHub/user/mocks"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	repo := new(mocks.Repository)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, mock.AnythingOfType("string")).Return()

	s := user.NewService(repo, mockRegistry)
	mockRegistry.AssertCalled(t, "RegisterService", s, "user")
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	var testCases = []struct {
		testName  string
		id        string
		callsRepo bool
		expected  *models.User
		err       error
	}{
		{"call ok", bson.NewObjectId().Hex(), true, &models.User{}, nil},
		{"bad objectid", "This is a string", false, nil, errors.New("This is a string is not a valid ID")},
		{"call fails", bson.NewObjectId().Hex(), true, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			if tc.callsRepo {
				mockRepository.On("FindByID", bson.ObjectIdHex(tc.id)).Return(tc.expected, tc.err)
			}
			r := registry.NewRegistry()

			service := user.NewService(mockRepository, r)
			result, err := service.FindByID(tc.id)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestFindByEmail(t *testing.T) {
	t.Parallel()
	email := "jenarvaezg@gmail.com"
	u := &models.User{Email: email}
	var testCases = []struct {
		testName  string
		email     string
		mockUsers []*models.User
		expected  *models.User
		mockError error
		err       error
	}{
		{"call ok", email, []*models.User{u}, u, nil, nil},
		{"call fails", email, []*models.User{}, nil, errors.New("fail"), errors.New("find user by email: fail")},
		{"user not found", email, []*models.User{}, nil, nil, fmt.Errorf("user with email %s not found", email)},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			findMap := map[string]interface{}{"email": tc.email}
			mockRepository := new(mocks.Repository)
			mockRepository.On("FindBy", findMap).Return(tc.mockUsers, tc.mockError)

			r := registry.NewRegistry()
			service := user.NewService(mockRepository, r)
			result, err := service.FindByEmail(tc.email)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	username, email, firstName, lastName, imageURL := "uname", "a@a.com", "first", "last", "google.com"
	expectedUserOK := &models.User{Username: username, Email: email, FirstName: firstName, LastName: lastName, ImageURL: imageURL}
	var testCases = []struct {
		testName  string
		username  string
		email     string
		firstName string
		lastName  string
		imageURL  string
		expected  *models.User
		err       error
	}{
		{"call ok", username, email, firstName, lastName, imageURL, expectedUserOK, nil},
		{"call fails", username, email, firstName, lastName, imageURL, nil, errors.New("fail")},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepository := new(mocks.Repository)
			mockRepository.On("Store", mock.Anything).Return(bson.NewObjectId(), tc.err)
			r := registry.NewRegistry()

			service := user.NewService(mockRepository, r)
			result, err := service.CreateUser(tc.username, tc.email, tc.firstName, tc.lastName, tc.imageURL)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
			mockRepository.AssertExpectations(t)
		})
	}
}

func TestOnAllServicesRegistered(t *testing.T) {
	t.Parallel()
	mockRepository := new(mocks.Repository)
	mockRegistry := new(interfaceMocks.Registry)
	mockRegistry.On("RegisterService", mock.Anything, "user")

	s := user.NewService(mockRepository, mockRegistry)
	s.OnAllServicesRegistered(mockRegistry)

	mockRegistry.AssertCalled(t, "RegisterService", s, "user")

}
