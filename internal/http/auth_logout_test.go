package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/entry/sqlstore"
	_ "modernc.org/sqlite"
)

func TestLogoutClearsSessionCookie(t *testing.T) {
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

	// find session cookie from login
	var sessionCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "entrywatch_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatalf("expected entrywatch_session cookie to be set after login")
	}

	// now call logout with the session cookie
	req, _ := http.NewRequest("GET", ts.URL+"/auth/logout", nil)
	req.AddCookie(sessionCookie)
	resp2, err := client.Do(req)
	if err != nil {
		t.Fatalf("logout request failed: %v", err)
	}
	defer resp2.Body.Close()

	// Look for Set-Cookie header for entrywatch_session indicating deletion
	var logoutCookie *http.Cookie
	for _, c := range resp2.Cookies() {
		if c.Name == "entrywatch_session" {
			logoutCookie = c
			break
		}
	}
	if logoutCookie == nil {
		t.Fatalf("expected entrywatch_session cookie to be set on logout response")
	}

	// When MaxAge <= 0 the cookie is expired/deleted
	if logoutCookie.MaxAge > 0 {
		t.Fatalf("expected logout cookie MaxAge <= 0, got %d", logoutCookie.MaxAge)
	}
}
