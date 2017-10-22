package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

func getNoteRequest(r *http.Request) (models.NoteRequest, error) {
	var noteRequest models.NoteRequest
	err := json.NewDecoder(r.Body).Decode(&noteRequest)
	return noteRequest, err
}

// ListNotesHandler handles GET requests for a box's notes
func ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	user := getCurrentUser(r)

	if !box.IsUserRegistered(user) {
		utils.ResponseError(w, "You are not allowed to get notes from this box", http.StatusForbidden)
		return
	}

	notes, err := models.GetNoteListResponse(box)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusConflict)
	} else {
		utils.ResponseJSON(w, notes, true)
	}
}

// InsertNoteHandler handles PUT requests for inserting a note in a box
func InsertNoteHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	user := getCurrentUser(r)
	if !box.IsUserRegistered(user) {
		utils.ResponseError(w, "You are not allowed to insert notes into this box", http.StatusForbidden)
		return
	}

	noteRequest, err := getNoteRequest(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	note := models.NewNote(noteRequest, user)

	if err := note.Validate(); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := box.AddNote(*note); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(box)
	utils.ResponseCreated(w)
}

//DeleteNotesHandler handles DELETE requests for deletion of all the notes in the box
func DeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	box := getBox(r)
	user := getCurrentUser(r)
	if !box.IsUserRegistered(user) {
		utils.ResponseError(w, "You are not allowed to delete notes from this box", http.StatusForbidden)
		return
	}

	box.DeleteNotes()
	utils.ResponseNoContent(w)
}
