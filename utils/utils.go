package utils

import (
	"encoding/json"
	"net/http"

	"github.com/jenarvaezg/magicbox/serializers"
)

type serializerFn = func(serializers.Serializable) ([]byte, error)

// ResponseError returns writes to w a mesage and sets the status code to code
func ResponseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ResponseJSON serializes a object and sends the result to w
func ResponseJSON(w http.ResponseWriter, serializer serializerFn, object interface{}) {
	data, err := serializer(object)
	if err != nil {
		ResponseError(w, err.Error(), 400)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}

}
