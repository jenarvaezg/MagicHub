package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

// ListUsersHandler handles GET requests for listing users in database
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	//TODO filtering
	users := models.ListUsers()
	utils.ResponseJSON(w, users, true)

}

// CreateUserHandler handles POST requests for user creation
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Password == nil {
		utils.ResponseError(w, "Missing password field", http.StatusBadRequest)
		return
	}
	user.SetPassword(*user.Password)

	if err := user.Save(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
	} else {
		setLocationHeader(w, r, &user)
		utils.ResponseCreated(w)
	}
}

// UserDetailHandler handles GET requests for user detail
func UserDetailHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	user.Password = nil
	utils.ResponseJSON(w, user, false)

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

	jsonMap, err := utils.GetJSONMap(r.Body)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
	}

	if err := user.Update(jsonMap); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ResponseNoContent(w)
}
