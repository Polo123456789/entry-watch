package auth

import (
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/entry"
	templates "github.com/Polo123456789/entry-watch/internal/templates/user"
	"github.com/gorilla/sessions"
	"log/slog"
)

// Handle mounts auth routes under /auth/
func Handle(app *entry.App, store sessions.Store, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := templates.Login().Render(r.Context(), w); err != nil {
				logger.Info("failed to render login template", "error", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
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

		// set session with uid
		sess, err := store.Get(r, "entrywatch_session")
		if err != nil {
			logger.Info("failed to get session", "error", err)
		} else {
			sess.Values["uid"] = user.ID
			sess.Options = &sessions.Options{
				Path:     "/",
				HttpOnly: true,
				MaxAge:   60 * 60 * 12, // 12 hours
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			}
			if err := sess.Save(r, w); err != nil {
				logger.Info("failed to save session", "error", err)
			}
		}

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
		sess, err := store.Get(r, "entrywatch_session")
		if err == nil && sess != nil {
			sess.Options = &sessions.Options{Path: "/", MaxAge: -1, HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode}
			_ = sess.Save(r, w)
		}
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	})

	return mux
}
