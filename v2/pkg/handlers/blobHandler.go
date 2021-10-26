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

	encoder := json.NewEncoder(w)
	err := storage.UploadToStorage(&r.Body, filePath)

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
}
