package superadmin

import (
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/util"
)

func Handle(
	app *entry.App,
	logger *slog.Logger,
	store AdminStore,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /super/", hGet(app, logger))

	mux.Handle("GET /super/condos", hCondosList(app, logger))
	mux.Handle("GET /super/condos/new", hCondosNew(app, logger))
	mux.Handle("POST /super/condos", hCondosCreate(app, logger))
	mux.Handle("GET /super/condos/{id}/edit", hCondosEdit(app, logger))
	mux.Handle("POST /super/condos/{id}", hCondosUpdate(app, logger))
	mux.Handle("POST /super/condos/{id}/delete", hCondosDelete(app, logger))

	adminHandlers := NewAdminHandlers(store)
	mux.Handle("GET /super/admins", adminHandlers.List(logger))
	mux.Handle("GET /super/admins/new", adminHandlers.New(logger))
	mux.Handle("POST /super/admins", adminHandlers.Create(logger))
	mux.Handle("GET /super/admins/{id}/edit", adminHandlers.Edit(logger))
	mux.Handle("POST /super/admins/{id}", adminHandlers.Update(logger))
	mux.Handle("POST /super/admins/{id}/delete", adminHandlers.Delete(logger))

	var handler http.Handler = mux
	handler = authMiddleware(handler, logger)
	return handler
}

func authMiddleware(
	next http.Handler,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := entry.RequireRole(r.Context(), entry.RoleSuperAdmin)
		if err != nil {
			util.HandleError(w, r, logger, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
