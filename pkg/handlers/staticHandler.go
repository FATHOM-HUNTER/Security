package handlers

import (
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/markbates/pkger"

	"github.com/gorilla/mux"
)

// HandleStatic serves static files on route /
func HandleStatic(mux *mux.Router) {
	mux.PathPrefix("/").HandlerFunc(serveStatic)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	ext := path.Ext(r.URL.Path)
	if ext == ".html" || ext == "" {
		file, err := pkger.Open("/static/index.html")
		if err != nil {
			panic(err)
		}
		io.Copy(w, file)
	} else {
		file, err := pkger.Open(filepath.Join("/static", r.URL.Path))
		if err != nil {
			log.Println(r.URL.Path, err)
			return
		}

		if ext == ".js" {
			w.Header().Set("Content-Type", "application/javascript")
		}

		if ext == ".css" {
			w.Header().Set("Content-Type", "text/css")
		}

		if ext == ".svg" {
			w.Header().Set("Content-Type", "image/svg+xml")
		}

		if ext == ".ico" {
			w.Header().Set("Content-Type", "image/x-icon")
		}

		io.Copy(w, file)
	}
}
