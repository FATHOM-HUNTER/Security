package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/db"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type successResponse struct {
	success bool
}

type errorResponse struct {
	data string
}

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
		fmt.Println(err)
		return "", nil
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	saltedPass, err := hashAndSalt(password)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	return username, saltedPass
}

func userExists(username string, password []byte) bool {
	rows, err := db.DbConnection.Query(`SELECT EXISTS(SELECT username FROM auth WHERE password = $1 AND username = $2)`, password, username)
	if err != nil {
		fmt.Println(err)
		return false
	}

	exists := false

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&exists)

		if err != nil {
			fmt.Println(err)
			return exists
		}
		break
	}

	return exists
}

func writeUser(username string, password []byte) error {
	_, err := db.DbConnection.Exec(`INSERT INTO auth (username, password) VALUES ($1, $2)`, username, password)
	return err
}

func login(resp http.ResponseWriter, req *http.Request) {
	username, saltedPass := parseForm(req)
	userExists := userExists(username, saltedPass)

	success := successResponse{userExists}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(success)
}

func signUp(resp http.ResponseWriter, req *http.Request) {
	username, saltedPass := parseForm(req)

	var ret interface{}
	if userExists(username, saltedPass) {
		ret = errorResponse{"user already exists"}
	} else {
		err := writeUser(username, saltedPass)
		if err != nil {
			ret = errorResponse{err.Error()}
		} else {
			ret = successResponse{true}
		}
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(ret)
}
