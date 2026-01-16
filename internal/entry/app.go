package entry

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Config Config
	store  Store
	logger *slog.Logger
}

func NewApp(logger *slog.Logger, store Store) *App {
	return &App{
		store:  store,
		logger: logger,
		Config: Config{},
	}
}

type Store interface {
	CondominiumStore
	VisitStore
	UserStore
}

type Config struct{}

type Valid interface {
	Valid() error
}

// BootstrapSuperadmin ensures at least one superadmin user exists. If none
// exists, it creates a default one with email "superadmin@local" and
// password "password". This is intended for development/bootstrap only.
func (a *App) BootstrapSuperadmin(ctx context.Context) error {
	if a.store == nil {
		return nil
	}
	c, err := a.store.UserCountByRole(ctx, RoleSuperAdmin)
	if err != nil {
		return err
	}
	if c > 0 {
		return nil
	}
	// create default
	pw := "password"
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := &StoreUser{
		CondominiumID: 0,
		FirstName:     "Super",
		LastName:      "Admin",
		Email:         "superadmin@local",
		PasswordHash:  string(hash),
		Role:          RoleSuperAdmin,
		Enabled:       true,
		Hidden:        false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, err = a.store.UserCreate(ctx, u)
	return err
}

// UserGetByID is a small wrapper that exposes the underlying store's
// UserGetByID implementation to packages outside the entry package. This
// allows HTTP middleware to lookup users without accessing internal store
// fields directly.
func (a *App) UserGetByID(ctx context.Context, id int64) (*StoreUser, error) {
	if a.store == nil {
		return nil, ErrNotFound
	}
	return a.store.UserGetByID(ctx, id)
}

// UserGetByEmail exposes the underlying store's UserGetByEmail implementation
// so HTTP handlers and middleware can lookup users by email.
func (a *App) UserGetByEmail(ctx context.Context, email string) (*StoreUser, error) {
	if a.store == nil {
		return nil, ErrNotFound
	}
	return a.store.UserGetByEmail(ctx, email)
}
