package main

import (
	"log"
	"net/http"

	db "github.com/Ovenoboyo/basic_webserver/v2/pkg/database"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/handlers"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/middleware"
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/storage"
	"github.com/joho/godotenv"
	"github.com/markbates/pkger"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env")
	}

	db.ConnectToDB()
	storage.InitializeStorage()

	pkger.Include("/static")

	r := mux.NewRouter()
	apiRouter := mux.NewRouter()
	apiRouterNegroni := middleware.GetJWTWrappedNegroni(apiRouter)

	r.PathPrefix("/api").Handler(apiRouterNegroni)

	http.Handle("/", r)

	handlers.HandleStatic(r)
	handlers.HandleBlobs(apiRouter)
	handlers.HandleLogin(r)

	n := negroni.Classic()
	n.UseHandler(r)

	n.Run(":8080")
}
