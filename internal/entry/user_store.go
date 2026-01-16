package entry

import (
	"context"
	"time"
)

// Store-facing user representation
type StoreUser struct {
	ID            int64
	CondominiumID int64
	FirstName     string
	LastName      string
	Email         string
	PasswordHash  string
	Role          UserRole
	Enabled       bool
	Hidden        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CreatedBy     int64
	UpdatedBy     int64
}

type UserStore interface {
	UserGetByID(ctx context.Context, id int64) (*StoreUser, error)
	UserGetByEmail(ctx context.Context, email string) (*StoreUser, error)
	UserCreate(ctx context.Context, u *StoreUser) (*StoreUser, error)
	UserCountByRole(ctx context.Context, role UserRole) (int64, error)
}
