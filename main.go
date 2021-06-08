package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Cat struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Gender string `json:"gender"`
}

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

	// Routes for cats
	r.Route("/cats", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			// initiate cats variable
			var cats []*Cat

			// read cats.json file
			rawCats, err := ioutil.ReadFile("data/cats.json")

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("not found"))
			}

			// unmarshall the cats
			// from []byte to json format
			err = json.Unmarshal(rawCats, &cats)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}

			// send response to client
			render.JSON(w, r, cats)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			// initiate cat variable
			var cats []*Cat
			// var cat *Cat

			// get the params
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 8)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			// read cats.json file
			rawCats, err := ioutil.ReadFile("data/cats.json")

			// if cats.json cannot be read
			// send status not found to client
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("not found"))
			}

			// unmarshal the rawCats into json
			// and placed it inside cats slices
			err = json.Unmarshal(rawCats, &cats)

			// if unmarshal error send server error
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}

			// get the matching cat from cats
			// using id as identifier inside for loop
			for _, v := range cats {
				if v.ID == id {
					// get the selected cat/value

					// json.NewEncoder().Encode(v)
					render.JSON(w, r, v)
				}
			}

		})
	})

	// serve the route
	http.ListenAndServe(":3000", r)
}
