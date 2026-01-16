package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"log/slog"
)

// Handle mounts auth routes under /auth/
func Handle(app *entry.App, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body><form method=POST><input name=\"email\" placeholder=\"email\"/><input name=\"password\" type=\"password\" placeholder=\"password\"/><button type=submit>Login</button></form></body></html>"))
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if app == nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		if email == "" || password == "" {
			http.Error(w, "email and password required", http.StatusBadRequest)
			return
		}

		user, err := app.UserGetByEmail(r.Context(), email)
		if err != nil || user == nil {
			logger.Info("login failed: user not found", "email", email)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// compare bcrypt hash
		if err := entry.ComparePassword(user.PasswordHash, password); err != nil {
			logger.Info("login failed: bad password", "email", email)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// set cookie with uid
		c := &http.Cookie{
			Name:     "entrywatch_uid",
			Value:    fmt.Sprintf("%d", user.ID),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   60 * 60 * 12, // 12 hours
			Expires:  time.Now().Add(12 * time.Hour),
		}
		http.SetCookie(w, c)

		// redirect based on role
		switch user.Role {
		case entry.RoleSuperAdmin:
			http.Redirect(w, r, "/super/", http.StatusSeeOther)
		case entry.RoleAdmin:
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
		case entry.RoleGuardian:
			http.Redirect(w, r, "/guard/", http.StatusSeeOther)
		default:
			http.Redirect(w, r, "/neighbor/", http.StatusSeeOther)
		}
	})

	mux.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		// clear cookie
		c := &http.Cookie{
			Name:     "entrywatch_uid",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	})

	return mux
}
