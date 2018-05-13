package box

import (
	"fmt"

	"github.com/jenarvaezg/MagicHub/user"
)

// Note is an embedded document which holds information about a note
type Note struct {
	Title          string    `bson:"title"`
	Detail         string    `bson:"detail"`
	MentionedUsers user.User `bson:"mentionedUsers"`
}

func (n Note) String() string {
	return fmt.Sprintf("Note: %s. Detail: %s", n.Title, n.Detail)
}
