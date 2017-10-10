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
	utils.ResponseJSON(w, models.GetBoxListResponse(), true)
}

// CreateBoxHandler handles POST requests for box creation
func CreateBoxHandler(w http.ResponseWriter, r *http.Request) {
	var boxRequest models.BoxRequest
	if err := json.NewDecoder(r.Body).Decode(&boxRequest); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return

	}
	box := models.NewBox(boxRequest)
	if err := box.Save(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
	} else {
		setLocationHeader(w, r, box)
		utils.ResponseCreated(w)
	}
}

// BoxDetailHandler handles GET requests for box detail
func BoxDetailHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	box.RefreshStatus()
	utils.ResponseJSON(w, box.GetResponse(), false)
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
