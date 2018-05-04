package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jenarvaezg/MagicHub/auth"
	"github.com/jenarvaezg/MagicHub/utils"
)

func loginWithOwnUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token, err := auth.GetAuthTokenFromForm(r.Form)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func loginWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req auth.GoogleFrontendRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := auth.GetAuthTokenFromGoogleToken(req)
	if err != nil {
		log.Println(err)
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

// LoginRequestHandler handles request for token issuing
func LoginRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Google-Login") != "" {
		loginWithGoogle(w, r)
	} else {
		loginWithOwnUser(w, r)
	}
}
