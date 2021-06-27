package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/ArkjuniorK/store_app/api"
	"github.com/ArkjuniorK/store_app/static"
)

func main() {

	// define the router
	r := chi.NewRouter()

	// working directory
	// wd, _ := os.Getwd()

	// base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// // base route
	// r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	// 	// render.PlainText(w, r, "Welcome!")
	// 	filedir := http.Dir(filepath.Join(wd, "view"))
	// 	rctx := chi.RouteContext(r.Context())
	// 	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	// 	fs := http.StripPrefix(pathPrefix, http.FileServer(filedir))
	// 	fs.ServeHTTP(w, r)
	// })

	// r.Get("/source/*", func(w http.ResponseWriter, r *http.Request) {
	// 	filedir := http.Dir(filepath.Join(wd, "source"))
	// 	rctx := chi.RouteContext(r.Context())
	// 	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	// 	fs := http.StripPrefix(pathPrefix, http.FileServer(filedir))
	// 	fs.ServeHTTP(w, r)
	// })

	// api endpoints to "/api" endpoint to create more convienent
	// way of managing the endpoint structure, this endpoint would
	// used to access all api request to backend
	r.Mount("/api", api.Entry{}.Routes())

	// static endpoints to "/static" endpoint to manage static assets
	r.Mount("/static", static.Entry{}.Routes())

	// serve the route
	http.ListenAndServe(":3000", r)
}
