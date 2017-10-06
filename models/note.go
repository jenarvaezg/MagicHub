package models

import "errors"

// Note is an embedded document which holds information about a note
type Note struct {
	From   string `json:"from"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Notes is a list of Note embedded documents
type Notes = []Note

//NewNote returns a Note
func NewNote() Note {
	return Note{}
}

// Validate returns an error if any field is missing
func (n *Note) Validate() error {
	if n.From == "" {
		return errors.New("Missing from field")
	}
	if n.Title == "" {
		return errors.New("Missing title field")
	}
	if n.Detail == "" {
		return errors.New("Missing detail field")
	}
	return nil
}
