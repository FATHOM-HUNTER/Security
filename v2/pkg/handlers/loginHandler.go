package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/crypto"
	db "github.com/Ovenoboyo/basic_webserver/v2/pkg/database"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/middleware"

	"github.com/gorilla/mux"
)

// HandleLogin handles login and signUp route
func HandleLogin(router *mux.Router) {
	router.HandleFunc("/login", login)
	router.HandleFunc("/register", signUp)
}

func parseForm(req *http.Request) (string, []byte) {
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

func login(resp http.ResponseWriter, req *http.Request) {
	username, password := parseForm(req)
	userExists := db.UserExists(username)

	resp.Header().Set("Content-Type", "application/json")

	if userExists {
		validated, uid, err := db.ValidateUser(username, password)
		if err != nil {
			log.Println("here")
			resp.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(resp).Encode(errorResponse{
				Error: err.Error(),
			})
			return
		}
		if validated {
			token, err := middleware.GenerateToken()
			if err != nil {
				log.Println("here1")
				resp.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(resp).Encode(errorResponse{
					Error: err.Error(),
				})
				return
			}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(successResponse{
				Success: true,
				Data: authResponse{
					UID:   uid,
					Token: token,
				},
			})
			return
		} else {
			resp.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(resp).Encode(errorResponse{
				Error: "Invalid username/password",
			})
		}
	}

	json.NewEncoder(resp).Encode(errorResponse{
		Error: "User does not exist",
	})
}

func signUp(resp http.ResponseWriter, req *http.Request) {
	username, password := parseForm(req)
	var ret interface{}

	if len(username) > 0 && len(password) > 0 {
		if db.UserExists(username) {
			ret = errorResponse{"user already exists"}
		} else {
			saltedPass, err := crypto.HashAndSalt(string(password))
			if err != nil {
				ret = errorResponse{err.Error()}
			} else {
				err = db.WriteUser(username, saltedPass)
				if err != nil {
					ret = errorResponse{err.Error()}
				} else {
					ret = successResponse{
						Success: true,
					}
				}
			}
		}
	} else {
		ret = errorResponse{"username or password cant be empty"}
	}
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(ret)
}
