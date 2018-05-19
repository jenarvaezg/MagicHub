package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/zebresel-com/mongodm"
)

// Status is the status of a box, open or closed
type Status string

const (
	// BoxStatusOpen is the status of an open box
	BoxStatusOpen = Status("open")
	// BoxStatusClosed is the status of a closed box
	BoxStatusClosed = Status("closed")
)

// Box is a document which holds information about a box
type Box struct {
	mongodm.DocumentBase `bson:",inline"`

	Name     string      `json:"name" bson:"name"`
	Notes    []Note      `json:"notes" bson:"notes"`
	Status   Status      `json:"status" bson:"status"`
	OpenDate time.Time   `json:"openDate" bson:"openDate"`
	Team     interface{} `json:"box" bson:"team" model:"Team" relation:"11" autosave:"true"`
}

func (b *Box) String() string {
	return fmt.Sprintf("Box name: %q notes %s, opens at %s", b.Name, b.Notes, b.OpenDate)
}

// RefreshStatus updates the box status if conditions are met
func (b *Box) RefreshStatus() {
	if time.Now().After(b.OpenDate) {
		b.Status = BoxStatusOpen
	}
}

// AddNote adds a note to a box
func (b *Box) AddNote(note Note) error {
	if b.Status == BoxStatusOpen {
		return errors.New("Only closed boxes can get new notes")
	}
	b.Notes = append(b.Notes, note)
	return b.Save()
}
