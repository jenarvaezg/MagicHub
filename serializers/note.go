package serializers

import (
	"encoding/json"

	"github.com/jenarvaezg/magicbox/models"
)

type notesSerializer struct {
	Results models.Notes `json:"results"`
}

// SerializeNotes returns a serialized json as a byte slice for a BoxList
func SerializeNotes(notesSerializable Serializable) ([]byte, error) {
	notes := notesSerializable.(models.Notes)
	serialized, err := json.MarshalIndent(notesSerializer{Results: notes}, "", identation)

	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}
