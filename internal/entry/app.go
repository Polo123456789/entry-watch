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
