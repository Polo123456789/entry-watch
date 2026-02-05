package auth

import (
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

// EnsureSuperAdminExists checks if there is at least one enabled superadmin.
// If not, it creates a default superadmin with credentials:
// - Email: admin@localhost
// - Password: changeme
// This is called at application startup.
func EnsureSuperAdminExists(ctx context.Context, store UserStore, logger *slog.Logger) error {
	count, err := store.CountSuperAdmins(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		logger.Info("Superadmin exists, skipping bootstrap")
		return nil
	}

	logger.Warn("No superadmin found - creating default account",
		"email", "admin@localhost",
		"warning", "CHANGE PASSWORD IMMEDIATELY AFTER FIRST LOGIN",
	)

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("changeme"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &User{
		FirstName:     "Super",
		LastName:      "Admin",
		Email:         "admin@localhost",
		Role:          entry.RoleSuperAdmin,
		Enabled:       true,
		CondominiumID: 0,
		Hidden:        false,
	}

	createdUser, err := store.CreateUser(ctx, user, string(passwordHash))
	if err != nil {
		return err
	}

	logger.Warn("Default superadmin created successfully",
		"email", "admin@localhost",
		"id", createdUser.ID,
	)

	return nil
}
