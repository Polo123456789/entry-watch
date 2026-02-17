package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
)

type UserStore struct {
	queries *Queries
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		queries: New(db),
	}
}

func (s *UserStore) GetByEmailForAuth(ctx context.Context, email string) (auth.UserWithPassword, bool, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.UserWithPassword{}, false, nil
		}
		return auth.UserWithPassword{}, false, err
	}

	return auth.UserWithPassword{
		User:         user.unmarshall(),
		PasswordHash: user.Password,
	}, true, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*auth.User, bool, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return user.unmarshall(), true, nil
}

func (s *UserStore) CreateUser(ctx context.Context, user *auth.User, passwordHash string) (*auth.User, error) {
	now := time.Now().Unix()

	var condoID sql.NullInt64
	if user.CondominiumID != 0 {
		condoID = sql.NullInt64{Int64: user.CondominiumID, Valid: true}
	}

	var phone sql.NullString
	if user.Phone != "" {
		phone = sql.NullString{String: user.Phone, Valid: true}
	}

	createdUser, err := s.queries.CreateUser(ctx, CreateUserParams{
		CondominiumID: condoID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Phone:         phone,
		Role:          string(user.Role),
		Password:      passwordHash,
		Enabled:       user.Enabled,
		Hidden:        user.Hidden,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedBy:     sql.NullInt64{},
		UpdatedBy:     sql.NullInt64{},
	})
	if err != nil {
		return nil, err
	}

	return createdUser.unmarshall(), nil
}

func (s *UserStore) CountSuperAdmins(ctx context.Context) (int64, error) {
	return s.queries.CountSuperAdmins(ctx)
}

type UserWithCondo struct {
	*auth.User
	CondoName string
}

func (s *UserStore) UserListByRole(ctx context.Context, role entry.UserRole) ([]UserWithCondo, error) {
	rows, err := s.queries.UserListByRole(ctx, string(role))
	if err != nil {
		return nil, err
	}
	result := make([]UserWithCondo, len(rows))
	for i, row := range rows {
		result[i] = UserWithCondo{
			User:      row.unmarshall(),
			CondoName: validNullString(row.CondoName),
		}
	}
	return result, nil
}

func (s *UserStore) UserUpdate(ctx context.Context, id int64, user *auth.User, updatedBy int64) (*auth.User, error) {
	var condoID sql.NullInt64
	if user.CondominiumID != 0 {
		condoID = sql.NullInt64{Int64: user.CondominiumID, Valid: true}
	}

	var phone sql.NullString
	if user.Phone != "" {
		phone = sql.NullString{String: user.Phone, Valid: true}
	}

	updated, err := s.queries.UserUpdate(ctx, UserUpdateParams{
		ID:            id,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Phone:         phone,
		CondominiumID: condoID,
		Enabled:       user.Enabled,
		UpdatedAt:     time.Now().Unix(),
		UpdatedBy:     sql.NullInt64{Int64: updatedBy, Valid: updatedBy != 0},
	})
	if err != nil {
		return nil, err
	}
	return updated.unmarshall(), nil
}

func (s *UserStore) UserDelete(ctx context.Context, id int64) error {
	return s.queries.UserDelete(ctx, id)
}

func (u UserListByRoleRow) unmarshall() *auth.User {
	return &auth.User{
		ID:            u.ID,
		CondominiumID: validNullInt64(u.CondominiumID),
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		Phone:         validNullString(u.Phone),
		Role:          entry.UserRole(u.Role),
		Enabled:       u.Enabled,
		Hidden:        u.Hidden,
	}
}
