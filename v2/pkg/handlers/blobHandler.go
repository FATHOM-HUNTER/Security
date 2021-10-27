package handlers

import (
	"net/http"

	db "github.com/Ovenoboyo/basic_webserver/v2/pkg/database"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/storage"
	"github.com/gorilla/mux"
)

// HandleBlobs registers all blob related routes
func HandleBlobs(router *mux.Router) {
	router.HandleFunc("/api/upload", uploadBlob).Methods("POST")
	router.HandleFunc("/api/list", listBlobs).Methods("GET")
}

func uploadBlob(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	uid := parseJWTToken(r)

	if len(uid) > 0 && len(filePath) > 0 {
		err := storage.UploadToStorage(&r.Body, filePath, uid)

		if err != nil {
			encodeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		encodeSuccess(w, nil)
		return
	}

	encodeError(w, http.StatusBadRequest, "Must provide uid and path as query params")
}

func listBlobs(w http.ResponseWriter, r *http.Request) {
	uid := parseJWTToken(r)
	if len(uid) > 0 {
		data, err := db.ListFilesForUser(uid)
		if err != nil {
			encodeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		encodeSuccess(w, data)
	}
}
