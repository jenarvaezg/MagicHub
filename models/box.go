package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-bongo/bongo"
	"github.com/jenarvaezg/magicbox/utils"
	"gopkg.in/mgo.v2/bson"
)

var boxCollection *bongo.Collection

// Box is a document which holds information about a box
type Box struct {
	bongo.DocumentBase `bson:",inline"`
	Name               string `json:"name"`
	Notes              []Note `json:"-" bson:"notes"`
}

// BoxList is a list of Box Documents
type BoxList = []Box

// NewBox returns a pointer to a new instance of Box
func NewBox() *Box {
	box := &Box{}
	box.Notes = make(Notes, 0)
	return box
}

func (b *Box) String() string {
	return fmt.Sprintf("Box name: %q id %s notes %s", b.Name, b.Id, b.Notes)
}

func (b *Box) validate() error {
	if b.Name == "" {
		return errors.New("No name provided")
	}
	return nil
}

// Save saves a Box instance into database
func (b *Box) Save() error {
	if err := b.validate(); err != nil {
		return err
	}
	return boxCollection.Save(b)
}

// Delete deletes a box instance from database
func (b *Box) Delete() error {
	return boxCollection.DeleteDocument(b)
}

// Update updates a box instance from database
func (b *Box) Update(updateMap utils.JSONMap) error {
	updateBytes, err := json.Marshal(updateMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(updateBytes, b)
	if err != nil {
		return err
	}

	return b.Save()
}

// AddNote adds a note to a box
func (b *Box) AddNote(note Note) error {
	b.Notes = append(b.Notes, note)
	return b.Save()
}

// GetNotes returns a list of notes from a Box instance
func (b *Box) GetNotes() Notes {
	return b.Notes
}

// DeleteNotes deletes all the notes inside a box
func (b *Box) DeleteNotes() {
	b.Notes = make(Notes, 0)
	b.Save()
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

func newBoxList() BoxList {
	return make([]Box, 0)
}

//ListBoxes returns all boxes in the box collection
func ListBoxes() (boxes BoxList) {
	boxes = newBoxList()
	results := boxCollection.Find(bson.M{})

	box := Box{}
	for results.Next(&box) {
		boxes = append(boxes, box)
	}

	return
}
