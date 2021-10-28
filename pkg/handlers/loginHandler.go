package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Ovenoboyo/basic_webserver/pkg/crypto"
	db "github.com/Ovenoboyo/basic_webserver/pkg/database"
	"github.com/Ovenoboyo/basic_webserver/pkg/middleware"

	"github.com/gorilla/mux"
)

// HandleLogin handles login and signUp route
func HandleLogin(router *mux.Router) {
	router.HandleFunc("/login", login)
	router.HandleFunc("/register", signUp)
	router.HandleFunc("/validate", validateToken)
	router.HandleFunc("/health", healthCheck)
}

func parseAuthForm(req *http.Request) (string, []byte) {
	err := req.ParseForm()
	if err != nil {
		return "", nil
	}

	var a authBody
	err = json.NewDecoder(req.Body).Decode(&a)
	if err != nil {
		return "", nil
	}

	return a.Username, []byte(a.Password)
}

func login(w http.ResponseWriter, r *http.Request) {
	username, password := parseAuthForm(r)
	userExists := db.UserExists(username)

	w.Header().Set("Content-Type", "application/json")

	if userExists {

		validated, uid, err := db.ValidateUser(username, password)
		if err != nil {
			encodeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if validated {
			token, err := middleware.GenerateToken(uid)
			if err != nil {
				encodeError(w, http.StatusInternalServerError, err.Error())
				return
			}

			encodeSuccess(w, authResponse{
				Token: token,
			})
			return
		}

		encodeError(w, http.StatusUnauthorized, "Invalid username/password")
		return
	}

	encodeError(w, http.StatusUnauthorized, "User does not exist")
}

func signUp(w http.ResponseWriter, r *http.Request) {
	username, password := parseAuthForm(r)

	if len(username) > 0 && len(password) > 0 {
		if db.UserExists(username) {
			encodeError(w, http.StatusBadRequest, "User already exists")
			return
		}

		saltedPass, err := crypto.HashAndSalt(string(password))
		if err != nil {
			encodeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = db.WriteUser(username, saltedPass)
		if err != nil {
			encodeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		encodeSuccess(w, nil)
		return

	}
	encodeError(w, http.StatusInternalServerError, "Username or password cannot be empty")
}

func validateToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	encodeSuccessHeader(w)
	json.NewEncoder(w).Encode(successResponse{
		Success: middleware.ValidateToken(token),
	})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	encodeSuccessHeader(w)
}
