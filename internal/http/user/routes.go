package user

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

	mux.Handle("/neighbor/", hGet(app, logger))

	var handler http.Handler = mux
	handler = authMiddleware(handler, app, logger)
	return handler
}

func authMiddleware(
	next http.Handler,
	app *entry.App,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := entry.RequireRole(r.Context(), entry.RoleUser)
		if err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
