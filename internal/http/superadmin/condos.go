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
		condos, err := app.CondoList(r.Context())
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

		_, err := app.CondoCreate(
			r.Context(),
			r.FormValue("name"),
			r.FormValue("address"),
		)
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

		condo, ok, err := app.CondoGetByID(r.Context(), id)
		if err != nil {
			return err
		}
		if !ok {
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

		err = app.CondoUpdate(
			r.Context(),
			id,
			r.FormValue("name"),
			r.FormValue("address"),
		)
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

		if err := app.CondoDelete(r.Context(), id); err != nil {
			return err
		}

		http.Redirect(w, r, "/super/condos", http.StatusSeeOther)
		return nil
	})
}
