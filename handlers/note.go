package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

// ListNotesHandler handles GET requests for a box's notes
func ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)

	notes, err := box.GetNotes()
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusConflict) //TODO check
	} else {
		utils.ResponseJSON(w, notes, true)
	}
}

// InsertNoteHandler handles PUT requests for inserting a note in a box
func InsertNoteHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	currentUser := getCurrentUser(r)
	note := models.NewNote()
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := note.Validate(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	note.From = currentUser.GetId().Hex()
	box.AddNote(note)
	utils.ResponseCreated(w)

}

//DeleteNotesHandler handles DELETE requests for deletion of all the notes in the box
func DeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	box.DeleteNotes()
	utils.ResponseNoContent(w)
}
