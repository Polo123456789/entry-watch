package superadmin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/util"
	templates "github.com/Polo123456789/entry-watch/internal/templates/superadmin"
)

func hCondosList(app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		condos, err := app.Store().CondoList(r.Context())
		if err != nil {
			return err
		}
		return templates.CondosList(condos).Render(r.Context(), w)
	})
}

func hCondosNew(app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		return templates.CondoForm(nil).Render(r.Context(), w)
	})
}

func hCondosCreate(app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		name := r.FormValue("name")
		address := r.FormValue("address")

		if err := validateCondoInput(name, address); err != nil {
			return err
		}

		user := entry.UserFromCtx(r.Context())

		condo := &entry.Condominium{
			Name:      name,
			Address:   address,
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		}

		_, err := app.Store().CondoCreate(r.Context(), condo)
		if err != nil {
			return err
		}

		http.Redirect(w, r, "/super/condos", http.StatusSeeOther)
		return nil
	})
}

func hCondosEdit(app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		condo, err := app.Store().CondoGetByID(r.Context(), id)
		if err != nil {
			return err
		}
		if condo == nil {
			return util.NewErrorWithCode("Condominio no encontrado", http.StatusNotFound)
		}

		return templates.CondoForm(condo).Render(r.Context(), w)
	})
}

func hCondosUpdate(app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		if err := r.ParseForm(); err != nil {
			return err
		}

		name := r.FormValue("name")
		address := r.FormValue("address")

		if err := validateCondoInput(name, address); err != nil {
			return err
		}

		user := entry.UserFromCtx(r.Context())

		err = app.Store().CondoUpdate(r.Context(), id, func(condo *entry.Condominium) (*entry.Condominium, error) {
			condo.Name = name
			condo.Address = address
			condo.UpdatedBy = user.ID
			return condo, nil
		})
		if err != nil {
			return err
		}

		http.Redirect(w, r, "/super/condos", http.StatusSeeOther)
		return nil
	})
}

func hCondosDelete(app *entry.App, logger *slog.Logger) http.Handler {
	return util.Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			return util.NewErrorWithCode("ID inválido", http.StatusBadRequest)
		}

		if err := app.Store().CondoDelete(r.Context(), id); err != nil {
			return err
		}

		http.Redirect(w, r, "/super/condos", http.StatusSeeOther)
		return nil
	})
}

func validateCondoInput(name, address string) error {
	if len(name) == 0 || len(name) > 200 {
		return entry.NewUserSafeError("El nombre debe tener entre 1 y 200 caracteres")
	}
	if len(address) == 0 || len(address) > 500 {
		return entry.NewUserSafeError("La dirección debe tener entre 1 y 500 caracteres")
	}
	return nil
}
