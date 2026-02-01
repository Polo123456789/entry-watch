package auth

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/util"
)

// AuthMiddleware validates the session and injects the user into the context.
// If the user is not authenticated, it redirects to the login page.
func AuthMiddleware(
	session sessions.Store,
	store UserStore,
	logger *slog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := CurrentUser(session, r)
			if !ok {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			ctx := entry.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RedirectIfAuthenticated redirects authenticated users away from the login page.
// This is used on the login page to redirect already logged-in users to their dashboard.
func RedirectIfAuthenticated(
	session sessions.Store,
	logger *slog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := CurrentUser(session, r)
			if ok {
				redirectURL := getRedirectForRole(user.Role)
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole creates a middleware that ensures the user has the required role.
// Superadmin bypasses all role checks.
func RequireRole(
	role entry.UserRole,
	session sessions.Store,
	logger *slog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := entry.RequireRole(r.Context(), role)
			if err != nil {
				util.HandleError(w, r, logger, err)
				return
			}

			// Re-inject user to ensure it's fresh
			ctx := entry.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRoleAndCondo creates a middleware that ensures the user has the required role
// and belongs to the specified condominium.
// Superadmin bypasses all checks.
func RequireRoleAndCondo(
	role entry.UserRole,
	condoID int64,
	session sessions.Store,
	logger *slog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := entry.RequireRoleAndCondo(r.Context(), role, condoID)
			if err != nil {
				util.HandleError(w, r, logger, err)
				return
			}

			ctx := entry.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
