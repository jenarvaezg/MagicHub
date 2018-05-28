package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jenarvaezg/mongodm"
)

// Status is the status of a box, open or closed
type Status string

// Box is a document which holds information about a box
type Box struct {
	mongodm.DocumentBase `bson:",inline"`

	Name     string      `json:"name" bson:"name"`
	Notes    []*Note     `json:"notes" bson:"notes"`
	OpenDate time.Time   `json:"openDate" bson:"openDate"`
	Team     interface{} `json:"team" bson:"team" model:"Team" relation:"11" autosave:"true"`
}

func (b *Box) String() string {
	return fmt.Sprintf("Box name: %q notes %s, opens at %s", b.Name, b.Notes, b.OpenDate)
}

// IsOpen returns if a box is open
func (b *Box) IsOpen() bool {
	return time.Now().After(b.OpenDate)
}

// AddNote adds a note to a box
func (b *Box) AddNote(note Note) error {
	if b.IsOpen() {
		return errors.New("can't add note because box is open")
	}
	b.Notes = append(b.Notes, &note)
	return nil
}
