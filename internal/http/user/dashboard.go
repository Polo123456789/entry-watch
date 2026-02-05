package user

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/util"
	templates "github.com/Polo123456789/entry-watch/internal/templates/user"
)

func hGet(
	app *entry.App,
	logger *slog.Logger,
) http.Handler {
	return util.Handler(logger, func(
		w http.ResponseWriter, r *http.Request,
	) error {
		today := time.Now()

		return templates.Dashboard(entry.Visit{
			MaxUses:   1,
			ValidFrom: today,
			ValidTo:   today,
		}).Render(r.Context(), w)
	})
}

func hPost(
	app *entry.App,
	logger *slog.Logger,
) http.Handler {
	return util.Handler(logger, func(
		w http.ResponseWriter, r *http.Request,
	) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		visitor := r.FormValue("visitor_name")
		limitUses := r.FormValue("limit_uses") == "on"
		maxUses := 0

		if limitUses {
			var err error
			maxUses, err = strconv.Atoi(r.FormValue("max_uses"))
			if err != nil {
				return err
			}
		}

		validFrom, err := time.Parse(
			time.DateOnly, r.FormValue("valid_from"),
		)
		if err != nil {
			return err
		}

		validTo, err := time.Parse(
			time.DateOnly, r.FormValue("valid_to"),
		)
		if err != nil {
			return err
		}

		// visit, err := app.CreateVisit( ... )
		// http.Redirect( ... )
		// TODO: Come back to this
		return nil
	})
}
