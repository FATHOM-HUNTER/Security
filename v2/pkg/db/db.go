package db

import (
	"database/sql"
	"fmt"
	"os"

	// Driver for postgresql
	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("host")
	port     = os.Getenv("port")
	user     = os.Getenv("user")
	password = os.Getenv("password")
	dbname   = os.Getenv("dbname")
)

// DbConnection is the Global object for database connection
var DbConnection *sql.DB

// ConnectToDB Connects to postgresql db
func ConnectToDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	DbConnection, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	migrate()
}

func migrate() {
	_, err := DbConnection.Exec(`CREATE TABLE IF NOT EXISTS auth (
		username text PRIMARY KEY,
		password text NOT NULL
		);`)

	if err != nil {
		panic(err)
	}
}
