package superadmin

import (
	"context"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

type AdminStore interface {
	UserListByRole(ctx context.Context, role entry.UserRole) ([]*entry.AdminUser, error)
	UserGetByID(ctx context.Context, id int64) (*entry.AdminUser, bool, error)
	UserCreate(ctx context.Context, user *entry.AdminUser, password string) (*entry.AdminUser, error)
	UserUpdate(ctx context.Context, id int64, user *entry.AdminUser, updatedBy int64) (*entry.AdminUser, error)
	UserDelete(ctx context.Context, id int64) error
	CondoList(ctx context.Context) ([]*entry.Condominium, error)
}
