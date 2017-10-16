package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

func getRegisterRequest(r *http.Request) (models.BoxRegisterRequest, error) {
	var registerRequest models.BoxRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	return registerRequest, err
}

// RegisterInBoxHandler handles POST requests for adding user into a box
func RegisterInBoxHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	registerRequest, err := getRegisterRequest(r) //boxRequest, err := getBoxRequest(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if box.ChallengePassword(registerRequest.Passphrase) {
		utils.ResponseError(w, "Provided passphrase is not valid for this box", http.StatusBadRequest)
		return
	}
	if err := box.AddUser(getCurrentUser(r)); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusConflict)
		return
	}
	box.Save()
	w.WriteHeader(http.StatusOK)

}

// RemoveFromBoxHandler handles DELETE requests for user deletion from a box
func RemoveFromBoxHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)

	if err := box.RemoveUser(getCurrentUser(r)); err != nil {
		utils.ResponseError(w, "Provided passphrase is not valid for this box", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
