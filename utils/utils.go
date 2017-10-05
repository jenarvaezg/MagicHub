package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSONMap is an alias to map[string]*json.RawMessage
type JSONMap = map[string]*json.RawMessage

type listSerializer struct {
	Results interface{} `json:"results"`
}

//RemoveForbiddenFields removes id created_at and modified at from JSONMap
func RemoveForbiddenFields(jm JSONMap) JSONMap {
	delete(jm, "_id")
	delete(jm, "_created_at")
	delete(jm, "_updated_at")
	return jm
}

func getJSONEncoder(w http.ResponseWriter) *json.Encoder {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder
}

// ResponseError returns writes to w a mesage and sets the status code to code
func ResponseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	getJSONEncoder(w).Encode(map[string]string{"error": message})
}

// ResponseJSON serializes a object and sends the result to w
func ResponseJSON(w http.ResponseWriter, object interface{}, many bool) { //serializer serializerFn, object interface{}) {
	var err error
	encoder := getJSONEncoder(w)
	if many {
		err = encoder.Encode(listSerializer{Results: object})
	} else {
		err = encoder.Encode(object)
	}
	if err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
}

// GetJSONMap returns a JSONMap which is a map of string to *json.RawMessage
func GetJSONMap(r io.Reader) (objmap JSONMap, err error) {
	err = json.NewDecoder(r).Decode(&objmap)
	return objmap, err
}

// ResponseCreated sets header to 201 Created
func ResponseCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

// ResponseNoContent sets header to 204 NoContent
func ResponseNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
