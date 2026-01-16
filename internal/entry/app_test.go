package entry

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestBootstrapSuperadmin_CreatesAndIsIdempotent(t *testing.T) {
	ctx := context.Background()
	logger := newTestLogger()
	store := NewMemStore()
	app := NewApp(logger, store)

	// First bootstrap should create one superadmin
	if err := app.BootstrapSuperadmin(ctx); err != nil {
		t.Fatalf("BootstrapSuperadmin failed: %v", err)
	}

	c, err := store.UserCountByRole(ctx, RoleSuperAdmin)
	if err != nil {
		t.Fatalf("UserCountByRole failed: %v", err)
	}
	if c != 1 {
		t.Fatalf("expected 1 superadmin after bootstrap, got %d", c)
	}

	// Second bootstrap should be idempotent
	if err := app.BootstrapSuperadmin(ctx); err != nil {
		t.Fatalf("second BootstrapSuperadmin failed: %v", err)
	}
	c2, err := store.UserCountByRole(ctx, RoleSuperAdmin)
	if err != nil {
		t.Fatalf("UserCountByRole failed: %v", err)
	}
	if c2 != 1 {
		t.Fatalf("expected 1 superadmin after second bootstrap, got %d", c2)
	}

	// Verify that the created user exists and has expected email
	u, err := store.UserGetByEmail(ctx, "superadmin@local")
	if err != nil {
		t.Fatalf("UserGetByEmail failed: %v", err)
	}
	if u == nil {
		t.Fatalf("expected user, got nil")
	}
	if u.Role != RoleSuperAdmin {
		t.Fatalf("expected role %s, got %s", RoleSuperAdmin, u.Role)
	}
	// Ensure password hash is non-empty and ComparePassword accepts "password"
	if u.PasswordHash == "" {
		t.Fatalf("expected password hash to be set")
	}
	if err := ComparePassword(u.PasswordHash, "password"); err != nil {
		t.Fatalf("ComparePassword failed for bootstrap password: %v", err)
	}
}
