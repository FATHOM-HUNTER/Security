package handlers

import (
	"encoding/json"
	"log"
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
	encoder := json.NewEncoder(w)

	if len(uid) > 0 && len(filePath) > 0 {
		err := storage.UploadToStorage(&r.Body, filePath, uid)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(errorResponse{
				Error: err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		encoder.Encode(successResponse{
			Success: true,
		})
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	encoder.Encode(errorResponse{
		Error: "Must provide uid and path as query params",
	})
}
