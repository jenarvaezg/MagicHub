package models

import (
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2/bson"
)

// Box is a document whic holds information about a box
type Box struct {
	bongo.DocumentBase `bson:",inline"`
	Name               string
}

// BoxList is a list of Box Documents
type BoxList = []*Box

func newBoxList() BoxList {
	return make([]*Box, 0)
}

//ListBoxes returns all boxes in the box collection
func ListBoxes() (boxes BoxList) {
	boxes = newBoxList()
	results := connection.Collection("box").Find(bson.M{})

	box := &Box{}
	for results.Next(box) {
		boxes = append(boxes, box)
	}

	return
}
