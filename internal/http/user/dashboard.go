package user

import (
	"log/slog"
	"net/http"
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

// TODO: Make the POST, and the display of the generated code
