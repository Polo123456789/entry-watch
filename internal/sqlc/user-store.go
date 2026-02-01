package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

// UserStore wraps SQLC queries to provide user-related operations.
// This is used by the auth package to implement auth.UserStore.
type UserStore struct {
	queries *Queries
}

// NewUserStore creates a new UserStore that wraps the SQLC queries.
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		queries: New(db),
	}
}

// GetByEmail retrieves a user by email along with the password hash.
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*entry.User, string, bool, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", false, nil
		}
		return nil, "", false, err
	}

	return convertUserToDomain(user), user.Password, true, nil
}

// GetByID retrieves a user by ID.
func (s *UserStore) GetByID(ctx context.Context, id int64) (*entry.User, bool, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return convertUserToDomain(user), true, nil
}

// CreateUser creates a new user with the given password hash.
func (s *UserStore) CreateUser(ctx context.Context, email, firstName, lastName string, user *entry.User, passwordHash string) error {
	now := time.Now().Unix()

	var condoID sql.NullInt64
	if user.CondominiumID != 0 {
		condoID = sql.NullInt64{Int64: user.CondominiumID, Valid: true}
	}

	return s.queries.CreateUser(ctx, CreateUserParams{
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
}

// CountSuperAdmins returns the number of enabled superadmins.
func (s *UserStore) CountSuperAdmins(ctx context.Context) (int64, error) {
	return s.queries.CountSuperAdmins(ctx)
}

// convertUserToDomain converts a SQLC User to an entry.User.
func convertUserToDomain(u User) *entry.User {
	user := &entry.User{
		ID:      u.ID,
		Role:    entry.UserRole(u.Role),
		Enabled: u.Enabled,
	}

	if u.CondominiumID.Valid {
		user.CondominiumID = u.CondominiumID.Int64
	}

	return user
}
