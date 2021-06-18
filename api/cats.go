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
	"fmt"
	"math"
	"net/http"
	"strconv"

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
	r.Get("/test/{page}/{limit}", func(w http.ResponseWriter, r *http.Request) {

		// dummy data
		nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
		var slicedNums []int

		// get params
		// page, _ := strconv.ParseInt(chi.URLParam(r, "page"), 0, 8)
		limit, _ := strconv.ParseFloat(chi.URLParam(r, "limit"), 0)

		// total page
		totalPage := float64(len(nums)) / limit

		for i := 0; i < int(totalPage); i++ {
			fmt.Println("index,", i)
			// for j := 0; j < len(nums); j++ {
			// 	if j < int(limit) {
			// 		slicedNums = append(slicedNums, nums[j])
			// 	}
			// }
		}

		nums = nums[0:3]

		fmt.Println(math.Ceil(totalPage), slicedNums, nums)
	})
}
