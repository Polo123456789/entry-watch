package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/util"
	templates "github.com/Polo123456789/entry-watch/internal/templates/auth"
)

const authSessionKey = "entry-watch-auth"

// UserSafeError creates user-safe error messages in Spanish
type UserSafeError struct {
	msg string
}

func (e *UserSafeError) Error() string {
	return e.msg
}

func NewUserSafeError(msg string) *UserSafeError {
	return &UserSafeError{msg: msg}
}

func hGetLogin(
	session sessions.Store,
	logger *slog.Logger,
) http.Handler {
	return util.Handler(logger, func(
		w http.ResponseWriter, r *http.Request,
	) error {
		user, ok := CurrentUser(session, r)
		if ok {
			redirectURL := getRedirectForRole(user.Role)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return nil
		}

		errorParam := r.URL.Query().Get("error")
		hasError := errorParam == "1"

		return templates.Login(hasError).Render(r.Context(), w)
	})
}

func hPostLogin(
	session sessions.Store,
	store UserStore,
	logger *slog.Logger,
) http.Handler {
	return util.Handler(logger, func(
		w http.ResponseWriter, r *http.Request,
	) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := attemptLogin(r.Context(), store, email, password)
		if err != nil {
			return err
		}

		if err := setCurrentUser(w, r, session, user); err != nil {
			return err
		}

		redirectURL := getRedirectForRole(user.Role)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return nil
	})
}

func hGetLogout(
	session sessions.Store,
	logger *slog.Logger,
) http.Handler {
	return util.Handler(logger, func(
		w http.ResponseWriter, r *http.Request,
	) error {
		s, _ := session.Get(r, authSessionKey)
		s.Options.MaxAge = -1
		_ = s.Save(r, w)
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return nil
	})
}

func attemptLogin(
	ctx context.Context,
	store UserStore,
	email string,
	password string,
) (*entry.User, error) {
	wrongCredsErr := NewUserSafeError("Correo electrónico o contraseña incorrectos")

	userWithPass, ok, err := store.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, wrongCredsErr
	}

	if !userWithPass.Enabled {
		return nil, NewUserSafeError("La cuenta está deshabilitada")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(userWithPass.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, wrongCredsErr
	}

	return userWithPass.User, nil
}

func setCurrentUser(
	w http.ResponseWriter,
	r *http.Request,
	session sessions.Store,
	user *entry.User,
) error {
	s, err := session.Get(r, authSessionKey)
	if err != nil {
		return err
	}

	s.Values["user_id"] = user.ID
	s.Values["role"] = string(user.Role)
	s.Values["condominium_id"] = user.CondominiumID
	s.Values["enabled"] = user.Enabled

	s.Options.MaxAge = 60 * 60 * 12
	s.Options.HttpOnly = true
	s.Options.Secure = false
	s.Options.SameSite = http.SameSiteLaxMode

	return s.Save(r, w)
}

func CurrentUser(
	session sessions.Store,
	r *http.Request,
) (*entry.User, bool) {
	s, _ := session.Get(r, authSessionKey)

	userID, ok := s.Values["user_id"].(int64)
	if !ok {
		return nil, false
	}

	roleStr, ok := s.Values["role"].(string)
	if !ok {
		return nil, false
	}

	condoID, _ := s.Values["condominium_id"].(int64)
	enabled, _ := s.Values["enabled"].(bool)

	return &entry.User{
		ID:            userID,
		Role:          entry.UserRole(roleStr),
		CondominiumID: condoID,
		Enabled:       enabled,
	}, true
}

func getRedirectForRole(role entry.UserRole) string {
	switch role {
	case entry.RoleSuperAdmin:
		return "/super/"
	case entry.RoleAdmin:
		return "/admin/"
	case entry.RoleGuardian:
		return "/guard/"
	case entry.RoleUser:
		return "/neighbor/"
	default:
		return "/auth/login"
	}
}
