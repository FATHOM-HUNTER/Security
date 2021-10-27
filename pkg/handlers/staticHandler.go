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
			log.Println(err)
		}

		if ext == ".js" {
			w.Header().Set("Content-Type", "application/javascript")
		}

		if ext == ".css" {
			w.Header().Set("Content-Type", "text/css")
		}

		io.Copy(w, file)
	}
}
