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
	// StatusOpen is the status of an open box
	StatusOpen = Status("open")
	// StatusClosed is the status of a closed box
	StatusClosed = Status("closed")
)

// Box is a document which holds information about a box
type Box struct {
	mongodm.DocumentBase `bson:",inline"`

	Name     string    `bson:"name"`
	Notes    []Note    `bson:"notes"`
	Status   Status    `bson:"status"`
	OpenDate time.Time `bson:"openDate"`
}

func (b *Box) String() string {
	return fmt.Sprintf("Box name: %q notes %s, opens at %s", b.Name, b.Notes, b.OpenDate)
}

// RefreshStatus updates the box status if conditions are met
func (b *Box) RefreshStatus() {
	if time.Now().After(b.OpenDate) {
		b.Status = StatusOpen
	}
}

// AddNote adds a note to a box
func (b *Box) AddNote(note Note) error {
	if b.Status == StatusOpen {
		return errors.New("Only closed boxes can get new notes")
	}
	b.Notes = append(b.Notes, note)
	return b.Save()
}

// // GetBoxByID returns a box searching by id
// func GetBoxByID(id string) (box Box, err error) {
// 	if !bson.IsObjectIdHex(id) {
// 		return box, fmt.Errorf("%s is not a valid id}", id)
// 	}

// 	err = boxCollection.FindById(bson.ObjectIdHex(id), &box)
// 	if err != nil {
// 		log.Panic("WTF", err.Error())
// 	}
// 	return
// }
