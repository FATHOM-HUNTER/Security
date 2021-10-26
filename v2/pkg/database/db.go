package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	// Driver for Azure MySQL
	"github.com/Ovenoboyo/basic_webserver/v2/pkg/crypto"

	// Driver for microsoft sqlserver
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"

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

// ValidateUser validates user credentials from database
func ValidateUser(username string, password []byte) (bool, string, error) {
	rows, err := dbConnection.Query(`SELECT username, password, uid FROM auth WHERE username = @p1`, username)
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

	success := crypto.ValidatePassword(password, passwordP)
	if success {
		return true, uidP, nil
	}

	return false, "", errors.New("Invalid username or password")
}

// UserExists checks if user exists in database
func UserExists(username string) bool {
	rows, err := dbConnection.Query(`SELECT username FROM auth WHERE username = @p1`, username)
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

func WriteUser(username string, password []byte) error {
	uid := uuid.New()
	_, err := dbConnection.Exec(`INSERT INTO auth (username, uid, password) VALUES (@p1, @p2, @p3)`, username, uid, string(password))
	return err
}

// AddFileMetaToDB adds uploaded file to database
func AddFileMetaToDB(fileName string, md5 string, uid string, contents int) error {
	lastModified := strconv.FormatInt(time.Now().UnixMilli(), 10)
	id := uuid.New().String()
	_, err := dbConnection.Exec(`INSERT INTO file_meta (id, file_name, uid, last_modified, md5_hash, file_contents, version) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`, id, fileName, uid, lastModified, md5, contents, 1)

	return err
}
