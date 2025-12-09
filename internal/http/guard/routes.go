package guard

import (
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

func Handle(
	app *entry.App,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	// Setup routes
	// Middlerewares

	return mux
}
