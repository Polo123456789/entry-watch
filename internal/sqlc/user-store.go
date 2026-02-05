package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
)

// UserStore wraps SQLC queries to provide user-related operations.
// This implements auth.UserStore interface.
type UserStore struct {
	queries *Queries
}

// NewUserStore creates a new UserStore that wraps the SQLC queries.
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		queries: New(db),
	}
}

// GetByEmailForAuth retrieves a user by email along with the password hash.
// Implements auth.UserStore.
// WARNING: Returns password hash - only use for authentication!
func (s *UserStore) GetByEmailForAuth(ctx context.Context, email string) (auth.UserWithPassword, bool, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.UserWithPassword{}, false, nil
		}
		return auth.UserWithPassword{}, false, err
	}

	return auth.UserWithPassword{
		User: &auth.User{
			ID:            user.ID,
			CondominiumID: getInt64FromNullInt64(user.CondominiumID),
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			Phone:         getStringFromNullString(user.Phone),
			Role:          entry.UserRole(user.Role),
			Enabled:       user.Enabled,
			Hidden:        user.Hidden,
		},
		PasswordHash: user.Password,
	}, true, nil
}

// GetByID retrieves a user by ID.
// Implements auth.UserStore.
func (s *UserStore) GetByID(ctx context.Context, id int64) (*auth.User, bool, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &auth.User{
		ID:            user.ID,
		CondominiumID: getInt64FromNullInt64(user.CondominiumID),
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Phone:         getStringFromNullString(user.Phone),
		Role:          entry.UserRole(user.Role),
		Enabled:       user.Enabled,
		Hidden:        user.Hidden,
	}, true, nil
}

// CreateUser creates a new user with the given password hash.
// Implements auth.UserStore.
// Returns the created user with the assigned ID from the database.
func (s *UserStore) CreateUser(ctx context.Context, email, firstName, lastName string, user *auth.User, passwordHash string) (*auth.User, error) {
	now := time.Now().Unix()

	var condoID sql.NullInt64
	if user.CondominiumID != 0 {
		condoID = sql.NullInt64{Int64: user.CondominiumID, Valid: true}
	}

	createdUser, err := s.queries.CreateUser(ctx, CreateUserParams{
		CondominiumID: condoID,
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Phone:         sql.NullString{},
		Role:          string(user.Role),
		Password:      passwordHash,
		Enabled:       user.Enabled,
		Hidden:        false,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedBy:     sql.NullInt64{},
		UpdatedBy:     sql.NullInt64{},
	})
	if err != nil {
		return nil, err
	}

	return &auth.User{
		ID:            createdUser.ID,
		CondominiumID: getInt64FromNullInt64(createdUser.CondominiumID),
		FirstName:     createdUser.FirstName,
		LastName:      createdUser.LastName,
		Email:         createdUser.Email,
		Phone:         getStringFromNullString(createdUser.Phone),
		Role:          entry.UserRole(createdUser.Role),
		Enabled:       createdUser.Enabled,
		Hidden:        createdUser.Hidden,
	}, nil
}

// CountSuperAdmins returns the number of enabled superadmins.
// Implements auth.UserStore.
func (s *UserStore) CountSuperAdmins(ctx context.Context) (int64, error) {
	return s.queries.CountSuperAdmins(ctx)
}

// getInt64FromNullInt64 safely extracts int64 from sql.NullInt64.
func getInt64FromNullInt64(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

// getStringFromNullString safely extracts string from sql.NullString.
func getStringFromNullString(n sql.NullString) string {
	if n.Valid {
		return n.String
	}
	return ""
}
