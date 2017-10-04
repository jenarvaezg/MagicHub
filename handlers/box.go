package handlers

import (
	"net/http"

	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/serializers"
	"github.com/jenarvaezg/magicbox/utils"
)

//ListBoxesHandler handles requests for box listing
func ListBoxesHandler(w http.ResponseWriter, r *http.Request) {
	//TODO filtering
	boxes := models.ListBoxes()
	utils.ResponseJSON(w, serializers.SerializeBoxList, boxes)
}

func CreateBoxHandler(w http.ResponseWriter, r *http.Request) {
}
