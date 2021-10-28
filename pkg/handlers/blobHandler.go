package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	db "github.com/Ovenoboyo/basic_webserver/pkg/database"
	"github.com/Ovenoboyo/basic_webserver/pkg/storage"
	"github.com/gorilla/mux"
)

// HandleBlobs registers all blob related routes
func HandleBlobs(router *mux.Router) {
	router.HandleFunc("/api/upload", uploadBlob).Methods("POST")
	router.HandleFunc("/api/list", listBlobs).Methods("GET")
	router.HandleFunc("/api/download", downloadBlobs).Methods("GET")
	router.HandleFunc("/api/delete", deleteBlobs).Methods("POST")

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

func parseDeleteForm(req *http.Request) (string, string) {
	err := req.ParseForm()
	if err != nil {
		return "", ""
	}

	var a deleteBody
	err = json.NewDecoder(req.Body).Decode(&a)
	if err != nil {
		return "", ""
	}

	return a.FileName, a.Version
}

func deleteBlobs(w http.ResponseWriter, r *http.Request) {
	uid := parseJWTToken(r)
	fileName, version := parseDeleteForm(r)

	if len(fileName) > 0 && len(version) > 0 {

		err := storage.DeleteBlob(fileName, uid, version)
		if err != nil {
			encodeError(w, http.StatusBadRequest, err.Error())
			return
		}

		encodeSuccess(w, nil)
		return
	}
	encodeError(w, http.StatusBadRequest, "Filename and version must be provided")
}

func downloadBlobs(w http.ResponseWriter, r *http.Request) {
	uid := parseJWTToken(r)
	fileName := r.URL.Query().Get("path")
	version := r.URL.Query().Get("version")

	if len(fileName) > 0 && len(version) > 0 {
		stream, err := storage.DownloadBlob(fileName, uid, version)
		if err != nil {
			encodeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		encodeSuccessHeader(w)
		io.Copy(w, stream)
		return
	}
	encodeError(w, http.StatusBadRequest, "Filename and version must be provided")
}
