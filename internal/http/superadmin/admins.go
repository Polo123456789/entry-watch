package superadmin

import (
	"log/slog"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
	"github.com/Polo123456789/entry-watch/internal/http/util"
	templates "github.com/Polo123456789/entry-watch/internal/templates/superadmin"
	"golang.org/x/crypto/bcrypt"
)

func hAdminsList(userStore auth.UserStore, app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		admins, err := userStore.UserListByRole(r.Context(), entry.RoleAdmin)
		if err != nil {
			return err
		}
		condos, err := app.CondoList(r.Context())
		if err != nil {
			return err
		}
		return templates.AdminsList(admins, condos).Render(r.Context(), w)
	})
}

func hAdminsNew(app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		condos, err := app.CondoList(r.Context())
		if err != nil {
			return err
		}
		return templates.AdminForm(nil, condos).Render(r.Context(), w)
	})
}

func hAdminsCreate(userStore auth.UserStore, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		if err := validateAdminInput(firstName, lastName, email, password, true); err != nil {
			return err
		}

		condoID, err := strconv.ParseInt(r.FormValue("condominium_id"), 10, 64)
		if err != nil || condoID == 0 {
			return entry.NewUserSafeError("Debe seleccionar un condominio válido")
		}

		user := &auth.User{
			CondominiumID: condoID,
			FirstName:     firstName,
			LastName:      lastName,
			Email:         email,
			Phone:         phone,
			Role:          entry.RoleAdmin,
			Enabled:       r.FormValue("enabled") == "on",
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = userStore.CreateUser(r.Context(), user, string(hash))
		if err != nil {
			return err
		}

		http.Redirect(w, r, "/super/admins", http.StatusSeeOther)
		return nil
	})
}

func hAdminsEdit(userStore auth.UserStore, app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		user, ok, err := userStore.GetByID(r.Context(), id)
		if err != nil {
			return err
		}
		if !ok {
			return util.NewErrorWithCode("Administrador no encontrado", http.StatusNotFound)
		}

		condos, err := app.CondoList(r.Context())
		if err != nil {
			return err
		}

		return templates.AdminForm(user, condos).Render(r.Context(), w)
	})
}

func hAdminsUpdate(userStore auth.UserStore, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		if err := r.ParseForm(); err != nil {
			return err
		}

		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		if err := validateAdminInput(firstName, lastName, email, password, false); err != nil {
			return err
		}

		currentUser := entry.UserFromCtx(r.Context())

		condoID, err := strconv.ParseInt(r.FormValue("condominium_id"), 10, 64)
		if err != nil || condoID == 0 {
			return entry.NewUserSafeError("Debe seleccionar un condominio válido")
		}

		user := &auth.User{
			CondominiumID: condoID,
			FirstName:     firstName,
			LastName:      lastName,
			Email:         email,
			Phone:         phone,
			Enabled:       r.FormValue("enabled") == "on",
		}

		_, err = userStore.UserUpdate(r.Context(), id, user, currentUser.ID)
		if err != nil {
			return err
		}

		if password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			if err := userStore.UserUpdatePassword(r.Context(), id, string(hash), currentUser.ID); err != nil {
				return err
			}
		}

		http.Redirect(w, r, "/super/admins", http.StatusSeeOther)
		return nil
	})
}

func hAdminsDelete(userStore auth.UserStore, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		currentUser := entry.UserFromCtx(r.Context())

		if id == currentUser.ID {
			return entry.NewUserSafeError("No puedes eliminar tu propia cuenta")
		}

		if err := userStore.UserDelete(r.Context(), id); err != nil {
			return err
		}

		http.Redirect(w, r, "/super/admins", http.StatusSeeOther)
		return nil
	})
}

func validateAdminInput(firstName, lastName, email, password string, requirePassword bool) error {
	if len(firstName) == 0 || len(firstName) > 100 {
		return entry.NewUserSafeError("El nombre debe tener entre 1 y 100 caracteres")
	}
	if len(lastName) == 0 || len(lastName) > 100 {
		return entry.NewUserSafeError("El apellido debe tener entre 1 y 100 caracteres")
	}
	if email == "" {
		return entry.NewUserSafeError("El email es requerido")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return entry.NewUserSafeError("El email no tiene un formato válido")
	}
	if len(email) > 255 {
		return entry.NewUserSafeError("El email no puede exceder 255 caracteres")
	}
	if requirePassword {
		if password == "" {
			return entry.NewUserSafeError("La contraseña es requerida")
		}
		if len(password) < 8 {
			return entry.NewUserSafeError("La contraseña debe tener al menos 8 caracteres")
		}
		if len(password) > 72 {
			return entry.NewUserSafeError("La contraseña no puede exceder 72 caracteres")
		}
	} else if password != "" {
		if len(password) < 8 {
			return entry.NewUserSafeError("La contraseña debe tener al menos 8 caracteres")
		}
		if len(password) > 72 {
			return entry.NewUserSafeError("La contraseña no puede exceder 72 caracteres")
		}
	}
	return nil
}
