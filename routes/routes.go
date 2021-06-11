// main file for routes

package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// type to hold the Routes() function
// type Entry struct{}

func Routes() chi.Router {
	// initiate new chi instance
	r := chi.NewRouter()

	// entry
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "Hello to API")
	})

	// Route for cats endpoint
	r.Route("/cats", Cats)

	// return the route so main file could mounted it
	return r
}
