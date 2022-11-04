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

func checkUsername(username string) bool {
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

func checkEmail(email string) bool {
	rows, err := dbConnection.Query(`SELECT email FROM auth WHERE email = @p1`, email)
	if err != nil {
		log.Println(err)
		return false
	}

	var emailP string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&emailP)

		if err != nil {
			log.Println(err)
			return true
		}
		break
	}

	return email == emailP
}

// UserExists checks if user exists in database
func UsernameExists(username string) bool {
	return checkUsername(username)
}

func UsernameAndEmailExists(username string, email string) bool {
	return checkEmail(email) && checkUsername(username)
}

func WriteUser(username string, email string, password []byte) error {
	uid := uuid.New()
	_, err := dbConnection.Exec(`INSERT INTO auth (username, uid, email, password) VALUES (@p1, @p2, @p3, @p4)`, username, uid, email, string(password))
	return err
}
