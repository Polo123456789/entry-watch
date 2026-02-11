package admin

import (
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/util"
)

func Handle(
	app *entry.App,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	// Setup routes
	mux.Handle("/admin/", hGet(app, logger))

	var handler http.Handler = mux
	handler = authMiddleware(handler, logger)
	return handler
}

func authMiddleware(
	next http.Handler,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := entry.RequireRole(r.Context(), entry.RoleAdmin)
		if err != nil {
			util.HandleError(w, r, logger, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
