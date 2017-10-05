package serializers

import (
	"encoding/json"

	"github.com/jenarvaezg/magicbox/models"
)

type boxListSerializer struct {
	Results Serializable `json:"results"`
}

// SerializeBoxList returns a serialized json as a byte slice for a BoxList
func SerializeBoxList(boxlistSerializable Serializable) ([]byte, error) {
	boxlist := boxlistSerializable.(models.BoxList)
	serialized, err := json.MarshalIndent(boxListSerializer{Results: boxlist}, "", identation)

	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}

// SerializeBox returns a seriaized json as a byte slice for a Box
func SerializeBox(boxSerializable Serializable) ([]byte, error) {
	box := boxSerializable.(models.Box)
	serialized, err := json.MarshalIndent(box, "", identation)

	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}
