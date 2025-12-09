package user

import (
	"log/slog"
	"net/http"

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
		return templates.Dashboard().Render(r.Context(), w)
	})
}
