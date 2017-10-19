package models

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/go-bongo/bongo"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *bongo.Collection

// UserStatus is a string that
type userStatus string

const (
	userActive   = userStatus("ACTIVE")
	userInactive = userStatus("INACTIVE")
)

// User is a document which holds information about a user
type User struct {
	bongo.DocumentBase `bson:",inline"`
	Username           string     `bson:"username"`
	Password           string     `bson:"password"`
	Email              string     `bson:"email"`
	FirstName          string     `bson:"firstName"`
	LastName           string     `bson:"lastName"`
	Status             userStatus `bson:"status"`
}

// UserRequest is a struct that resembles a request performed by users to edit or create a user
type UserRequest struct {
	Username  string  `json:"username"`
	Password  *string `json:"password,omitempty"`
	Email     string  `json:"email"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
}

//UserResponse is a struct that resembles a response for user detail and listing
type UserResponse struct {
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Status    userStatus    `json:"status"`
	ID        bson.ObjectId `json:"id"`
}

// UserList is a list of User Documents
type UserList []User

// UserListResponse is a list of User Documents
type UserListResponse []UserResponse

func validatePassword(password string) error {
	if password == "" {
		return errors.New("Password is required")
	}
	if len(password) < 8 {
		return errors.New("Password must have at least 8 characters")
	}

	return nil
}

// NewUser returns an User instance, with status set to inactive
func NewUser(request UserRequest) (*User, error) {
	user := &User{
		Status:    userInactive,
		Username:  request.Username,
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}
	if request.Password != nil {
		if err := validatePassword(*request.Password); err != nil {
			return user, err
		}
		user.SetPassword(*request.Password)
	}
	return user, nil
}

// GetUserByEmail return an user from database if the email exists
func GetUserByEmail(email string) (User, error) {
	user := User{}
	err := userCollection.FindOne(bson.M{"email": email}, &user)
	return user, err

}

// GetUserByUsername return an user from database if the email exists
func GetUserByUsername(username string) (User, error) {
	user := User{}
	err := userCollection.FindOne(bson.M{"username": username}, &user)

	return user, err

}

// GetUserByID return an user from database if an user with the specified ID exists.
func GetUserByID(id string) (user User, err error) {
	if !bson.IsObjectIdHex(id) {
		return user, fmt.Errorf("%s is not a valid id}", id)
	}

	err = userCollection.FindById(bson.ObjectIdHex(id), &user)
	if err != nil {
		if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
			return user, dnfError
		}
		log.Panic("WTF", err.Error())
	}
	return
}

func (u *User) String() string {
	return fmt.Sprintf("User: %q id %s email %q", u.Username, u.Id, u.Email)
}

// Save saves a User instance into database
func (u *User) Save() error {
	if err := u.validate(); err != nil {
		return err
	}
	u.Status = userActive
	return userCollection.Save(u)
}

func (u *User) validate() error {
	if err := u.validateUsername(); err != nil {
		return err
	}
	if u.FirstName == "" {
		return errors.New("Field firstName is required")
	}
	if err := u.validateEmail(); err != nil {
		return err
	}
	return validatePassword(u.Password)
}

func (u *User) validateEmail() error {
	emailRegexp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if u.Email == "" {
		return errors.New("Field email is required")
	}
	if !emailRegexp.MatchString(u.Email) {
		return errors.New("Invalid email format")
	}

	if otherU, err := GetUserByEmail(u.Email); err == nil && u.GetId() != otherU.GetId() { //ensure unique email
		return errors.New("Email already exists")
	}
	return nil
}

func (u *User) validateUsername() error {
	if u.Username == "" {
		return errors.New("Field username is required")
	}
	if otherU, err := GetUserByUsername(u.Username); err == nil && u.GetId() != otherU.GetId() { //ensure unique email
		return errors.New("Username already exists")
	}
	return nil
}

// SetPassword sets the privided password to the user, but using PBKDF2 cypher
func (u *User) SetPassword(password string) {
	dk := getPBKDF2([]byte(password))
	u.Password = base64.StdEncoding.EncodeToString(dk)
}

func getPBKDF2(passphrase []byte) []byte {
	salt := []byte("MagicBox")
	iterations := 4096
	keylen := 64

	return pbkdf2.Key([]byte(passphrase), salt, iterations, keylen, sha512.New)
}

// Delete deletes a box instance from database
func (u *User) Delete() error {
	return userCollection.DeleteDocument(u)
}

// Update updates a User instance from database
func (u *User) Update(request UserRequest) error {
	u.Username = request.Username
	u.Email = request.Email
	u.FirstName = request.FirstName
	u.LastName = request.LastName
	log.Println(request.Password)
	log.Println(*request.Password)
	log.Println(u.Password)
	if request.Password != nil {
		if err := validatePassword(*request.Password); err != nil {
			return err
		}
		u.SetPassword(*request.Password)
	}

	return u.Save()
}

// GetResponse returns a BoxResponse
func (u *User) GetResponse() UserResponse {
	response := UserResponse{
		Username:  u.Username,
		Status:    u.Status,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ID:        u.GetId(),
	}
	return response
}

// ChallengePassword returns whether a provided password equals the user's passowrd
func (u *User) ChallengePassword(password string) bool {
	ciphered := getPBKDF2([]byte(password))
	return base64.StdEncoding.EncodeToString(ciphered) == u.Password
}

func newUserList() UserList {
	return make([]User, 0)
}

//ListUsers returns all boxes in the box collection
func ListUsers() (users UserList) {
	users = newUserList()
	results := userCollection.Find(bson.M{})

	user := User{}
	for results.Next(&user) {
		users = append(users, user)
	}

	return
}

//GetUserListResponse returns a UserListResponse which represent a the users in the database
func GetUserListResponse() UserListResponse {
	users := ListUsers()
	responses := make(UserListResponse, len(users))
	for i, user := range users {
		responses[i] = user.GetResponse()
	}
	return responses
}
