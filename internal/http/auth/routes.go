package auth

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

// Handle sets up all authentication routes.
// Unauthenticated routes: /auth/login, /auth/logout
// The session store is passed in to be used by all auth handlers.
func Handle(
	logger *slog.Logger,
	session sessions.Store,
	userStore UserStore,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"GET /auth/login",
		hGetLogin(session, logger),
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
