package models

import (
	"fmt"
)

// Note is an embedded document which holds information about a note
type Note struct {
	Title          string `bson:"title"`
	Detail         string `bson:"detail"`
	MentionedUsers User   `bson:"mentionedUsers"`
}

func (n Note) String() string {
	return fmt.Sprintf("Note: %s. Detail: %s", n.Title, n.Detail)
}
