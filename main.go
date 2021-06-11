package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/ArkjuniorK/store_app/routes"
)

func main() {

	// define the router
	r := chi.NewRouter()

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

	// base route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "Welcome!")
	})

	// Routes from routes folder, here we actually didn't have to
	// include the bracket at the end of function since the exported func
	// itself already contain the chi.Router params inside
	//
	// dont know if it was a bug or something, I expect to pass an argument
	// inside of the rts functions because it actually a function
	// r.
	// r.Route("/cats", routes.Cats)
	r.Mount("/api", routes.Routes())

	// serve the route
	http.ListenAndServe(":3000", r)
}
