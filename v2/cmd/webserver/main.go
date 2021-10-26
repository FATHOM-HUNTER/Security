package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Ovenoboyo/basic_webserver/v2/pkg/handlers"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/storage"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env")
	}

	// db.ConnectToDB()
	storage.InitializeStorage()

	r := mux.NewRouter()
	http.Handle("/", r)

	handlers.HandleStatic(r)
	handlers.HandleBlobs(r)
	// handlers.HandleLogin(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
