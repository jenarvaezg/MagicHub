package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jenarvaezg/MagicHub/models"
	"github.com/jenarvaezg/MagicHub/utils"
)

func getUserRequest(r *http.Request) (models.UserRequest, error) {
	var userRequest models.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	return userRequest, err
}

// ListUsersHandler handles GET requests for listing users in database
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	//TODO filtering
	utils.ResponseJSON(w, models.GetUserListResponse(), true)

}

// CreateUserHandler handles POST requests for user creation
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userRequest, err := getUserRequest(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := models.NewUser(userRequest)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user.Save(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
	} else {
		setLocationHeader(w, r, user)
		utils.ResponseCreated(w)
	}
}

// UserDetailHandler handles GET requests for user detail
func UserDetailHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	utils.ResponseJSON(w, user.GetResponse(), false)

}

// UserDeleteHandler handles GET requests for user detail
func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if err := user.Delete(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponseNoContent(w)
}

// UserPatchHandler handles PATCH requests for user updating
func UserPatchHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	userRequest, err := getUserRequest(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.Update(userRequest); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponseNoContent(w)
}
