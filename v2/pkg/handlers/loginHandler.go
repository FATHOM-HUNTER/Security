package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/db"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type successResponse struct {
	Success bool `json:"success"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type authBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HandleLogin handles login and signUp route
func HandleLogin(router *mux.Router) {
	router.HandleFunc("/api/login", login)
	router.HandleFunc("/api/register", signUp)
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

func validateUser(username string, password []byte) (bool, error) {
	rows, err := db.DbConnection.Query(`SELECT username, password FROM auth WHERE username = @p1`, username)
	if err != nil {
		return false, err
	}

	var usernameP string
	var passwordP string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&usernameP, &passwordP)

		if err != nil {
			return false, err
		}
		break
	}

	return bcrypt.CompareHashAndPassword([]byte(passwordP), password) == nil, nil
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

	var validated bool
	var err error
	var success interface{}

	if userExists {
		validated, err = validateUser(username, password)
		if err != nil {
			success = errorResponse{err.Error()}
		}
	}

	success = successResponse{validated}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(success)
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
					ret = successResponse{true}
				}
			}
		}
	} else {
		ret = errorResponse{"username or password cant be empty"}
	}
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(ret)
}
