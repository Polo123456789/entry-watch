package superadmin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/util"
	templates "github.com/Polo123456789/entry-watch/internal/templates/superadmin"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandlers struct {
	store AdminStore
}

func NewAdminHandlers(store AdminStore) *AdminHandlers {
	return &AdminHandlers{store: store}
}

func (h *AdminHandlers) List(logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		admins, err := h.store.UserListByRole(r.Context(), entry.RoleAdmin)
		if err != nil {
			return err
		}
		condos, err := h.store.CondoList(r.Context())
		if err != nil {
			return err
		}
		return templates.AdminsList(admins, condos).Render(r.Context(), w)
	})
}

func (h *AdminHandlers) New(logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		condos, err := h.store.CondoList(r.Context())
		if err != nil {
			return err
		}
		return templates.AdminForm(nil, condos).Render(r.Context(), w)
	})
}

func (h *AdminHandlers) Create(logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		condoID, _ := parseInt64(r.FormValue("condominium_id"))

		admin := &entry.AdminUser{
			FirstName:     r.FormValue("first_name"),
			LastName:      r.FormValue("last_name"),
			Email:         r.FormValue("email"),
			Phone:         r.FormValue("phone"),
			CondominiumID: condoID,
			Enabled:       r.FormValue("enabled") == "on",
		}

		password := r.FormValue("password")
		if password == "" {
			return entry.NewUserSafeError("La contraseña es requerida")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = h.store.UserCreate(r.Context(), admin, string(hash))
		if err != nil {
			return err
		}

		http.Redirect(w, r, "/super/admins", http.StatusSeeOther)
		return nil
	})
}

func (h *AdminHandlers) Edit(logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		admin, ok, err := h.store.UserGetByID(r.Context(), id)
		if err != nil {
			return err
		}
		if !ok {
			return util.NewErrorWithCode("Administrador no encontrado", http.StatusNotFound)
		}

		condos, err := h.store.CondoList(r.Context())
		if err != nil {
			return err
		}

		return templates.AdminForm(admin, condos).Render(r.Context(), w)
	})
}

func (h *AdminHandlers) Update(logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		if err := r.ParseForm(); err != nil {
			return err
		}

		user := entry.UserFromCtx(r.Context())
		condoID, _ := parseInt64(r.FormValue("condominium_id"))

		admin := &entry.AdminUser{
			FirstName:     r.FormValue("first_name"),
			LastName:      r.FormValue("last_name"),
			Email:         r.FormValue("email"),
			Phone:         r.FormValue("phone"),
			CondominiumID: condoID,
			Enabled:       r.FormValue("enabled") == "on",
		}

		_, err = h.store.UserUpdate(r.Context(), id, admin, user.ID)
		if err != nil {
			return err
		}

		http.Redirect(w, r, "/super/admins", http.StatusSeeOther)
		return nil
	})
}

func (h *AdminHandlers) Delete(logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		if err := h.store.UserDelete(r.Context(), id); err != nil {
			return err
		}

		http.Redirect(w, r, "/super/admins", http.StatusSeeOther)
		return nil
	})
}
