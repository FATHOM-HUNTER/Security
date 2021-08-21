package main

import (
	"net/http"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/handlers"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/db"

	"github.com/gorilla/mux"
)

func main() {
	db.ConnectToDB()

	r := mux.NewRouter()
	http.Handle("/", r)
	handlers.HandleLogin(r)
}
