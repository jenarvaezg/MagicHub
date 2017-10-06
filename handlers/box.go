package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

// ListBoxesHandler handles GET requests for box listing
func ListBoxesHandler(w http.ResponseWriter, r *http.Request) {
	//TODO filtering
	boxes := models.ListBoxes()
	utils.ResponseJSON(w, boxes, true)
}

// CreateBoxHandler handles POST requests for box creation
func CreateBoxHandler(w http.ResponseWriter, r *http.Request) {
	box := models.NewBox()
	if err := json.NewDecoder(r.Body).Decode(&box); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := box.Save(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
	} else {
		utils.ResponseCreated(w)
	}
}

// BoxDetailHandler handles GET requests for box detail
func BoxDetailHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	utils.ResponseJSON(w, box, false)
}

// BoxDeleteHandler handles DELETE requests for box deletion
func BoxDeleteHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)

	if err := box.Delete(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponseNoContent(w)
}

// BoxPatchHandler handles PATCH requests for box updating
func BoxPatchHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)

	jsonMap, err := utils.GetJSONMap(r.Body)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
	}

	if err := box.Update(jsonMap); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ResponseNoContent(w)
}
