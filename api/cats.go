// =======================
// This package is package to store routes for cats
// each routes would have their own controller which
// would be imported from the controllers package
//
// Todo:
// - Add pagination for all cats
// - Search cat by name
// - Search cat by variety
// - Search cat by age
// - Create custom middleware for uploading image for cat
// - Upload multiple image for cat
// - Delete image by image ID
//
// =======================

package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/ArkjuniorK/store_app/controllers"
	"github.com/ArkjuniorK/store_app/middleware"
)

// define controller
var (
	CatWrapper *controllers.Cat           = new(controllers.Cat)
	Cat        controllers.CatControllers = *CatWrapper
)

// Cats router function that would be exported to main.go
// and used by "/cats" endpoint
func Cats(r chi.Router) {
	r.Get("/{page}/{limit}", Cat.GetCats)
	r.Post("/add", Cat.AddCat)
	r.Get("/{id}", Cat.GetCat)
	r.Put("/{id}", Cat.UpdateCat)
	r.Delete("/{id}", Cat.DeleteCat)
	r.With(middleware.SetImage).Post("/{id}", Cat.UploadImageCat)
	r.Delete("/{id}/{id_image}", Cat.DeleteImageCat)
}
