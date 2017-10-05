package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

// ListNotesHandler handles GET requests for a box's notes
func ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	box, err := getBox(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusNotFound)
	}

	notes := box.GetNotes()
	utils.ResponseJSON(w, notes, true)
}

// InsertNoteHandler handles PUT requests for inserting a note in a box
func InsertNoteHandler(w http.ResponseWriter, r *http.Request) {
	box, err := getBox(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusNotFound)
	}

	note := models.NewNote()
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := note.Validate(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	box.AddNote(note)

}
