package crypto

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash, err
}

func ValidatePassword(provided []byte, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), (provided)) == nil
}
