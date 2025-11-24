package guard

import (
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/app"
)

func Handle(
	app *app.App,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	// Setup routes
	// Middlerewares

	return mux
}
