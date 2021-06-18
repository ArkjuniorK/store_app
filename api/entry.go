// ==================
// Entry file for route. Here all the api routes would be handled
// each subroute of api would be defined by Route function
//
// ex: to add "/sample" route to api.
// r.Route("/sample", Sample)
// Sample would be a function that wrapped all the "/sample" endpoint
// ==================

package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// type to hold the Routes() function
type Entry struct{}

func (e Entry) Routes() chi.Router {
	// initiate new chi instance
	r := chi.NewRouter()

	// entry
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "Welcome to API")
	})

	// Route for cats endpoint
	r.Route("/cats", Cats)

	// return the route so main file could mounted it
	return r
}
