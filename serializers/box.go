package serializers

import (
	"encoding/json"

	"github.com/jenarvaezg/magicbox/models"
)

type boxListSerializer struct {
	Results models.BoxList `json:"results"`
}

// Serializable is an alias for interface{}
type Serializable interface {
}

const identation string = "  "

// SerializeBoxList returns a seriaized json as a byte slice for a BoxList
func SerializeBoxList(boxesInterface Serializable) ([]byte, error) {
	boxes := boxesInterface.(models.BoxList)
	serialized, err := json.MarshalIndent(boxListSerializer{Results: boxes}, "", identation)

	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}
