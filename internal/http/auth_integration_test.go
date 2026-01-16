package http

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"

	"github.com/Polo123456789/entry-watch/db"
	"github.com/Polo123456789/entry-watch/internal/entry"
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

func TestLoginSetsSessionCookie(t *testing.T) {
	ctx := context.Background()
	logger := newTestLogger()
	dbConn := openTestDB(t)
	store := sqlstore.New(dbConn, logger)
	app := entry.NewApp(logger, store)

	// ensure superadmin exists
	if err := app.BootstrapSuperadmin(ctx); err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	// Create server and use its handler in an httptest server
	srv := NewServer("127.0.0.1", 0, app, logger)
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	// Prepare form data for login
	form := url.Values{}
	form.Set("email", "superadmin@local")
	form.Set("password", "password")

	client := &http.Client{
		// don't follow redirects so we can inspect Set-Cookie
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	resp, err := client.Post(ts.URL+"/auth/login", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusFound {
		t.Fatalf("expected redirect status, got %d", resp.StatusCode)
	}

	// Check Set-Cookie header for session cookie
	var sessionCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "entrywatch_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatalf("expected entrywatch_session cookie to be set")
	}

	// Decode session using same dev key as server
	cookieStore := sessions.NewCookieStore([]byte("dev-secret-key-please-change"))
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.AddCookie(sessionCookie)
	sess, err := cookieStore.Get(req, "entrywatch_session")
	if err != nil {
		t.Fatalf("failed to decode session cookie: %v", err)
	}
	v, ok := sess.Values["uid"]
	if !ok {
		t.Fatalf("session does not contain uid")
	}
	// ensure uid parses to a positive integer
	switch val := v.(type) {
	case int64:
		if val <= 0 {
			t.Fatalf("invalid uid value: %d", val)
		}
	case int:
		if int64(val) <= 0 {
			t.Fatalf("invalid uid value: %d", val)
		}
	case float64:
		if int64(val) <= 0 {
			t.Fatalf("invalid uid value: %v", val)
		}
	default:
		t.Fatalf("unexpected uid type: %T", v)
	}
}
