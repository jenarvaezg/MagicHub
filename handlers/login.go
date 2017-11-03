package handlers

import (
	"net/http"

	"github.com/jenarvaezg/magicbox/auth"
	"github.com/jenarvaezg/magicbox/utils"
)

// LoginRequestHandler handles request for token issuing
func LoginRequestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token, err := auth.GetAuthToken(r.Form)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))

}
