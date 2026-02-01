package auth

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/sqlc"
)

// Handle sets up all authentication routes.
// Unauthenticated routes: /auth/login, /auth/logout
// The session store is passed in to be used by all auth handlers.
func Handle(
	app *entry.App,
	logger *slog.Logger,
	session sessions.Store,
	db *sql.DB,
) http.Handler {
	mux := http.NewServeMux()

	// Create user store using the database connection
	sqlcUserStore := sqlc.NewUserStore(db)
	userStore := NewSQLCUserStore(sqlcUserStore)

	// Unauthenticated routes
	mux.Handle(
		"GET /auth/login",
		RedirectIfAuthenticated(session, logger)(
			hGetLogin(session, logger),
		),
	)
	mux.Handle(
		"POST /auth/login",
		hPostLogin(session, userStore, logger),
	)
	mux.Handle(
		"GET /auth/logout",
		hGetLogout(session, logger),
	)

	return mux
}
