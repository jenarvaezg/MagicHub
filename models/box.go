package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

var boxCollection *bongo.Collection

// BoxStatus is a string that determines the box' state
type BoxStatus string

const (
	boxStatusOpen   = BoxStatus("open")
	boxStatusClosed = BoxStatus("closed")
)

// Box is a document which holds information about a box
type Box struct {
	bongo.DocumentBase `bson:",inline"`
	Name               string          `bson:"name"`
	Notes              []Note          `bson:"notes"`
	Users              []bson.ObjectId `bson:"users"`
	Status             BoxStatus       `bson:"status"`
	OpenDate           time.Time       `bson:"openDate"`
}

//BoxResponse is a struct that resembles a response for box detail and listing
type BoxResponse struct {
	Name          string        `json:"name"`
	Status        BoxStatus     `json:"status"`
	OpenDate      time.Time     `json:"openDate"`
	NumberOfNotes int           `json:"numberOfNotes"`
	ID            bson.ObjectId `json:"id"`
}

// BoxRequest is a struct that resembles a request performed by users to edit or create a box instance
type BoxRequest struct {
	Name     string    `json:"name"`
	OpenDate time.Time `json:"openDate"`
}

// BoxList is a list of Box Documents
type BoxList = []Box

// BoxListResponse is a list of BoxResponse
type BoxListResponse = []BoxResponse

// NewBox returns a pointer to a new instance of Box
func NewBox(request BoxRequest) *Box {
	box := &Box{Status: boxStatusClosed}
	box.Notes = make(Notes, 0)
	box.Name = request.Name
	box.OpenDate = request.OpenDate
	return box
}

func (b *Box) String() string {
	return fmt.Sprintf("Box name: %q notes %s, opens at %s", b.Name, b.Notes, b.OpenDate)
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
	b.RefreshStatus()
	return boxCollection.Save(b)
}

// RefreshStatus updates the box status if conditions are met
func (b *Box) RefreshStatus() {
	if time.Now().After(b.OpenDate) {
		b.Status = boxStatusOpen
	}
}

// Delete deletes a box instance from database
func (b *Box) Delete() error {
	return boxCollection.DeleteDocument(b)
}

// Update updates a box instance from database
func (b *Box) Update(request BoxRequest) error {
	b.Name = request.Name
	b.OpenDate = request.OpenDate

	return b.Save()
}

// AddNote adds a note to a box
func (b *Box) AddNote(note Note) error {
	if b.Status == boxStatusOpen {
		return errors.New("Only closed boxes can get new notes")
	}
	b.Notes = append(b.Notes, note)
	return b.Save()
}

// GetNotes returns a list of notes from a Box instance
func (b *Box) GetNotes() (Notes, error) {
	log.Println(b.Status)
	if b.Status != boxStatusOpen {
		return Notes{}, errors.New("Can't get notes from a closed box")
	}
	return b.Notes, nil
}

// DeleteNotes deletes all the notes inside a box
func (b *Box) DeleteNotes() {
	b.Notes = make(Notes, 0)
	b.Save()
}

// GetResponse returns a BoxResponse
func (b *Box) GetResponse() BoxResponse {
	response := BoxResponse{
		Name:          b.Name,
		Status:        b.Status,
		OpenDate:      b.OpenDate,
		NumberOfNotes: len(b.Notes),
		ID:            b.GetId(),
	}
	return response
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

//GetBoxListResponse returns a BoxListResponse which represent a the boxes in the database
func GetBoxListResponse() BoxListResponse {
	boxes := ListBoxes()
	responses := make(BoxListResponse, len(boxes))
	for i, box := range boxes {
		box.RefreshStatus()
		responses[i] = box.GetResponse()
	}
	return responses
}
