package http

import (
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/app"
	"github.com/Polo123456789/entry-watch/internal/http/admin"
	"github.com/Polo123456789/entry-watch/internal/http/guard"
	"github.com/Polo123456789/entry-watch/internal/http/superadmin"
	"github.com/Polo123456789/entry-watch/internal/http/user"
	"github.com/Polo123456789/entry-watch/web"
)

func setupRoutes(
	mux *http.ServeMux,
	app *app.App,
	logger *slog.Logger,
) {
	mux.Handle("/super/", superadmin.Handle(app, logger))
	mux.Handle("/admin/", admin.Handle(app, logger))
	mux.Handle("/guard/", guard.Handle(app, logger))
	mux.Handle("/neighbor/", user.Handle(app, logger))
	mux.Handle("GET /static/", http.FileServerFS(web.StaticFiles))
}
