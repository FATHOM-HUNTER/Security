package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/db"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.ConnectToDB()

	r := mux.NewRouter()
	http.Handle("/", r)

	handlers.HandleStatic(r)
	handlers.HandleLogin(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
