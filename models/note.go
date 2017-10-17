package models

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2/bson"
)

// Note is an embedded document which holds information about a note
type Note struct {
	From   *bson.ObjectId `bson:"from,omitempty"`
	Title  string         `bson:"title"`
	Detail string         `bson:"detail"`
}

// NoteRequest is a struct that resembles a request performed by users to edit or create a note
type NoteRequest struct {
	Anonymous bool   `json:"anonymous"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
}

// NoteResponse is a struct that resembles a response for note detail and listing
type NoteResponse struct {
	From   *bson.ObjectId `json:"from,omitempty"`
	Title  string         `json:"title"`
	Detail string         `json:"detail"`
}

// NoteListResponse is a list of NoteResponse
type NoteListResponse []NoteResponse

// Notes is a list of Note embedded documents
type Notes []Note

//NewNote returns a Note
func NewNote(request NoteRequest, user User) *Note {
	note := &Note{
		Title:  request.Title,
		Detail: request.Detail,
	}
	if request.Anonymous {
		note.From = nil
	} else {
		userID := user.GetId()
		note.From = &userID
	}
	return note
}

// Validate returns an error if any field is missing
func (n *Note) Validate() error {
	if n.Title == "" {
		return errors.New("Missing title field")
	}
	if n.Detail == "" {
		return errors.New("Missing detail field")
	}
	return nil
}

// GetResponse returns a NoteResponse
func (n *Note) GetResponse() NoteResponse {
	response := NoteResponse{
		Title:  n.Title,
		Detail: n.Detail,
	}
	log.Println(n.From)
	if n.From != nil {
		response.From = n.From
	}
	return response
}

//GetNoteListResponse returns a NoteListResponse which represent a the notes in a box
func GetNoteListResponse(box *Box) (NoteListResponse, error) {
	notes, err := box.GetNotes()
	if err != nil {
		return NoteListResponse{}, err
	}

	responses := make(NoteListResponse, len(notes))
	for i, note := range notes {
		responses[i] = note.GetResponse()
	}
	return responses, nil
}
