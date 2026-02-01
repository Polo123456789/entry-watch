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

// hPost handles the POST request for creating a new visit.
// Currently unused - will be wired up when visit creation is implemented.
//
//nolint:unused
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

		_ = r.FormValue("visitor_name")
		limitUses := r.FormValue("limit_uses") == "on"
		_ = 0

		if limitUses {
			_, err := strconv.Atoi(r.FormValue("max_uses"))
			if err != nil {
				return err
			}
		}

		_, err := time.Parse(
			time.DateOnly, r.FormValue("valid_from"),
		)
		if err != nil {
			return err
		}

		_, err = time.Parse(
			time.DateOnly, r.FormValue("valid_to"),
		)
		if err != nil {
			return err
		}

		// visit, err := app.CreateVisit( ... )
		// http.Redirect( ... )
		// TODO: Come back to this
		_ = app // Use app to avoid unused parameter warning
		return nil
	})
}
