// =====================
// This package is package to store controllers for cats
// it would be used inside the routes package at their own respective route
// note that the controller did not include middleware function
//
// Todo
// - Add pagination for all cats
// - Search cat by name
// - Search cat by variety
// - Search cat by age
// - Upload multiple image for cat [done]
// - Delete image by image ID [done]
//
// =====================

package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"

	"github.com/ArkjuniorK/store_app/middleware"
	"github.com/ArkjuniorK/store_app/models"
)

// Define an interface for each cat controllers
// interface would be use inside routes packages to access each controller
// instead of using Cat struct
type CatControllers interface {
	// Controller to get all adopted cats
	GetCats(w http.ResponseWriter, r *http.Request)

	// Controller to add cat to be adopted
	AddCat(w http.ResponseWriter, r *http.Request)

	// Controller to get one cat based on given id
	GetCat(w http.ResponseWriter, r *http.Request)

	// Controller to update one cat based on given id
	UpdateCat(w http.ResponseWriter, r *http.Request)

	// Controller to delete cat based on given id
	DeleteCat(w http.ResponseWriter, r *http.Request)

	// Controller to post cat's image
	UploadImageCat(w http.ResponseWriter, r *http.Request)

	// Controller to delete cat's image
	DeleteImageCat(w http.ResponseWriter, r *http.Request)
}

// define type that would use as pothe controllers of cat
// it would be useful since controllers package would have
// more than one file and we need to wrap controllers that
// would used for specific route, note that type struct is ideal
// to hold functions inside
type Cat string

// Controller for root of "/cats" endpoint.
// Response is JSON Array take from the models.Cats slices.
// Accepted methods [GET]
func (c Cat) GetCats(w http.ResponseWriter, r *http.Request) {
	// initiate cats variable that would be hold cat entities
	// new would init an empty data inside referenced type
	var cats *models.Cats = new(models.Cats)
	var wd, _ = os.Getwd()

	// get params from request
	page, _ := strconv.ParseInt(chi.URLParam(r, "page"), 0, 8)
	limit, _ := strconv.ParseInt(chi.URLParam(r, "limit"), 0, 8)

	// create totalPage
	println(page, limit)

	// read cats.json file
	data, err := ioutil.ReadDir(wd + "/data/cats")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading cats"))
		return
	}

	for _, v := range data {
		var cat *models.Cat

		// read the file
		data, err := ioutil.ReadFile(wd + "/data/cats/" + v.Name())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error reading cat data"))
			return
		}

		// marshall json data inside file to struct
		err = json.Unmarshal(data, &cat)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error unmarshal cat"))
			return
		}

		// append the marshalled json to cats slices
		*cats = append(*cats, cat)
	}

	// send the response to client
	render.JSON(w, r, cats)
}

// Controller for post new cat at "/cats" endpoint.
// Response is JSON Object take from models.Cat struct.
// More specify the res would be the new cat that have been posted.
// Accepted methods [POST]
func (c Cat) AddCat(w http.ResponseWriter, r *http.Request) {
	// initiate cat variabels
	var cat *models.Cat

	// get the body from request
	// since it is in []byte format change it to JSON
	// using unmarshal method from json package
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading cats data"))
		return
	}

	// unmarshal the body to change the format to JSON
	err = json.Unmarshal(body, &cat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error unmarshal requested body"))
		return
	}

	// generate an id for requsted body
	// also with create and update key
	cat.ID = xid.New()
	cat.Create = time.Now()
	cat.Update = time.Now()

	// change the format of requsted body back to JSON
	data, err := json.Marshal(*cat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshal cat"))
		return
	}

	// then write it to file and save it with generated id as filename
	ioutil.WriteFile("data/cats/"+cat.ID.String()+".json", data, 0644)

	// send response
	render.JSON(w, r, cat)
}

// Controller for get cat entity at "/cats/{id}" endpoint.
// Response is JSON Object take from models.Cat struct.
// Accepted methods [GET]
func (c Cat) GetCat(w http.ResponseWriter, r *http.Request) {
	// initate cat variable
	var cat *models.Cat

	// get the params
	id := chi.URLParam(r, "id")

	// read data cat based on given id
	data, err := ioutil.ReadFile("data/cats/" + id + ".json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading cat data"))
		return
	}

	// unmarshall the type of data to struct
	json.Unmarshal(data, &cat)

	// send struct type data as json to client
	render.JSON(w, r, cat)
}

// Controller for update cat entity at /cats/{id} endpoint.
// Response is JSON Object from updated cat
// Accepted methods [PUT]
func (c Cat) UpdateCat(w http.ResponseWriter, r *http.Request) {
	// initiate cat variable
	var (
		mcat  models.CatMap // store from file
		mrcat models.CatMap // store from body
		cat   *models.Cat   // store updated data
	)

	// get the requested id and body
	id := chi.URLParam(r, "id")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading requested body"))
		return
	}

	// find data of cat using id
	file, err := ioutil.ReadFile("data/cats/" + id + ".json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading cat data"))
		return
	}

	// unmarshall body and file to struct type
	if err = json.Unmarshal(body, &mrcat); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error unmarshall requested body"))
		return
	}

	if err = json.Unmarshal(file, &mcat); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error unmarshall cat data"))
		return
	}

	// using map make us easy to compare each field
	// since it could using for loop
	for i, _ := range mcat {
		// check for requested field
		// if it's nil then do not loop
		if mrcat[i] != nil {
			// check if the field in cat and rcat isn't same
			// if condition fulfilled, do the update inside
			if mcat[i] != mrcat[i] {
				mcat[i] = mrcat[i]
			}
		}
	}

	// then update the value of cat Update key
	mcat["updated_at"] = time.Now()

	// chnage the format to []byte
	data, err := json.Marshal(mcat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshall cat data"))
		return
	}

	// write to file
	if err = ioutil.WriteFile("data/cats/"+id+".json", data, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error write updated cat data"))
		return
	}

	// from []byte change again to struct
	// so we could send the response to client
	if err = json.Unmarshal(data, &cat); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshal updated cat data"))
		return
	}

	// send cat struct to client as json
	render.JSON(w, r, cat)
}

// Controller for delete cat entity based on id
// Response is success message
// Accepted methods [DELETE]
func (c Cat) DeleteCat(w http.ResponseWriter, r *http.Request) {
	// get id from params
	id := chi.URLParam(r, "id")

	// find the cat data using id
	err := os.Remove("data/cats/" + id + ".json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error reading cat data"))
		return
	}

	render.PlainText(w, r, "Success deleting cat")
}

// Controller for post cat's image
// Delete image when error is occured
// Response is JSON cat data
// Accepted methods [POST]
func (c Cat) UploadImageCat(w http.ResponseWriter, r *http.Request) {
	var (
		cat   *models.Cat  = new(models.Cat)
		link  *models.Link = new(models.Link)
		wd, _              = os.Getwd()
	)

	// get id from url params
	id := chi.URLParam(r, "id")

	// get the context, since image middleware passing
	// image filename on context so we need to get the value
	filenameCxt := r.Context().Value(middleware.KeyName)

	// change format of filename to string using fmt
	filename := fmt.Sprintf("%v", filenameCxt)

	// read cat data
	catData, err := ioutil.ReadFile("data/cats/" + id + ".json")

	if err != nil {
		// remove image from storage
		err = os.Remove(wd + "/static/cats/" + filename + ".webp")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error delete cat image"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error read cat data"))
		return
	}

	// change cat data to struct
	err = json.Unmarshal(catData, &cat)

	if err != nil {
		// remove image from storage
		err = os.Remove(wd + "/static/cats/" + filename + ".webp")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error delete cat image"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error unmarshal cat data"))
		return
	}

	// assign link
	link.ID = xid.New()
	link.URL = r.Host + "/static/cats/" + filename + ".webp"

	// add image to cat
	// init cat.Image slices first then append link
	cat.Image = new(models.Picture)
	*cat.Image = append(*cat.Image, link)

	// change cat struct to byte
	data, err := json.Marshal(cat)

	if err != nil {
		// remove image from storage
		err = os.Remove(wd + "/static/cats/" + filename + ".webp")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error delete cat image"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error encode cat data to bytes"))
		return
	}

	// write update to file data
	err = ioutil.WriteFile(wd+"/data/cats/"+id+".json", data, 0644)

	if err != nil {
		// remove image from storage
		err = os.Remove(wd + "/static/cats/" + filename + ".webp")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error delete cat image"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error write cat data"))
		return
	}

	// send response
	render.JSON(w, r, cat)
}

// Controller to delete cat's image from data and disk
// Response is JSON cat data
// Accepted methods [DELETE]
func (c Cat) DeleteImageCat(w http.ResponseWriter, r *http.Request) {
	var (
		filename *[]string = new([]string)
		cat      *models.Cat
		id       = chi.URLParam(r, "id")
		id_image = chi.URLParam(r, "id_image")
		wd, _    = os.Getwd()
	)

	// find cat data
	file, err := ioutil.ReadFile(wd + "/data/cats/" + id + ".json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error read cat data"))
		return
	}

	// format cat data to struct
	err = json.Unmarshal(file, &cat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error unmarshal cat file"))
		return
	}

	// delete image data from db by iterating
	// cat.Image slices and find matching ID
	// assign filename of image to filename
	// then delete current index of image using append
	for i, v := range *cat.Image {
		if v.ID.String() == id_image {
			*filename = strings.SplitAfter(v.URL, r.Host)
			*cat.Image = append((*cat.Image)[:i], (*cat.Image)[i+1:]...)
		}
	}

	// change cat to byte
	data, err := json.Marshal(cat)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshal cat data"))
		return
	}

	// save changes by write to file
	err = ioutil.WriteFile(wd+"/data/cats/"+id+".json", data, 0644)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error write cat data"))
		return
	}

	// check if image is written in disk
	_, err = ioutil.ReadFile(wd + (*filename)[1])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error read cat's image"))
		return
	}

	// delete image from disk
	err = os.Remove(wd + (*filename)[1])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error deleting cat's image"))
		return
	}

	render.JSON(w, r, cat)
}
