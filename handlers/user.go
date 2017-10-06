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

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.Save(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
	} else {
		utils.ResponseCreated(w)
	}
}
