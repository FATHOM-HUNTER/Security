package database // ValidateUser validates user credentials from database
import (
	"errors"
	"log"

	"github.com/Ovenoboyo/basic_webserver/pkg/crypto"
	"github.com/google/uuid"
)

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
		log.Println(err)
		return false
	}

	var usernameP string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&usernameP)

		if err != nil {
			log.Println(err)
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
