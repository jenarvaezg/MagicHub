package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jenarvaezg/magicbox/models"
	"github.com/jenarvaezg/magicbox/utils"
)

// RequireJSONFunc is a MatcherFunc for gorilla mux, which specifies that a method is accesed with json
func RequireJSONFunc(r *http.Request, rm *mux.RouteMatch) bool {
	if r.Method == "POST" && r.Header.Get("content-type") != "application/json" {
		rm.Handler = badRequestHandler()
	}
	return true
}

func badRequestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseError(w, "Expected content-type to be application/json", http.StatusBadRequest)
	})
}

func getBox(r *http.Request) (models.Box, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return models.GetBoxByID(id)
}
