package handlers

import (
	"encoding/json"
	"net/http"

	db "github.com/Ovenoboyo/basic_webserver/pkg/database"
	"github.com/gorilla/mux"
)

// HandleShares registers all share related routes
func HandleShares(router *mux.Router) {
	router.HandleFunc("/api/shareFile", addSharedEmail).Methods("POST")
}

type shareFileParams struct {
	fileName        string
	sharedWithEmail string
}

func addSharedEmail(w http.ResponseWriter, r *http.Request) {
	ownerUID := parseJWTToken(r)

	var s shareFileParams
	json.NewDecoder(r.Body).Decode(&s)

	if len(ownerUID) > 0 && len(s.fileName) > 0 && len(ownerUID) > 0 && len(s.sharedWithEmail) > 0 {
		err := db.AddSharedEmail(s.fileName, ownerUID, s.sharedWithEmail)
		if err != nil {
			encodeError(w, http.StatusBadRequest, "Failed to write to database")
		}
		encodeSuccessHeader(w)
		return
	}

	encodeError(w, http.StatusBadRequest, "uid, filename and email to share to must be provided")
}

func getSharedEmails(w http.ResponseWriter, r *http.Request) {
	ownerUID := parseJWTToken(r)

	var s shareFileParams
	json.NewDecoder(r.Body).Decode(&s)

	if len(ownerUID) > 0 && len(s.fileName) > 0 {
		files, err := db.GetSharedEmails(s.fileName, ownerUID)
		if err != nil {
			encodeError(w, http.StatusBadRequest, "Failed to get shared emails from database")
		}
		encodeSuccess(w, files)
		return
	}

	encodeError(w, http.StatusBadRequest, "uid and filename to must be provided")
}

func removeSharedEmail(w http.ResponseWriter, r *http.Request) {
	ownerUID := parseJWTToken(r)

	var s shareFileParams
	json.NewDecoder(r.Body).Decode(&s)

	if len(ownerUID) > 0 && len(s.fileName) > 0 && len(s.sharedWithEmail) > 0 {
		err := db.RemoveSharedEmail(s.fileName, ownerUID, s.sharedWithEmail)
		if err != nil {
			encodeError(w, http.StatusBadRequest, "Failed to remove shared email from database")
		}
		encodeSuccessHeader(w)
		return
	}

	encodeError(w, http.StatusBadRequest, "UID and filename to must be provided")
}
