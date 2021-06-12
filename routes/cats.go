// =======================
// This package is package to store routes for cats
// each routes would have their own controller which
// would be imported from the controllers package
// =======================

package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/ArkjuniorK/store_app/controllers"
)

// define controller
var cat = controllers.Cat{}

// Cats router function that would be exported to main.go
// and used by "/cats" endpoint
func Cats(r chi.Router) {
	r.Get("/", cat.GetCats)
	r.Post("/", cat.AddCat)
	r.Get("/{id}", cat.GetCat)
	r.Put("/{id}", cat.UpdateCat)
}
