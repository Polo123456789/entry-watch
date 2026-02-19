package superadmin

import (
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
	"github.com/Polo123456789/entry-watch/internal/http/util"
)

func Handle(
	app *entry.App,
	logger *slog.Logger,
	userStore auth.UserStore,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /super/", hGet(app, logger))

	mux.Handle("GET /super/condos", hCondosList(app, logger))
	mux.Handle("GET /super/condos/new", hCondosNew(app, logger))
	mux.Handle("POST /super/condos", hCondosCreate(app, logger))
	mux.Handle("GET /super/condos/{id}/edit", hCondosEdit(app, logger))
	mux.Handle("POST /super/condos/{id}", hCondosUpdate(app, logger))
	mux.Handle("POST /super/condos/{id}/delete", hCondosDelete(app, userStore, logger))

	mux.Handle("GET /super/admins", hAdminsList(userStore, logger))
	mux.Handle("GET /super/admins/new", hAdminsNew(app, logger))
	mux.Handle("POST /super/admins", hAdminsCreate(userStore, logger))
	mux.Handle("GET /super/admins/{id}/edit", hAdminsEdit(userStore, app, logger))
	mux.Handle("POST /super/admins/{id}", hAdminsUpdate(userStore, logger))
	mux.Handle("POST /super/admins/{id}/delete", hAdminsDelete(userStore, logger))

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
