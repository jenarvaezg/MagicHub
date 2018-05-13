package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-bongo/bongo"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

var boxCollection *bongo.Collection

type status string

const (
	statusOpen   = status("open")
	statusClosed = status("closed")
)

// Box is a document which holds information about a box
type Box struct {
	mongodm.DocumentBase `bson:",inline"`

	Name     string    `bson:"name"`
	Notes    []Note    `bson:"notes"`
	Status   status    `bson:"status"`
	OpenDate time.Time `bson:"openDate"`
}

func (b *Box) String() string {
	return fmt.Sprintf("Box name: %q notes %s, opens at %s", b.Name, b.Notes, b.OpenDate)
}

// RefreshStatus updates the box status if conditions are met
func (b *Box) RefreshStatus() {
	if time.Now().After(b.OpenDate) {
		b.Status = statusOpen
	}
}

// AddNote adds a note to a box
func (b *Box) AddNote(note Note) error {
	if b.Status == statusOpen {
		return errors.New("Only closed boxes can get new notes")
	}
	b.Notes = append(b.Notes, note)
	return b.Save()
}

// GetBoxByID returns a box searching by id
func GetBoxByID(id string) (box Box, err error) {
	if !bson.IsObjectIdHex(id) {
		return box, fmt.Errorf("%s is not a valid id}", id)
	}

	err = boxCollection.FindById(bson.ObjectIdHex(id), &box)
	if err != nil {
		if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
			return box, dnfError
		}
		log.Panic("WTF", err.Error())
	}
	return
}
