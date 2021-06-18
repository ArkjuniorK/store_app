package middleware

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/h2non/bimg"
	"github.com/rs/xid"
)

// How to work:
// - Get the file image form
// - Resize the file image
// - Save it as .png

type Key int

// KeyName is variable with custom type to assign inside
// context so filename could be accessed.
// KeyName would be exported so controller could get the key
// for filename context
const KeyName Key = iota

// Function that act as middleware for file request,
// this middleware would read the requested file, write it to
// static folder as image in webp format
func SetImage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the requset url path
		// dir would be specific directory where image
		// would be saved
		url := r.URL.Path
		dir := strings.Split(url, "/")[2] + "/"

		// generate xid for filename
		// later it would be used to create the ID
		filename := xid.New().String()

		// get the requsted id
		id := chi.URLParam(r, "id")

		if len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error getting param id"))
			return
		}

		// get the requsted body
		err := r.ParseForm()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error parsing form"))
			return
		}

		// get the image form value
		file, _, err := r.FormFile("image")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting form file"))
			return
		}

		buff, err := io.ReadAll(file)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error read form file"))
			return
		}

		// resize image buffer
		buff, err = bimg.Resize(buff, bimg.Options{
			Width:       800,
			Height:      0,
			Quality:     80,
			Compression: 80,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error resize buffer"))
			return
		}

		// write the buffer to file and save it as webp
		err = bimg.Write("static/"+dir+string(filename)+"."+bimg.ImageTypeName(bimg.WEBP), buff)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error write image"))
			return
		}

		ctx := context.WithValue(r.Context(), KeyName, filename)

		// next to controller
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
