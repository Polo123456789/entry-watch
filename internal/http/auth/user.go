package auth

import (
	"context"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

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

type UserWithPassword struct {
	*User
	PasswordHash string
}

type UserWithCondo struct {
	*User
	CondoName string
}

type UserStore interface {
	GetByEmailForAuth(ctx context.Context, email string) (UserWithPassword, bool, error)
	GetByID(ctx context.Context, id int64) (*User, bool, error)
	CreateUser(ctx context.Context, user *User, passwordHash string) (*User, error)
	CountSuperAdmins(ctx context.Context) (int64, error)
	UserListByRole(ctx context.Context, role entry.UserRole) ([]UserWithCondo, error)
	UserUpdate(ctx context.Context, id int64, user *User, updatedBy int64) (*User, error)
	UserUpdatePassword(ctx context.Context, id int64, passwordHash string, updatedBy int64) error
	UserDelete(ctx context.Context, id int64) error
}

func (u *User) ToEntryUser() *entry.User {
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
