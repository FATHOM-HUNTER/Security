package handlers

import (
	"net/http"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/storage"
	"github.com/gorilla/mux"
)

// HandleBlobs registers all blob related routes
func HandleBlobs(router *mux.Router) {
	router.HandleFunc("/api/upload", uploadBlob).Methods("POST")
}

func uploadBlob(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	uid := r.URL.Query().Get("uid")

	if len(uid) > 0 && len(filePath) > 0 {
		err := storage.UploadToStorage(&r.Body, filePath, uid)

		if err != nil {
			encodeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		encodeSuccess(w)
		return
	}

	encodeError(w, http.StatusBadRequest, "Must provide uid and path as query params")
}
