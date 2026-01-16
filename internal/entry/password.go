package entry

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword compares a stored bcrypt hash with the provided password.
// Returns nil on success, or an error if they don't match.
func ComparePassword(hash string, password string) error {
	if hash == "" {
		return errors.New("empty password hash")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
