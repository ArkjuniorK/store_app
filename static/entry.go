// ==================
// Entry file that would be manage the route for static file/assets.
// This entry file would not using models and controllers, it only serving
// static file to client without read/write data.
// ==================

package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// struct to hold function for serving static file
type Entry struct{}

//
func (e Entry) Routes() chi.Router {
	// init new chi router
	r := chi.NewRouter()

	// set working directory
	workdir, _ := os.Getwd()

	// handle static assets for "/cats"
	r.Get("/cats/*", func(w http.ResponseWriter, r *http.Request) {
		filedir := http.Dir(filepath.Join(workdir, "static/cats"))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filedir))
		fs.ServeHTTP(w, r)
	})

	// return the router
	return r
}
