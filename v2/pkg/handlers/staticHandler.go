package handlers

import (
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"

	"github.com/gorilla/mux"
)

// HandleStatic serves static files on route /
func HandleStatic(mux *mux.Router) {
	mux.HandleFunc("/", serveStatic)
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
		w.Header().Add("Cache-Control", "max-age=604800000")
		split := strings.Split(r.URL.Path, "/")
		file, err := pkger.Open(filepath.Join("/static/", split[len(split)-2], split[len(split)-1]))
		if err != nil {
			panic(err)
		}
		io.Copy(w, file)
	}
}
