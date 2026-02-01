package auth

import (
	"context"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

// UserStore defines the interface for user storage operations.
// Implementations are responsible for converting between the SQLC model
// and the domain model (entry.User).
type UserStore interface {
	// GetByEmail retrieves a user by their email address.
	// Returns (nil, false, nil) if the user is not found.
	// Returns (nil, false, error) if there is a database error.
	GetByEmail(ctx context.Context, email string) (*entry.User, bool, error)

	// GetByID retrieves a user by their ID.
	// Returns (nil, false, nil) if the user is not found.
	// Returns (nil, false, error) if there is a database error.
	GetByID(ctx context.Context, id int64) (*entry.User, bool, error)

	// CreateUser creates a new user with the given password hash.
	// The password must already be hashed before calling this method.
	CreateUser(ctx context.Context, email, firstName, lastName string, user *entry.User, passwordHash string) error

	// CountSuperAdmins returns the number of enabled superadmins.
	CountSuperAdmins(ctx context.Context) (int64, error)
}
