package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-bongo/bongo"
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

func getBox(r *http.Request) *models.Box {
	ctx := r.Context()
	box := ctx.Value(utils.ContextKeyBox).(models.Box)
	box.RefreshStatus()
	return &box
}

func getUser(r *http.Request) models.User {
	ctx := r.Context()
	return ctx.Value(utils.ContextKeyUser).(models.User)
}

func getCurrentUser(r *http.Request) models.User {
	ctx := r.Context()
	return ctx.Value(utils.ContextKeyCurrentUser).(models.User)
}

func setLocationHeader(w http.ResponseWriter, r *http.Request, document bongo.Document) {
	url, _ := mux.CurrentRoute(r).URL()
	id := document.GetId().Hex()
	w.Header().Set("Locaton", fmt.Sprintf("%s/%s", url.RequestURI(), id))
}
