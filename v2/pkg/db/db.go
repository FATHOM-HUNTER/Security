package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	// Driver for Azure MySQL
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"

	// Driver for golang-migrate to read files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// DbConnection is the Global object for database connection
var DbConnection *sql.DB

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

	DbConnection, err = sql.Open("sqlserver", psqlInfo)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = DbConnection.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected!")

	migrateDB()
}

func migrateDB() {
	fmt.Println("Migrating")
	driver, err := sqlserver.WithInstance(DbConnection, &sqlserver.Config{})
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

// AddFileMetaToDB adds uploaded file to database
func AddFileMetaToDB(fileName string, md5 string, uid string, contents int) error {
	lastModified := strconv.FormatInt(time.Now().UnixMilli(), 10)
	id := uuid.New().String()
	_, err := DbConnection.Exec(`INSERT INTO file_meta (id, file_name, uid, last_modified, md5_hash, file_contents, version) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`, id, fileName, uid, lastModified, md5, contents, 1)

	return err
}
