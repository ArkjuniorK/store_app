// =====================
// This package is package to store controllers for cats
// it would be used inside the routes package at their own respective route
// note that the controller did not include middleware function
//
// Todo
// - Add pagination for all cats [done]
// - Zip code [done]
// - Filter cat by zip code [done]
// - Search cat by name using query [done]
// - Filter cat by variety using query [done]
// - Filter cat by gender using query
// - Filter cat by age using query
// - Upload multiple image for cat [done]
// - Delete image by image ID [done]
//
// =====================

package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
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
	// Controller to get all cats that would be adopted.
	// With pagination and filters
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
	// initiate cats variable that would be hold cat entities,
	// new would init an empty data inside referenced type
	var (
		cats        *models.Cats  = new(models.Cats)
		catsWrapper []models.Cats // hold models.Cats as chunck
	)

	// get working directory
	wd, err := os.Getwd()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error get working directory"))
		return
	}

	// get params from request
	page, _ := strconv.ParseFloat(chi.URLParam(r, "page"), 64)
	limit, _ := strconv.ParseFloat(chi.URLParam(r, "limit"), 64)

	// query for filtering and searching cats
	// type map
	query := r.URL.Query()

	// read cats directory
	catsDir, err := ioutil.ReadDir(wd + "/data/cats")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading cats data"))
		return
	}

	// total page for pagination
	totalPage := math.Ceil(float64(len(catsDir)) / limit)

	if page > totalPage {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error page is bigger than total page"))
		return
	}

	// first read all the file inside cat's data
	// then set it to *cats
	for i := 0; i < len(catsDir); i++ {

		var cat models.Cat

		// read each file
		catData, err := ioutil.ReadFile(wd + "/data/cats/" + catsDir[i].Name())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error read cat data"))
			return
		}

		// unmarshall the catData
		err = json.Unmarshal(catData, &cat)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error unmarshal cat data"))
			return
		}

		*cats = append(*cats, &cat)
	}

	// second filter the cat by it's zip code
	// the query for zip_code would always present
	// to make sure it easy to find adopt cat by location/region
	catsByRegion, err := func(cats *models.Cats, zip []string) (*models.Cats, error) {
		var catsWrapper models.Cats

		// check the zip_code
		// if it's not present in requested query
		// send an error
		if len(zip[0]) == 0 {
			return nil, errors.New("error zip_code not present in request")
		}

		// process the filter here
		// change the format of zip to int16
		zipCode, err := strconv.ParseInt(zip[0], 0, 16)

		if err != nil {
			return nil, errors.New("error parse zip_code to int")
		}

		// loop the cats data to compare each
		// cat zip_code to requested zip_code
		// then append the filtered cat to catsWrapper
		for i := 0; i < len(*cats); i++ {
			if zipCode == int64((*cats)[i].ZipCode) {
				catsWrapper = append(catsWrapper, (*cats)[i])
			}
		}

		return &catsWrapper, nil
	}(cats, query["zip_code"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// third verify if query["name"] is enabled
	// if it does then search for match cat name
	// but it doesn't return cats param
	catsByName := func(cats *models.Cats, name []string) *models.Cats {
		// name is defined
		// search the cat by given name
		if len(name) != 0 {
			// wrapper for cat name that match with the filter
			var catsWrapper models.Cats

			// set regexp for name
			validName := regexp.MustCompile(name[0])

			// loop the cats to get each detail data
			// to compare the validName with cat name
			for i := 0; i < len(*cats); i++ {
				if validName.MatchString(strings.ToLower((*cats)[i].Name)) {
					catsWrapper = append(catsWrapper, (*cats)[i])
				}
			}

			return &catsWrapper
		}

		// when name is not defined
		return cats
	}(catsByRegion, query["name"])

	// next filter the the catsByName
	// with veriety if it defined/enabled
	catsByVariety := func(cats *models.Cats, variety []string) *models.Cats {
		// variety is defined
		// then filter the cats
		if len(variety) != 0 {
			// wrapper for filtered cats
			var catsWrapper models.Cats

			// loop the cats to find the match variety
			// then append it to catsWrapper so it could be returned
			for i := 0; i < len(*cats); i++ {
				if variety[0] == (*cats)[i].Variety {
					catsWrapper = append(catsWrapper, (*cats)[i])
				}
			}

			return &catsWrapper
		}

		// when variety is not defined
		return cats
	}(catsByName, query["variety"])

	// then filter cats by it's gender if
	// query["gender"] is defined
	catsByGender := func(cats *models.Cats, gender []string) *models.Cats {
		// filter by gender is defined
		if len(gender) != 0 {
			// wrapper for filtered cats
			var catsWrapper models.Cats

			// loop the cats to find the match gender
			// then append it to catsWrapper so it could be returned
			for i := 0; i < len(*cats); i++ {
				if gender[0] == (*cats)[i].Gender {
					catsWrapper = append(catsWrapper, (*cats)[i])
				}
			}

			return &catsWrapper
		}

		// gender is not defined
		return cats
	}(catsByVariety, query["gender"])

	// last filter the cats by age
	// the query takes two value [min, max]
	// so the logic would compare that ages is not more
	// than the min and less than the max number
	catsByAge, err := func(cats *models.Cats, age []string) (*models.Cats, error) {
		// age is defined
		if len(age) != 0 {
			// wrapper for filtered cats
			var catsWrapper models.Cats

			// split the age to get the min and max value
			splittedAge := strings.Split(age[0], ",")

			// change the format of min (string) to int
			min, err := strconv.ParseInt(splittedAge[0], 0, 16)

			if err != nil {
				return nil, errors.New("error parse min to int")
			}

			// change the format of max (string) to int
			max, err := strconv.ParseInt(splittedAge[1], 0, 16)

			if err != nil {
				return nil, errors.New("error parse max to int")
			}

			// loop the cats data to get equivalent
			// age of cats
			for i := 0; i < len(*cats); i++ {
				if min <= int64((*cats)[i].Age) && max >= int64((*cats)[i].Age) {
					catsWrapper = append(catsWrapper, (*cats)[i])
				}
			}

			return &catsWrapper, nil
		}

		// age is not defined
		return cats, nil
	}(catsByGender, query["age"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// loop through cats data and slice it into chunk
	// that would be use to paginate cats data
	for i := 0; i < len(*catsByAge); i += int(limit) {

		// length of data
		length := i + int(limit)

		// when length is bigger than len of data
		// then assign length to len of data
		// to avoid bound of length
		if length > len(*catsByAge) {
			length = len(*catsByAge)
		}

		// we slice the data into chuck and append it
		// to catWrapper
		wrapperChunk := (*catsByAge)[i:length]
		catsWrapper = append(catsWrapper, wrapperChunk)
	}

	// send the response to client
	render.JSON(w, r, catsWrapper[int(page)-1])
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

	// check for zip code
	if cat.ZipCode == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error zip code is required"))
		return
	}

	// check for address
	if cat.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error address is required"))
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
	// get the params
	id := chi.URLParam(r, "id")

	// read data cat based on given id
	data, err := ioutil.ReadFile("data/cats/" + id + ".json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reading cat data"))
		return
	}

	// send struct type data as json to client
	render.JSON(w, r, data)
}

// Controller for update cat entity at /cats/{id} endpoint.
// Response is JSON Object from updated cat
// Accepted methods [PUT]
func (c Cat) UpdateCat(w http.ResponseWriter, r *http.Request) {
	// initiate cat variable
	var (
		mcat  models.CatMap // store from file
		mrcat models.CatMap // store from body
	)

	// get the requested id and body
	id := chi.URLParam(r, "id")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	// unmarshall body and file to map type
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
	// since it could use for loop
	for k := range mcat {
		// check for requested field
		// if it's nil then do not loop
		if mrcat[k] != nil {
			// check if the field in cat and rcat isn't same
			// if condition fulfilled, do the update inside
			if mcat[k] != mrcat[k] {
				mcat[k] = mrcat[k]
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

	// send cat struct to client as json
	render.JSON(w, r, data)
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
