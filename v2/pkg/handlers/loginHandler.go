package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/db"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/middleware"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

// HandleLogin handles login and signUp route
func HandleLogin(router *mux.Router) {
	router.HandleFunc("/login", login)
	router.HandleFunc("/register", signUp)
}

func hashAndSalt(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash, err
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

func validateUser(username string, password []byte) (bool, string, error) {
	rows, err := db.DbConnection.Query(`SELECT username, password, uid FROM auth WHERE username = @p1`, username)
	if err != nil {
		return false, "", err
	}

	var usernameP string
	var passwordP string
	var uidP string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&usernameP, &passwordP, &uidP)

		if err != nil {
			return false, "", err
		}
		break
	}

	success := bcrypt.CompareHashAndPassword([]byte(passwordP), password) == nil
	if success {
		return true, uidP, nil
	}

	return false, "", errors.New("Invalid username or password")
}

func userExists(username string) bool {
	rows, err := db.DbConnection.Query(`SELECT username FROM auth WHERE username = @p1`, username)
	if err != nil {
		panic(err)
		return false
	}

	var usernameP string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&usernameP)

		if err != nil {
			panic(err)
			return true
		}
		break
	}

	return username == usernameP
}

func writeUser(username string, password []byte) error {
	uid := uuid.New()
	_, err := db.DbConnection.Exec(`INSERT INTO auth (username, uid, password) VALUES (@p1, @p2, @p3)`, username, uid, string(password))
	return err
}

func login(resp http.ResponseWriter, req *http.Request) {
	username, password := parseForm(req)
	userExists := userExists(username)

	resp.Header().Set("Content-Type", "application/json")

	if userExists {
		validated, uid, err := validateUser(username, password)
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
		if userExists(username) {
			ret = errorResponse{"user already exists"}
		} else {
			saltedPass, err := hashAndSalt(string(password))
			if err != nil {
				ret = errorResponse{err.Error()}
			} else {
				err = writeUser(username, saltedPass)
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
