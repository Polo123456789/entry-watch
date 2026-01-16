package entry_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"

	"github.com/Polo123456789/entry-watch/db"
	entrypkg "github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/entry/sqlstore"
	_ "modernc.org/sqlite"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// openTestDB opens an in-memory sqlite DB, runs migrations and returns the *sql.DB.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	logger := newTestLogger()
	dbConn, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", err)
	}
	// ensure deterministic single connection
	dbConn.SetMaxOpenConns(1)
	if err := db.AutoMigrate(dbConn, logger); err != nil {
		dbConn.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}
	t.Cleanup(func() { dbConn.Close() })
	return dbConn
}

func TestBootstrapSuperadmin_CreatesAndIsIdempotent(t *testing.T) {
	ctx := context.Background()
	logger := newTestLogger()
	// use sqlite-backed store for tests
	dbConn := openTestDB(t)
	store := sqlstore.New(dbConn, logger)
	app := entrypkg.NewApp(logger, store)

	// First bootstrap should create one superadmin
	if err := app.BootstrapSuperadmin(ctx); err != nil {
		t.Fatalf("BootstrapSuperadmin failed: %v", err)
	}

	c, err := store.UserCountByRole(ctx, entrypkg.RoleSuperAdmin)
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
	c2, err := store.UserCountByRole(ctx, entrypkg.RoleSuperAdmin)
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
	if u.Role != entrypkg.RoleSuperAdmin {
		t.Fatalf("expected role %s, got %s", entrypkg.RoleSuperAdmin, u.Role)
	}
	// Ensure password hash is non-empty and ComparePassword accepts "password"
	if u.PasswordHash == "" {
		t.Fatalf("expected password hash to be set")
	}
	if err := entrypkg.ComparePassword(u.PasswordHash, "password"); err != nil {
		t.Fatalf("ComparePassword failed for bootstrap password: %v", err)
	}
}
