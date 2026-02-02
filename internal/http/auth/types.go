package auth

import (
	"context"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

// User represents an authenticated user in the auth system.
// This mirrors the database schema from db/migrations/20251123181233_init.sql
// This is separate from entry.User to maintain domain separation.
type User struct {
	ID            int64
	CondominiumID int64
	FirstName     string
	LastName      string
	Email         string
	Phone         string
	Role          entry.UserRole
	Enabled       bool
	Hidden        bool
}

// UserWithPassword extends User with the password hash for authentication.
type UserWithPassword struct {
	*User
	PasswordHash string
}

// UserSafeError represents an error that can be safely shown to users.
// These errors are typically caused by user input and don't expose system details.
type UserSafeError struct {
	msg string
}

func (e *UserSafeError) Error() string {
	return e.msg
}

// NewUserSafeError creates a new user-safe error with the given message.
func NewUserSafeError(msg string) *UserSafeError {
	return &UserSafeError{msg: msg}
}

// UserStore defines the interface for user storage operations.
// Implementations are responsible for converting between the SQLC model
// and the auth model.
type UserStore interface {
	// GetByEmail retrieves a user by their email address.
	// Returns the user with password hash for authentication.
	// Returns (UserWithPassword{}, false, nil) if the user is not found.
	// Returns (UserWithPassword{}, false, error) if there is a database error.
	GetByEmail(ctx context.Context, email string) (UserWithPassword, bool, error)

	// GetByID retrieves a user by their ID.
	// Returns (nil, false, nil) if the user is not found.
	// Returns (nil, false, error) if there is a database error.
	GetByID(ctx context.Context, id int64) (*User, bool, error)

	// CreateUser creates a new user with the given password hash.
	// The password must already be hashed before calling this method.
	CreateUser(ctx context.Context, email, firstName, lastName string, user *User, passwordHash string) error

	// CountSuperAdmins returns the number of enabled superadmins.
	CountSuperAdmins(ctx context.Context) (int64, error)
}

// toEntryUser converts an auth.User to an entry.User.
// This is used when passing user data to the domain layer.
// Only copies the fields needed for domain-level authorization.
func toEntryUser(u *User) *entry.User {
	if u == nil {
		return nil
	}
	return &entry.User{
		ID:            u.ID,
		CondominiumID: u.CondominiumID,
		Role:          u.Role,
		Enabled:       u.Enabled,
	}
}
