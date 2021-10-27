package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	// Driver for microsoft sqlserver
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"

	// Driver for golang-migrate to read files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// dbConnection is the Global object for database connection
var dbConnection *sql.DB

// ConnectToDB Connects to postgresql db
func ConnectToDB() {
	var (
		server   = os.Getenv("DB_SERVER")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DATABASE")
	)

	psqlInfo := fmt.Sprintf("server=%s;user id=%s;password=%s;port=1433;database=%s;",
		server, user, password, database)

	var err error

	dbConnection, err = sql.Open("sqlserver", psqlInfo)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = dbConnection.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected!")

	migrateDB()
}

func migrateDB() {
	fmt.Println("Migrating")
	driver, err := sqlserver.WithInstance(dbConnection, &sqlserver.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/", "sqlserver", driver)
	if err != nil {
		panic(err)
	}

	err = m.Migrate(2)
	if err != nil {
		log.Println(err)
	}
}
