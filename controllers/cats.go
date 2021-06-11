// =====================
// This package is package to store controllers for cats
// it would be used inside the routes package at their own respective route
// note that the controller did not include middleware function
// =====================

package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/ArkjuniorK/store_app/models"
)

// define an interface for each cat controllers
type CatControllers interface {
	// Controller to get all adopted cats
	GetCats(w http.ResponseWriter, r *http.Request)

	// Controller to add cat to be adopted
	AddCat(w http.ResponseWriter, r *http.Request)

	// Controller to get one cat based on given id
	GetCat(w http.ResponseWriter, r *http.Request)
}

// define type that would hold all the controllers of cat
// it would be useful since controllers package would have
// more than one file and we need to wrap controllers that
// would used for specific route, note that type struct is ideal
// to hold functions inside
type Cat struct{}

// Controller for root of "/cats" endpoint.
// Response is JSON Array take from the models.Cats slices.
// Accepted methods [GET]
func (c *Cat) GetCats(w http.ResponseWriter, r *http.Request) {
	// initiate cats variable that would be hold cat entities
	// new would init an empty data inside referenced type
	var cats *models.Cats = new(models.Cats)

	// read cats.json file
	// rawCats, err := ioutil.ReadFile("data/cats.json")
	data, err := ioutil.ReadDir("data/cats")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}

	for _, v := range data {
		var cat *models.Cat

		// read the file
		data, err := ioutil.ReadFile("data/cats/" + v.Name())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
		}

		// marshall json data inside file to struct
		err = json.Unmarshal(data, &cat)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
		}

		// append the marshalled json to cats slices
		*cats = append(*cats, cat)
		// fmt.Println(*cats)
	}

	// send the response to client
	render.JSON(w, r, cats)
}

// Controller for post new cat at "/cats" endpoint.
// Response is JSON Object take from models.Cat struct.
// More specify the res would be the new cat that have been posted.
// Accepted methods [POST]
func (c *Cat) AddCat(w http.ResponseWriter, r *http.Request) {
	// initiate cat variabels
	var cat *models.Cat

	// get the body from request
	// since it is in []byte format change it to JSON
	// using unmarshal method from json package
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("bad server"))
	}

	// unmarshal the body to change the format to JSON
	err = json.Unmarshal(body, &cat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}

	// generate an id for requsted body
	cat.ID = uuid.New()

	// change the format of requsted body back to JSON
	data, err := json.Marshal(*cat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}

	// then write it to file and save it with generated id as filename
	ioutil.WriteFile("data/cats/"+cat.ID.String()+".json", data, 0644)

	// send response
	render.JSON(w, r, cat)
}

// Controller for get cat entity at "/cats/{id}" endpoint.
// Response is JSON Object take from models.Cat struct.
// Accepted methods [GET]
func (c *Cat) GetCat(w http.ResponseWriter, r *http.Request) {
	var cat *models.Cat

	// get the params
	id := chi.URLParam(r, "id")

	// read data cat based on given id
	data, err := ioutil.ReadFile("data/cats/" + id + ".json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}

	// unmarshall the type of data to struct
	json.Unmarshal(data, &cat)

	// send struct type data as json to client
	render.JSON(w, r, cat)
}
