package models

import (
	"fmt"
)

// Note is an embedded document which holds information about a note
type Note struct {
	Text string      `json:"text" bson:"text"`
	From interface{} `json:"from" bson:"from" model:"User" relation:"11" autosave:"true"`
}

func (n Note) String() string {
	return fmt.Sprintf("Note: from: %s text: %s", n.From, n.Text)
}
